<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		title,
		description,
		children,
		actions
	}: {
		title: string;
		description?: string;
		children?: Snippet;
		actions?: Snippet;
	} = $props();
</script>

<section class="rounded-xl border border-border bg-card p-4">
	<div class="mb-3 flex items-start justify-between gap-3">
		<div>
			<h3 class="text-base font-semibold">{title}</h3>
			{#if description}
				<p class="mt-1 text-sm text-muted-foreground">{description}</p>
			{/if}
		</div>
		<div>
			{@render actions?.()}
		</div>
	</div>
	<div class="space-y-3">
		{@render children?.()}
	</div>
</section>
