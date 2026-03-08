<script lang="ts">
	import { page } from '$app/state';
</script>

<main class="mx-auto w-full max-w-3xl p-6">
	<h1 class="text-2xl font-semibold">App error</h1>
	<p class="mt-2 text-muted-foreground">{page.error?.message ?? 'Unexpected application error'}</p>
	<p class="mt-2 text-sm text-muted-foreground">Status: {page.status}</p>
	<p class="mt-6 text-sm"><a class="underline" href="/app">Back to dashboard</a></p>
</main>
