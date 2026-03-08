<script lang="ts">
	import { onMount } from 'svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';

	let { children } = $props();
	let loading = $state(true);
	let allowed = $state(false);
	let error = $state('');

	onMount(async () => {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.users.usersMeGet());
			const me = (res as { data?: { is_superadmin?: boolean } | null })?.data;
			allowed = Boolean(me?.is_superadmin);
		} catch (e) {
			error = userMessageFromError(e, 'Failed to check permissions');
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<main class="mx-auto w-full max-w-6xl p-6">
		<p>Checking admin access...</p>
	</main>
{:else if error}
	<main class="mx-auto w-full max-w-6xl p-6">
		<p class="text-sm text-destructive">{error}</p>
	</main>
{:else if !allowed}
	<main class="mx-auto w-full max-w-6xl p-6">
		<h1 class="text-2xl font-semibold">Access denied</h1>
		<p class="mt-2 text-sm text-muted-foreground">This section is available only for superadmin users.</p>
		<p class="mt-4 text-sm"><a class="underline" href="/app">Back to dashboard</a></p>
	</main>
{:else}
	{@render children()}
{/if}
