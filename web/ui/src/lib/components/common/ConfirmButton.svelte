<script lang="ts">
	import { Button } from '$components/ui/button';

	let {
		label = 'Delete',
		title = 'Confirm action',
		description = 'Are you sure?',
		confirmLabel = 'Confirm',
		cancelLabel = 'Cancel',
		disabled = false,
		pending = false,
		onConfirm
	}: {
		label?: string;
		title?: string;
		description?: string;
		confirmLabel?: string;
		cancelLabel?: string;
		disabled?: boolean;
		pending?: boolean;
		onConfirm: () => void | Promise<void>;
	} = $props();

	let open = $state(false);
	let submitting = $state(false);

	async function confirm() {
		submitting = true;
		try {
			await onConfirm();
			open = false;
		} finally {
			submitting = false;
		}
	}
</script>

<Button variant="outline" onclick={() => (open = true)} disabled={disabled || pending}>{pending ? 'Processing...' : label}</Button>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" role="dialog" aria-modal="true">
		<div class="w-full max-w-sm rounded-lg border border-border bg-background p-4 shadow-lg">
			<h3 class="text-base font-semibold">{title}</h3>
			<p class="mt-2 text-sm text-muted-foreground">{description}</p>
			<div class="mt-4 flex justify-end gap-2">
				<Button variant="outline" onclick={() => (open = false)} disabled={submitting}>{cancelLabel}</Button>
				<Button onclick={confirm} disabled={submitting}>{submitting ? 'Please wait...' : confirmLabel}</Button>
			</div>
		</div>
	</div>
{/if}
