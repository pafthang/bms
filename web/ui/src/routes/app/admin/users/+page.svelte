<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import Pager from '$components/common/Pager.svelte';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';

	type AdminUser = {
		id: number;
		email: string;
		status: string;
		is_superadmin: boolean;
		created_at: string;
		updated_at: string;
	};

	const PAGE_SIZE = 20;
	let items = $state<AdminUser[]>([]);
	let total = $state(0);
	let offset = $state(0);
	let hasNext = $state(false);
	let search = $state('');
	let loading = $state(true);
	let error = $state('');

	async function load(nextOffset = offset) {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.admin.adminUsersList(PAGE_SIZE, nextOffset, search || undefined));
			const envelope = (res as { data?: AdminUser[]; meta?: { total?: number; offset?: number; has_next?: boolean } }) ?? {};
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

	async function applySearch(e: SubmitEvent) {
		e.preventDefault();
		await load(0);
	}
</script>

<main class="mx-auto w-full max-w-6xl p-6">
	<h1 class="text-2xl font-semibold">Admin Users</h1>
	<p class="mt-1 text-sm text-muted-foreground">Superadmin area</p>

	<form class="mt-4 grid gap-2 sm:grid-cols-[1fr_auto_auto]" onsubmit={applySearch}>
		<input class="rounded-md border border-border bg-background px-3 py-2" placeholder="Search by email" bind:value={search} />
		<Button type="submit" variant="outline">Search</Button>
		<Button type="button" variant="outline" onclick={() => { search = ''; void load(0); }}>Reset</Button>
	</form>

	{#if error}<StateBlock kind="error" message={error} />{/if}

	{#if loading}
		<StateBlock message="Loading..." />
	{:else if items.length === 0}
		<StateBlock kind="muted" message="No users found." />
	{:else}
		<ul class="mt-4 space-y-2">
			{#each items as u}
				<li class="rounded-lg border border-border p-3">
					<div class="flex items-start justify-between gap-3">
						<div>
							<a class="font-medium underline" href={`/app/admin/users/${u.id}`}>{u.email}</a>
							<p class="text-sm text-muted-foreground">id: {u.id} | status: {u.status}</p>
						</div>
						{#if u.is_superadmin}
							<span class="rounded border border-border px-2 py-1 text-xs">superadmin</span>
						{/if}
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
</main>
