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

	type Tag = { id: number; name: string; color?: string | null };
	const PAGE_SIZE = 12;
	const HEX_COLOR_RE = /^#([0-9a-fA-F]{6})$/;

	let items = $state<Tag[]>([]);
	let total = $state(0);
	let offset = $state(0);
	let hasNext = $state(false);

	let name = $state('');
	let color = $state('#3b82f6');
	let search = $state('');
	let error = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let deletingId = $state<number | null>(null);

	let editingId = $state<number | null>(null);
	let editName = $state('');
	let editColor = $state('#3b82f6');
	let editSaving = $state(false);

	function isHexColor(value: string): boolean {
		return HEX_COLOR_RE.test(value);
	}

	async function load(nextOffset = offset) {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.tags.tagsList(workspaceId, PAGE_SIZE, nextOffset));
			const envelope = (res as { data?: Tag[]; meta?: { total?: number; offset?: number; has_next?: boolean } }) ?? {};
			items = envelope.data ?? [];
			total = envelope.meta?.total ?? 0;
			offset = envelope.meta?.offset ?? nextOffset;
			hasNext = envelope.meta?.has_next ?? false;
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load tags');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		void load(0);
	});

	async function createTag(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		if (!isHexColor(color)) {
			error = 'Color must be a valid hex value (example: #3b82f6).';
			return;
		}
		saving = true;
		try {
			await withApi(() =>
				api.tags.tagsCreate(workspaceId, {
					workspace_id: workspaceId,
					name,
					color
				})
			);
			name = '';
			pushToast({ type: 'success', title: 'Tag created' });
			await load(0);
		} catch (e) {
			error = userMessageFromError(e, 'Create failed');
		} finally {
			saving = false;
		}
	}

	function startEdit(tag: Tag) {
		editingId = tag.id;
		editName = tag.name;
		editColor = tag.color ?? '#3b82f6';
	}

	function cancelEdit() {
		editingId = null;
		editName = '';
		editColor = '#3b82f6';
	}

	async function saveEdit(tagId: number) {
		error = '';
		if (!isHexColor(editColor)) {
			error = 'Color must be a valid hex value (example: #3b82f6).';
			return;
		}
		editSaving = true;
		try {
			await withApi(() =>
				api.tags.tagsPatch(workspaceId, tagId, {
					workspace_id: workspaceId,
					tag_id: tagId,
					name: editName || null,
					color: editColor
				})
			);
			cancelEdit();
			pushToast({ type: 'success', title: 'Tag updated' });
			await load(offset);
		} catch (e) {
			error = userMessageFromError(e, 'Update failed');
		} finally {
			editSaving = false;
		}
	}

	async function removeTag(tagId: number) {
		error = '';
		deletingId = tagId;
		try {
			await withApi(() => api.tags.tagsDelete(workspaceId, tagId));
			const needBack = items.length === 1 && offset > 0;
			pushToast({ type: 'success', title: 'Tag deleted' });
			await load(needBack ? Math.max(offset - PAGE_SIZE, 0) : offset);
		} catch (e) {
			error = userMessageFromError(e, 'Delete failed');
		} finally {
			deletingId = null;
		}
	}

	async function applySearch(e: SubmitEvent) {
		e.preventDefault();
		await load(0);
	}

	const filteredItems = $derived.by(() => {
		const q = search.trim().toLowerCase();
		if (!q) return items;
		return items.filter((item) => item.name.toLowerCase().includes(q));
	});
</script>

<h2 class="text-xl font-semibold">Tags</h2>

<div class="mt-4 space-y-3">
	<FormSection title="Search" description="Filter tags by name">
		<form class="grid gap-2 sm:grid-cols-[1fr_auto_auto]" onsubmit={applySearch}>
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Search by tag name" bind:value={search} />
			<Button type="submit" variant="outline">Apply</Button>
			<Button type="button" variant="outline" onclick={() => { search = ''; void load(0); }}>Reset</Button>
		</form>
	</FormSection>

	<FormSection title="Create Tag" description="Add a new tag to this workspace">
		<form class="grid gap-2 sm:grid-cols-[1fr_120px_auto]" onsubmit={createTag}>
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Tag name" bind:value={name} required />
			<input class="h-10 rounded-md border border-border bg-background px-2" type="color" bind:value={color} />
			<Button type="submit" disabled={saving}>{saving ? 'Creating...' : 'Add'}</Button>
		</form>
	</FormSection>
</div>

{#if error}<StateBlock kind="error" message={error} />{/if}

{#if loading}
	<StateBlock message="Loading..." />
{:else if filteredItems.length === 0}
	<StateBlock kind="muted" message="No tags found." />
{:else}
	<ul class="mt-4 space-y-2">
		{#each filteredItems as tag}
			<li class="rounded-lg border border-border p-3">
				{#if editingId === tag.id}
					<div class="grid gap-2 sm:grid-cols-[1fr_120px_auto_auto]">
						<input class="rounded-md border border-border bg-background px-3 py-2" bind:value={editName} />
						<input class="h-10 rounded-md border border-border bg-background px-2" type="color" bind:value={editColor} />
						<Button onclick={() => saveEdit(tag.id)} disabled={editSaving}>{editSaving ? 'Saving...' : 'Save'}</Button>
						<Button variant="outline" onclick={cancelEdit}>Cancel</Button>
					</div>
				{:else}
					<div class="flex items-center justify-between gap-3">
						<div class="flex items-center gap-2">
							<span class="inline-block size-3 rounded-full border" style={`background:${tag.color ?? '#9ca3af'}`}></span>
							<p class="font-medium">{tag.name}</p>
						</div>
						<div class="flex gap-2">
							<Button variant="outline" onclick={() => startEdit(tag)}>Edit</Button>
							<ConfirmButton
								label="Delete"
								title="Delete tag?"
								description={`Tag '${tag.name}' will be deleted.`}
								confirmLabel="Delete"
								pending={deletingId === tag.id}
								onConfirm={() => removeTag(tag.id)}
							/>
						</div>
					</div>
				{/if}
			</li>
		{/each}
	</ul>

	{#if search.trim()}
		<p class="mt-4 text-sm text-muted-foreground">Filtered: {filteredItems.length}</p>
	{:else}
		<Pager
			{offset}
			pageSize={PAGE_SIZE}
			{total}
			{hasNext}
			onPrev={() => load(Math.max(offset - PAGE_SIZE, 0))}
			onNext={() => load(offset + PAGE_SIZE)}
		/>
	{/if}
{/if}
