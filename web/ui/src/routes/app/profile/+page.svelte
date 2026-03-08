<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';
	import { pushToast } from '$lib/stores/toast';

	type Me = {
		id: number;
		email: string;
		status: string;
		is_superadmin: boolean;
	};

	let me = $state<Me | null>(null);
	let email = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	async function loadProfile() {
		loading = true;
		error = '';
		try {
			const res = await withApi(() => api.users.usersMeGet());
			const data = (res as { data?: Me | null }).data ?? null;
			me = data;
			email = data?.email ?? '';
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load profile');
		} finally {
			loading = false;
		}
	}

	onMount(loadProfile);

	async function save(e: SubmitEvent) {
		e.preventDefault();
		if (!me) return;
		error = '';
		saving = true;
		try {
			const res = await withApi(() => api.users.usersMePatch({ email: email || null }));
			const data = (res as { data?: Me | null }).data ?? null;
			me = data;
			email = data?.email ?? '';
			pushToast({ type: 'success', title: 'Profile updated' });
		} catch (e) {
			error = userMessageFromError(e, 'Failed to update profile');
		} finally {
			saving = false;
		}
	}
</script>

<main class="mx-auto w-full max-w-3xl p-6">
	<h1 class="text-2xl font-semibold">Profile</h1>

	{#if loading}
		<StateBlock message="Loading..." />
	{:else if !me}
		<StateBlock kind="muted" message="Profile is unavailable." />
	{:else}
		<div class="mt-4 rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground">User ID</p>
			<p class="font-medium">{me.id}</p>
			<p class="mt-3 text-sm text-muted-foreground">Status</p>
			<p class="font-medium">{me.status}</p>
			<p class="mt-3 text-sm text-muted-foreground">Superadmin</p>
			<p class="font-medium">{me.is_superadmin ? 'yes' : 'no'}</p>
		</div>

		<form class="mt-4 space-y-4" onsubmit={save}>
			<div class="space-y-1">
				<label class="text-sm" for="email">Email</label>
				<input id="email" class="w-full rounded-md border border-border bg-background px-3 py-2" bind:value={email} type="email" required />
			</div>

			{#if error}<StateBlock kind="error" message={error} />{/if}

			<Button type="submit" disabled={saving}>{saving ? 'Saving...' : 'Save changes'}</Button>
		</form>
	{/if}
</main>
