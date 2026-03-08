<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import ConfirmButton from '$components/common/ConfirmButton.svelte';
	import FormSection from '$components/common/FormSection.svelte';
	import Pager from '$components/common/Pager.svelte';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';
	import { pushToast } from '$lib/stores/toast';

	type Workspace = {
		id: number;
		name: string;
		description?: string | null;
	};

	const PAGE_SIZE = 10;
	let items = $state<Workspace[]>([]);
	let total = $state(0);
	let offset = $state(0);
	let hasNext = $state(false);
	let search = $state('');

	let name = $state('');
	let description = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	let editingId = $state<number | null>(null);
	let editName = $state('');
	let editDescription = $state('');
	let editSaving = $state(false);
	let deletingId = $state<number | null>(null);

	async function load(nextOffset = offset) {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.workspaces.workspacesList(PAGE_SIZE, nextOffset, search || undefined));
			const envelope = (res as { data?: Workspace[]; meta?: { total?: number; has_next?: boolean; offset?: number } }) ?? {};
			items = envelope.data ?? [];
			total = envelope.meta?.total ?? 0;
			hasNext = envelope.meta?.has_next ?? false;
			offset = envelope.meta?.offset ?? nextOffset;
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load workspaces');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		void load(0);
	});

	async function createWorkspace(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		saving = true;
		try {
			await withApi(() => api.workspaces.workspacesCreate({ name, description: description || null }));
			name = '';
			description = '';
			pushToast({ type: 'success', title: 'Workspace created' });
			await load(0);
		} catch (e) {
			error = userMessageFromError(e, 'Create failed');
		} finally {
			saving = false;
		}
	}

	async function runSearch(e: SubmitEvent) {
		e.preventDefault();
		await load(0);
	}

	function startEdit(workspace: Workspace) {
		editingId = workspace.id;
		editName = workspace.name;
		editDescription = workspace.description ?? '';
	}

	function cancelEdit() {
		editingId = null;
		editName = '';
		editDescription = '';
	}

	async function saveEdit(workspaceId: number) {
		error = '';
		editSaving = true;
		try {
			await withApi(() =>
				api.workspaces.workspacesPatch(workspaceId, {
					workspace_id: workspaceId,
					name: editName || null,
					description: editDescription || null
				})
			);
			cancelEdit();
			pushToast({ type: 'success', title: 'Workspace updated' });
			await load(offset);
		} catch (e) {
			error = userMessageFromError(e, 'Update failed');
		} finally {
			editSaving = false;
		}
	}

	async function removeWorkspace(workspaceId: number) {
		error = '';
		deletingId = workspaceId;
		try {
			await withApi(() => api.workspaces.workspacesDelete(workspaceId));
			const needBack = items.length === 1 && offset > 0;
			const nextOffset = needBack ? Math.max(offset - PAGE_SIZE, 0) : offset;
			pushToast({ type: 'success', title: 'Workspace deleted' });
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

<main class="mx-auto w-full max-w-6xl p-6">
	<h1 class="text-2xl font-semibold">Workspaces</h1>
	<div class="mt-4 space-y-3">
		<FormSection title="Search" description="Find workspaces by name">
			<form class="grid gap-2 sm:grid-cols-[1fr_auto]" onsubmit={runSearch}>
				<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Search workspaces" bind:value={search} />
				<Button type="submit" variant="outline">Search</Button>
			</form>
		</FormSection>

		<FormSection title="Create Workspace" description="Create a new workspace">
			<form class="grid gap-2 sm:grid-cols-[1fr_1fr_auto]" onsubmit={createWorkspace}>
				<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Workspace name" bind:value={name} required />
				<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Description (optional)" bind:value={description} />
				<Button type="submit" disabled={saving}>{saving ? 'Creating...' : 'Create'}</Button>
			</form>
		</FormSection>
	</div>

	{#if error}<StateBlock kind="error" message={error} />{/if}

	{#if loading}
		<StateBlock message="Loading..." />
	{:else if items.length === 0}
		<StateBlock kind="muted" message="No workspaces yet." />
	{:else}
		<ul class="mt-4 space-y-2">
			{#each items as w}
				<li class="rounded-lg border border-border p-3">
					{#if editingId === w.id}
						<div class="grid gap-2 sm:grid-cols-[1fr_1fr_auto_auto]">
							<input class="rounded-md border border-border bg-background px-3 py-2" bind:value={editName} />
							<input class="rounded-md border border-border bg-background px-3 py-2" bind:value={editDescription} />
							<Button onclick={() => saveEdit(w.id)} disabled={editSaving}>{editSaving ? 'Saving...' : 'Save'}</Button>
							<Button variant="outline" onclick={cancelEdit}>Cancel</Button>
						</div>
					{:else}
						<div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
							<div>
								<a class="font-medium underline" href={`/app/workspaces/${w.id}`}>{w.name}</a>
								{#if w.description}
									<p class="text-sm text-muted-foreground">{w.description}</p>
								{/if}
							</div>
							<div class="flex gap-2">
								<Button variant="outline" onclick={() => startEdit(w)}>Edit</Button>
								<ConfirmButton
									label="Delete"
									title="Delete workspace?"
									description={`Workspace '${w.name}' will be deleted.`}
									confirmLabel="Delete"
									pending={deletingId === w.id}
									onConfirm={() => removeWorkspace(w.id)}
								/>
							</div>
						</div>
					{/if}
				</li>
			{/each}
		</ul>

		<Pager {offset} pageSize={PAGE_SIZE} {total} {hasNext} onPrev={prevPage} onNext={nextPage} />
	{/if}
</main>
