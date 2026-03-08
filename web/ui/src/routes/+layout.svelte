<script lang="ts">
	import { navigating } from '$app/stores';
	import ToastViewport from '$components/common/ToastViewport.svelte';
	import '../app.css';
	let { children } = $props();
</script>

{#if $navigating}
	<div class="fixed left-0 top-0 z-50 h-0.5 w-full animate-pulse bg-primary"></div>
{/if}

<ToastViewport />

{@render children()}
