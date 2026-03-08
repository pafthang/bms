<script lang="ts">
	let { data, children } = $props();
	const workspaceId = $derived(data.workspaceId as number);
</script>

<section class="mx-auto w-full max-w-6xl p-6">
	<header class="mb-4">
		<p class="text-sm text-muted-foreground">Workspace #{workspaceId}</p>
		<nav class="mt-2 flex gap-3 text-sm">
			<a class="rounded-md border border-border px-3 py-1" href={`/app/workspaces/${workspaceId}`}>Overview</a>
			<a class="rounded-md border border-border px-3 py-1" href={`/app/workspaces/${workspaceId}/bookmarks`}>Bookmarks</a>
			<a class="rounded-md border border-border px-3 py-1" href={`/app/workspaces/${workspaceId}/tags`}>Tags</a>
			<a class="rounded-md border border-border px-3 py-1" href={`/app/workspaces/${workspaceId}/users`}>Users</a>
		</nav>
	</header>

	{@render children()}
</section>
