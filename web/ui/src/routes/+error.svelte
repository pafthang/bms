<script lang="ts">
	import { page } from '$app/state';
</script>

<main class="mx-auto flex min-h-screen w-full max-w-3xl flex-col justify-center p-6">
	<h1 class="text-3xl font-semibold">Something went wrong</h1>
	<p class="mt-2 text-muted-foreground">{page.error?.message ?? 'Unexpected error'}</p>
	<p class="mt-2 text-sm text-muted-foreground">Status: {page.status}</p>
	<p class="mt-6 text-sm"><a class="underline" href="/">Go home</a></p>
</main>
