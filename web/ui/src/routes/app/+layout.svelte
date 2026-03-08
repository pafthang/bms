<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import { api, withApi } from '$lib/api/client';
	import { loginPath } from '$lib/auth/navigation';
	import { logout } from '$lib/auth/actions';
	import { session } from '$lib/auth/session';

	let { children } = $props();
	let ready = $state(false);
	let isSuperadmin = $state(false);

	onMount(() => {
		const unsub = session.subscribe((s) => {
			if (!s) {
				const next = `${window.location.pathname}${window.location.search}`;
				goto(loginPath(next), { replaceState: true });
				return;
			}
			ready = true;
		});
		return unsub;
	});

	onMount(async () => {
		try {
			const res = await withApi(() => api.users.usersMeGet());
			const me = (res as { data?: { is_superadmin?: boolean } | null })?.data;
			isSuperadmin = Boolean(me?.is_superadmin);
		} catch {
			isSuperadmin = false;
		}
	});

	async function onLogout() {
		logout();
		await goto('/login');
	}
</script>

{#if ready}
	<div class="min-h-screen bg-background">
		<header class="border-b border-border">
			<div class="mx-auto flex h-14 w-full max-w-6xl items-center justify-between px-4">
				<nav class="flex items-center gap-4 text-sm">
					<a href="/app">Home</a>
					<a href="/app/workspaces">Workspaces</a>
					{#if isSuperadmin}
						<a href="/app/admin/users">Admin</a>
					{/if}
					<a href="/app/profile">Profile</a>
					<a href="/app/settings">Settings</a>
				</nav>
				<Button variant="outline" onclick={onLogout}>Logout</Button>
			</div>
		</header>
		{@render children()}
	</div>
{:else}
	<p class="p-6">Checking session...</p>
{/if}
