<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import { login } from '$lib/auth/actions';
	import { normalizeNextPath } from '$lib/auth/navigation';
	import { session } from '$lib/auth/session';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;
	let next = '/app';

	onMount(() => {
		next = normalizeNextPath(new URLSearchParams(window.location.search).get('next'));
		const unsub = session.subscribe((s) => {
			if (s) goto(next, { replaceState: true });
		});
		return unsub;
	});

	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		loading = true;
		try {
			await login(email, password);
			await goto(next, { replaceState: true });
		} catch (err) {
			error = err instanceof Error ? err.message : 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<main class="mx-auto flex min-h-screen w-full max-w-md flex-col justify-center p-6">
	<h1 class="text-3xl font-semibold tracking-tight">Sign in</h1>
	<p class="mt-2 text-sm text-muted-foreground">Use your BMS credentials.</p>

	<form class="mt-6 space-y-4" on:submit={submit}>
		<div class="space-y-1">
			<label for="email" class="text-sm">Email</label>
			<input id="email" class="w-full rounded-md border border-border bg-background px-3 py-2" type="email" bind:value={email} required />
		</div>
		<div class="space-y-1">
			<label for="password" class="text-sm">Password</label>
			<input id="password" class="w-full rounded-md border border-border bg-background px-3 py-2" type="password" bind:value={password} required minlength={8} />
		</div>

		{#if error}
			<p class="text-sm text-destructive">{error}</p>
		{/if}

		<Button type="submit" class="w-full" disabled={loading}>{loading ? 'Signing in...' : 'Sign in'}</Button>
	</form>

	<p class="mt-4 text-sm text-muted-foreground">
		No account? <a class="underline" href="/register">Create one</a>
	</p>
</main>
