<script lang="ts">
	import { Button } from '$components/ui/button';

	let {
		offset,
		pageSize,
		total,
		hasNext,
		onPrev,
		onNext
	}: {
		offset: number;
		pageSize: number;
		total: number;
		hasNext: boolean;
		onPrev: () => void | Promise<void>;
		onNext: () => void | Promise<void>;
	} = $props();

	const from = $derived(total === 0 ? 0 : offset + 1);
	const to = $derived(Math.min(offset + pageSize, total));
</script>

<div class="mt-4 flex items-center justify-between">
	<p class="text-sm text-muted-foreground">Showing {from}-{to} of {total}</p>
	<div class="flex gap-2">
		<Button variant="outline" onclick={onPrev} disabled={offset === 0}>Previous</Button>
		<Button variant="outline" onclick={onNext} disabled={!hasNext}>Next</Button>
	</div>
</div>
