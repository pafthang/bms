<script lang="ts">
	import { onMount } from 'svelte';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { withApi, api, userMessageFromError } from '$lib/api/client';
	let { data } = $props();

	const workspaceId = $derived(data.workspaceId as number);
	let name = $state('');
	let description = $state('');
	let error = $state('');

	onMount(async () => {
		try {
			const res = await withApi(() => api.workspaces.workspacesGet(workspaceId));
			const w = (res as { data?: { name?: string; description?: string | null } }).data;
			name = w?.name ?? '';
			description = w?.description ?? '';
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load workspace');
		}
	});
</script>

{#if error}
	<StateBlock kind="error" message={error} />
{:else}
	<h1 class="text-2xl font-semibold">{name || `Workspace #${workspaceId}`}</h1>
	{#if description}
		<p class="mt-2 text-muted-foreground">{description}</p>
	{/if}
{/if}
