<script lang="ts">
	import { page } from '$app/state';
</script>

<main class="max-w-3xl p-6">
	<h1 class="text-2xl font-semibold">Workspace section error</h1>
	<p class="mt-2 text-muted-foreground">{page.error?.message ?? 'Unexpected workspace error'}</p>
	<p class="mt-2 text-sm text-muted-foreground">Status: {page.status}</p>
	<p class="mt-4 text-sm">
		<a class="underline" href="/app/workspaces">Back to workspaces</a>
	</p>
</main>
