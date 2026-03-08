<script lang="ts">
	import { toasts, removeToast } from '$lib/stores/toast';
</script>

<div class="pointer-events-none fixed right-4 top-4 z-50 flex w-[min(92vw,360px)] flex-col gap-2">
	{#each $toasts as t (t.id)}
		<div class="pointer-events-auto rounded-lg border border-border bg-card p-3 shadow-sm">
			<div class="flex items-start justify-between gap-3">
				<div>
					<p class={`text-sm font-medium ${t.type === 'error' ? 'text-destructive' : ''}`}>{t.title}</p>
					{#if t.message}
						<p class="mt-1 text-xs text-muted-foreground">{t.message}</p>
					{/if}
				</div>
				<button class="text-xs text-muted-foreground" onclick={() => removeToast(t.id)}>close</button>
			</div>
		</div>
	{/each}
</div>
