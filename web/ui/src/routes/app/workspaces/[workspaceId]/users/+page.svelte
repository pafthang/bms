<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import ConfirmButton from '$components/common/ConfirmButton.svelte';
	import FormSection from '$components/common/FormSection.svelte';
	import Pager from '$components/common/Pager.svelte';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';
	import { pushToast } from '$lib/stores/toast';
	import { WorkspaceUserCreateRequest, WorkspaceUserPatchRequest } from '@pafthang/bms-sdk';
	let { data } = $props();
	const workspaceId = $derived(data.workspaceId as number);

	type Member = { user_id: number; email: string; role: string };
	const PAGE_SIZE = 10;

	let items = $state<Member[]>([]);
	let total = $state(0);
	let offset = $state(0);
	let hasNext = $state(false);

	let search = $state('');
	let email = $state('');
	let role = $state<WorkspaceUserCreateRequest.role>(WorkspaceUserCreateRequest.role.VIEWER);
	let error = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let patchingId = $state<number | null>(null);
	let deletingId = $state<number | null>(null);

	async function load(nextOffset = offset) {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.workspaceUsers.workspaceUsersList(workspaceId, PAGE_SIZE, nextOffset));
			const envelope = (res as { data?: Member[]; meta?: { total?: number; offset?: number; has_next?: boolean } }) ?? {};
			items = envelope.data ?? [];
			total = envelope.meta?.total ?? 0;
			offset = envelope.meta?.offset ?? nextOffset;
			hasNext = envelope.meta?.has_next ?? false;
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load users');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		void load(0);
	});

	async function addUser(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		saving = true;
		try {
			await withApi(() =>
				api.workspaceUsers.workspaceUsersAdd(workspaceId, {
					workspace_id: workspaceId,
					email,
					role
				})
			);
			email = '';
			role = WorkspaceUserCreateRequest.role.VIEWER;
			pushToast({ type: 'success', title: 'User added to workspace' });
			await load(0);
		} catch (e) {
			error = userMessageFromError(e, 'Add failed');
		} finally {
			saving = false;
		}
	}

	async function setRole(userId: number, nextRole: WorkspaceUserPatchRequest.role) {
		error = '';
		patchingId = userId;
		try {
			await withApi(() =>
				api.workspaceUsers.workspaceUsersPatch(workspaceId, userId, {
					workspace_id: workspaceId,
					user_id: userId,
					role: nextRole
				})
			);
			pushToast({ type: 'success', title: 'Role updated' });
			await load(offset);
		} catch (e) {
			error = userMessageFromError(e, 'Role change failed');
		} finally {
			patchingId = null;
		}
	}

	async function removeUser(userId: number) {
		error = '';
		deletingId = userId;
		try {
			await withApi(() => api.workspaceUsers.workspaceUsersDelete(workspaceId, userId));
			const needBack = items.length === 1 && offset > 0;
			pushToast({ type: 'success', title: 'User removed from workspace' });
			await load(needBack ? Math.max(offset - PAGE_SIZE, 0) : offset);
		} catch (e) {
			error = userMessageFromError(e, 'Delete failed');
		} finally {
			deletingId = null;
		}
	}

	const filteredItems = $derived.by(() => {
		const q = search.trim().toLowerCase();
		if (!q) return items;
		return items.filter((item) => (item.email || '').toLowerCase().includes(q));
	});
</script>

<h2 class="text-xl font-semibold">Workspace Users</h2>

<div class="mt-4 space-y-3">
	<FormSection title="Search" description="Filter members by email">
		<form class="grid gap-2 sm:grid-cols-[1fr_auto]" onsubmit={(e) => e.preventDefault()}>
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Search by email" bind:value={search} />
			<Button type="button" variant="outline" onclick={() => load(0)}>Refresh</Button>
		</form>
	</FormSection>

	<FormSection title="Add Member" description="Invite user into workspace">
		<form class="grid gap-2 sm:grid-cols-[1fr_140px_auto]" onsubmit={addUser}>
			<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="user email" type="email" bind:value={email} required />
			<select class="rounded-md border border-border bg-background px-3 py-2" bind:value={role}>
				<option value={WorkspaceUserCreateRequest.role.VIEWER}>viewer</option>
				<option value={WorkspaceUserCreateRequest.role.EDITOR}>editor</option>
				<option value={WorkspaceUserCreateRequest.role.ADMIN}>admin</option>
			</select>
			<Button type="submit" disabled={saving}>{saving ? 'Adding...' : 'Add'}</Button>
		</form>
	</FormSection>
</div>

{#if error}<StateBlock kind="error" message={error} />{/if}

{#if loading}
	<StateBlock message="Loading..." />
{:else if filteredItems.length === 0}
	<StateBlock kind="muted" message="No users found." />
{:else}
	<ul class="mt-4 space-y-2">
		{#each filteredItems as u}
			<li class="flex flex-col gap-2 rounded-lg border border-border p-3 sm:flex-row sm:items-center sm:justify-between">
				<div>
					<p class="font-medium">{u.email || `user #${u.user_id}`}</p>
					<p class="text-sm text-muted-foreground">role: {u.role}</p>
				</div>
				<div class="flex flex-wrap gap-2">
					<Button variant="outline" onclick={() => setRole(u.user_id, WorkspaceUserPatchRequest.role.VIEWER)} disabled={patchingId === u.user_id}>viewer</Button>
					<Button variant="outline" onclick={() => setRole(u.user_id, WorkspaceUserPatchRequest.role.EDITOR)} disabled={patchingId === u.user_id}>editor</Button>
					<Button variant="outline" onclick={() => setRole(u.user_id, WorkspaceUserPatchRequest.role.ADMIN)} disabled={patchingId === u.user_id}>admin</Button>
					<ConfirmButton
						label="Remove"
						title="Remove user from workspace?"
						description={`User '${u.email || `#${u.user_id}`}' will be removed from workspace.`}
						confirmLabel="Remove"
						pending={deletingId === u.user_id}
						onConfirm={() => removeUser(u.user_id)}
					/>
				</div>
			</li>
		{/each}
	</ul>

	<Pager
		{offset}
		pageSize={PAGE_SIZE}
		{total}
		{hasNext}
		onPrev={() => load(Math.max(offset - PAGE_SIZE, 0))}
		onNext={() => load(offset + PAGE_SIZE)}
	/>
{/if}
