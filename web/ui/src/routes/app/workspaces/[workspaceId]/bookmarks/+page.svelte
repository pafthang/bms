<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import ConfirmButton from '$components/common/ConfirmButton.svelte';
	import FormSection from '$components/common/FormSection.svelte';
	import Pager from '$components/common/Pager.svelte';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';
	import { pushToast } from '$lib/stores/toast';
	let { data } = $props();
	const workspaceId = $derived(data.workspaceId as number);

	type Bookmark = {
		id: number;
		title: string;
		url: string;
		description?: string | null;
		is_archived: boolean;
		tag_ids?: number[];
	};
	type Tag = { id: number; name: string };

	const PAGE_SIZE = 10;
	let items = $state<Bookmark[]>([]);
	let tags = $state<Tag[]>([]);
	let total = $state(0);
	let offset = $state(0);
	let hasNext = $state(false);

	let title = $state('');
	let url = $state('');
	let description = $state('');
	let newTagIds = $state<number[]>([]);

	let search = $state('');
	let archived = $state<'all' | 'active' | 'archived'>('all');
	let tagId = $state<number | null>(null);
	let sort = $state('created_at');
	let order = $state('desc');

	let error = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let deletingId = $state<number | null>(null);
	let patchingId = $state<number | null>(null);

	let editingId = $state<number | null>(null);
	let editTitle = $state('');
	let editUrl = $state('');
	let editDescription = $state('');
	let editTagIds = $state<number[]>([]);
	let editSaving = $state(false);

	function toggleId(source: number[], id: number): number[] {
		return source.includes(id) ? source.filter((x) => x !== id) : [...source, id];
	}

	function archivedFilterValue(): boolean | null {
		if (archived === 'active') return false;
		if (archived === 'archived') return true;
		return null;
	}

	async function load(nextOffset = offset) {
		loading = true;
		error = '';
		try {
			const [bookmarksRes, tagsRes] = await Promise.all([
				withApi(() =>
					api.bookmarks.bookmarksList(
						workspaceId,
						PAGE_SIZE,
						nextOffset,
						search || undefined,
						undefined,
						tagId,
						archivedFilterValue(),
						sort,
						order
					)
				),
				withApi(() => api.tags.tagsList(workspaceId, 100, 0))
			]);
			const b = (bookmarksRes as { data?: Bookmark[]; meta?: { total?: number; has_next?: boolean; offset?: number } }) ?? {};
			items = b.data ?? [];
			total = b.meta?.total ?? 0;
			hasNext = b.meta?.has_next ?? false;
			offset = b.meta?.offset ?? nextOffset;
			tags = (tagsRes as { data?: Tag[] }).data ?? [];
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load bookmarks');
		} finally {
			loading = false;
		}
	}

	onMount(load);

	async function createBookmark(e: SubmitEvent) {
		e.preventDefault();
		saving = true;
		error = '';
		try {
			await withApi(() =>
				api.bookmarks.bookmarksCreate(workspaceId, {
					workspace_id: workspaceId,
					title,
					url,
					description: description || null,
					tag_ids: newTagIds,
					is_archived: false
				})
			);
			title = '';
			url = '';
			description = '';
			newTagIds = [];
			pushToast({ type: 'success', title: 'Bookmark created' });
			await load(0);
		} catch (e) {
			error = userMessageFromError(e, 'Create failed');
		} finally {
			saving = false;
		}
	}

	async function applyFilters(e: SubmitEvent) {
		e.preventDefault();
		await load(0);
	}

	async function clearFilters() {
		search = '';
		archived = 'all';
		tagId = null;
		sort = 'created_at';
		order = 'desc';
		await load(0);
	}

	function startEdit(item: Bookmark) {
		editingId = item.id;
		editTitle = item.title;
		editUrl = item.url;
		editDescription = item.description ?? '';
		editTagIds = item.tag_ids ?? [];
	}

	function cancelEdit() {
		editingId = null;
		editTitle = '';
		editUrl = '';
		editDescription = '';
		editTagIds = [];
	}

	async function saveEdit(bookmarkId: number) {
		error = '';
		editSaving = true;
		try {
			await withApi(() =>
				api.bookmarks.bookmarksPatch(workspaceId, bookmarkId, {
					workspace_id: workspaceId,
					bookmark_id: bookmarkId,
					title: editTitle || null,
					url: editUrl || null,
					description: editDescription || null,
					tag_ids: editTagIds
				})
			);
			cancelEdit();
			pushToast({ type: 'success', title: 'Bookmark updated' });
			await load(offset);
		} catch (e) {
			error = userMessageFromError(e, 'Update failed');
		} finally {
			editSaving = false;
		}
	}

	async function toggleArchive(item: Bookmark) {
		error = '';
		patchingId = item.id;
		try {
			await withApi(() =>
				api.bookmarks.bookmarksPatch(workspaceId, item.id, {
					workspace_id: workspaceId,
					bookmark_id: item.id,
					is_archived: !item.is_archived
				})
			);
			pushToast({ type: 'success', title: item.is_archived ? 'Bookmark unarchived' : 'Bookmark archived' });
			await load(offset);
		} catch (e) {
			error = userMessageFromError(e, 'Archive update failed');
		} finally {
			patchingId = null;
		}
	}

	async function removeBookmark(id: number) {
		deletingId = id;
		try {
			await withApi(() => api.bookmarks.bookmarksDelete(workspaceId, id));
			const needBack = items.length === 1 && offset > 0;
			const nextOffset = needBack ? Math.max(offset - PAGE_SIZE, 0) : offset;
			pushToast({ type: 'success', title: 'Bookmark deleted' });
			await load(nextOffset);
		} catch (e) {
			error = userMessageFromError(e, 'Delete failed');
		} finally {
			deletingId = null;
		}
	}

	async function prevPage() {
		if (offset === 0) return;
		await load(Math.max(offset - PAGE_SIZE, 0));
	}

	async function nextPage() {
		if (!hasNext) return;
		await load(offset + PAGE_SIZE);
	}
</script>

<h2 class="text-xl font-semibold">Bookmarks</h2>

<div class="mt-4 space-y-3">
	<FormSection title="Filters" description="Search and filter bookmarks">
		<form class="grid gap-2 sm:grid-cols-[1fr_120px_120px_140px_120px_auto_auto]" onsubmit={applyFilters}>
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Search" bind:value={search} />
			<select class="rounded-md border border-border bg-background px-3 py-2" bind:value={archived}>
				<option value="all">all</option>
				<option value="active">active</option>
				<option value="archived">archived</option>
			</select>
			<select
				class="rounded-md border border-border bg-background px-3 py-2"
				value={tagId === null ? '' : String(tagId)}
				onchange={(e) => {
					const value = (e.currentTarget as HTMLSelectElement).value;
					tagId = value ? Number(value) : null;
				}}
			>
				<option value="">all tags</option>
				{#each tags as t}
					<option value={t.id}>{t.name}</option>
				{/each}
			</select>
			<select class="rounded-md border border-border bg-background px-3 py-2" bind:value={sort}>
				<option value="created_at">created_at</option>
				<option value="updated_at">updated_at</option>
				<option value="title">title</option>
			</select>
			<select class="rounded-md border border-border bg-background px-3 py-2" bind:value={order}>
				<option value="desc">desc</option>
				<option value="asc">asc</option>
			</select>
			<Button type="submit" variant="outline">Apply</Button>
			<Button type="button" variant="outline" onclick={clearFilters}>Reset</Button>
		</form>
	</FormSection>

	<FormSection title="Create Bookmark" description="Add a link to this workspace">
		<form class="grid gap-2 sm:grid-cols-[1fr_1fr_1fr_auto]" onsubmit={createBookmark}>
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Title" bind:value={title} required />
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="https://..." type="url" bind:value={url} required />
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Description (optional)" bind:value={description} />
			<Button type="submit" disabled={saving}>{saving ? 'Creating...' : 'Add'}</Button>
		</form>
	</FormSection>
</div>

{#if tags.length > 0}
	<div class="mt-2 flex flex-wrap gap-2">
		{#each tags as tag}
			<label class="inline-flex items-center gap-2 rounded border border-border px-2 py-1 text-sm">
				<input
					type="checkbox"
					checked={newTagIds.includes(tag.id)}
					onchange={() => {
						newTagIds = toggleId(newTagIds, tag.id);
					}}
				/>
				{tag.name}
			</label>
		{/each}
	</div>
{/if}

{#if error}<StateBlock kind="error" message={error} />{/if}

{#if loading}
	<StateBlock message="Loading..." />
{:else}
	<ul class="mt-4 space-y-2">
		{#each items as b}
			<li class="rounded-lg border border-border p-3">
				{#if editingId === b.id}
					<div class="grid gap-2 sm:grid-cols-[1fr_1fr]">
						<input class="rounded-md border border-border bg-background px-3 py-2" bind:value={editTitle} />
						<input class="rounded-md border border-border bg-background px-3 py-2" bind:value={editUrl} />
					</div>
					<input class="mt-2 w-full rounded-md border border-border bg-background px-3 py-2" bind:value={editDescription} />
					{#if tags.length > 0}
						<div class="mt-2 flex flex-wrap gap-2">
							{#each tags as tag}
								<label class="inline-flex items-center gap-2 rounded border border-border px-2 py-1 text-sm">
									<input
										type="checkbox"
										checked={editTagIds.includes(tag.id)}
										onchange={() => {
											editTagIds = toggleId(editTagIds, tag.id);
										}}
									/>
									{tag.name}
								</label>
							{/each}
						</div>
					{/if}
					<div class="mt-2 flex gap-2">
						<Button onclick={() => saveEdit(b.id)} disabled={editSaving}>{editSaving ? 'Saving...' : 'Save'}</Button>
						<Button variant="outline" onclick={cancelEdit}>Cancel</Button>
					</div>
				{:else}
					<div class="flex items-start justify-between gap-3">
						<div>
							<p class="font-medium">
								{b.title}
								{#if b.is_archived}
									<span class="ml-2 rounded border border-border px-2 py-0.5 text-xs text-muted-foreground">archived</span>
								{/if}
							</p>
							<a class="text-sm text-muted-foreground underline" href={b.url} target="_blank" rel="noreferrer">{b.url}</a>
							{#if b.description}
								<p class="mt-1 text-sm text-muted-foreground">{b.description}</p>
							{/if}
						</div>
						<div class="flex gap-2">
							<Button variant="outline" onclick={() => startEdit(b)}>Edit</Button>
							<Button variant="outline" onclick={() => toggleArchive(b)} disabled={patchingId === b.id}>
								{patchingId === b.id ? 'Updating...' : b.is_archived ? 'Unarchive' : 'Archive'}
							</Button>
							<ConfirmButton
								label="Delete"
								title="Delete bookmark?"
								description={`Bookmark '${b.title}' will be deleted.`}
								confirmLabel="Delete"
								pending={deletingId === b.id}
								onConfirm={() => removeBookmark(b.id)}
							/>
						</div>
					</div>
				{/if}
			</li>
		{/each}
	</ul>

	<Pager {offset} pageSize={PAGE_SIZE} {total} {hasNext} onPrev={prevPage} onNext={nextPage} />
{/if}
