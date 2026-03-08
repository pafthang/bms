import { writable } from 'svelte/store';

export type ToastType = 'info' | 'success' | 'error';

export type Toast = {
	id: number;
	type: ToastType;
	title: string;
	message?: string;
	durationMs: number;
};

const DEFAULT_DURATION = 3500;
let seq = 1;

export const toasts = writable<Toast[]>([]);

export function pushToast(input: {
	type?: ToastType;
	title: string;
	message?: string;
	durationMs?: number;
}) {
	const toast: Toast = {
		id: seq++,
		type: input.type ?? 'info',
		title: input.title,
		message: input.message,
		durationMs: input.durationMs ?? DEFAULT_DURATION
	};
	toasts.update((prev) => [...prev, toast]);
	setTimeout(() => removeToast(toast.id), toast.durationMs);
}

export function removeToast(id: number) {
	toasts.update((prev) => prev.filter((t) => t.id !== id));
}
