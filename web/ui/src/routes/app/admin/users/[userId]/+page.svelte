<script lang="ts">
	import { onMount } from 'svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';

	let { data } = $props();
	const userId = $derived(data.userId as number);

	type AdminUser = {
		id: number;
		email: string;
		status: string;
		is_superadmin: boolean;
		created_at: string;
		updated_at: string;
	};

	let user = $state<AdminUser | null>(null);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.admin.adminUsersGet(userId));
			user = (res as { data?: AdminUser | null }).data ?? null;
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load user');
		} finally {
			loading = false;
		}
	});
</script>

<main class="mx-auto w-full max-w-3xl p-6">
	<h1 class="text-2xl font-semibold">Admin User Details</h1>
	<p class="mt-1 text-sm text-muted-foreground">User ID: {userId}</p>

	{#if loading}
		<p class="mt-4">Loading...</p>
	{:else if error}
		<p class="mt-4 text-sm text-destructive">{error}</p>
	{:else if !user}
		<p class="mt-4 text-muted-foreground">User not found.</p>
	{:else}
		<div class="mt-4 rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground">Email</p>
			<p class="font-medium">{user.email}</p>
			<p class="mt-3 text-sm text-muted-foreground">Status</p>
			<p class="font-medium">{user.status}</p>
			<p class="mt-3 text-sm text-muted-foreground">Superadmin</p>
			<p class="font-medium">{user.is_superadmin ? 'yes' : 'no'}</p>
			<p class="mt-3 text-sm text-muted-foreground">Created</p>
			<p class="font-medium">{new Date(user.created_at).toLocaleString()}</p>
			<p class="mt-3 text-sm text-muted-foreground">Updated</p>
			<p class="font-medium">{new Date(user.updated_at).toLocaleString()}</p>
		</div>
	{/if}

	<p class="mt-6 text-sm">
		<a class="underline" href="/app/admin/users">Back to users</a>
	</p>
</main>
