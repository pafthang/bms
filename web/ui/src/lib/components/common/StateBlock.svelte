<script lang="ts">
	let {
		kind = 'info',
		message
	}: {
		kind?: 'info' | 'error' | 'muted';
		message: string;
	} = $props();

	const cls = $derived(
		kind === 'error'
			? 'text-destructive'
			: kind === 'muted'
				? 'text-muted-foreground'
				: ''
	);
</script>

<p class={`mt-4 text-sm ${cls}`}>{message}</p>
