<script lang="ts">
	import { OpenAPI } from '@pafthang/bms-sdk';
	import { Button } from '$components/ui/button';
	import { onMount } from 'svelte';
	import { session } from '$lib/auth/session';

	let hasSession = false;

	onMount(() => {
		const unsub = session.subscribe((s) => {
			hasSession = Boolean(s);
		});
		return unsub;
	});
</script>

<main class="mx-auto flex min-h-screen w-full max-w-4xl flex-col justify-center gap-6 p-8">
	<div class="space-y-2">
		<p class="text-sm text-muted-foreground">BMS Frontend</p>
		<h1 class="text-4xl font-semibold tracking-tight">SvelteKit UI is ready</h1>
		<p class="max-w-2xl text-muted-foreground">
			Workspace setup is complete: Svelte 5 + SvelteKit + Tailwind + shadcn-svelte + typed SDK.
		</p>
	</div>

	<div class="rounded-2xl border border-border bg-card p-5">
		<p class="text-sm text-muted-foreground">API base URL</p>
		<p class="mt-1 font-mono text-sm">{OpenAPI.BASE}</p>
		<div class="mt-4 flex gap-3">
			{#if hasSession}
				<Button href="/app">Open App</Button>
			{:else}
				<Button href="/login">Sign In</Button>
				<Button variant="outline" href="/register">Create Account</Button>
			{/if}
			<Button variant="outline" href="http://localhost:8080/docs">Backend Docs</Button>
		</div>
	</div>
</main>
