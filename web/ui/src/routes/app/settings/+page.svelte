<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$components/ui/button';
	import StateBlock from '$components/common/StateBlock.svelte';
	import { api, userMessageFromError, withApi } from '$lib/api/client';
	import { pushToast } from '$lib/stores/toast';

	let theme = $state('');
	let locale = $state('');
	let timezone = $state('');
	let bookmarksPerPage = $state(20);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	onMount(async () => {
		try {
			const res = await withApi(() => api.settings.usersMeSettingsGet());
			const data = (res as { data?: { theme?: string | null; locale?: string | null; timezone?: string | null; bookmarks_per_page?: number | null } })?.data;
			theme = data?.theme ?? '';
			locale = data?.locale ?? '';
			timezone = data?.timezone ?? '';
			bookmarksPerPage = data?.bookmarks_per_page ?? 20;
		} catch (e) {
			error = userMessageFromError(e, 'Failed to load settings');
		} finally {
			loading = false;
		}
	});

	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		saving = true;
		try {
			await withApi(() =>
				api.settings.usersMeSettingsPut({
					theme: theme || null,
					locale: locale || null,
					timezone: timezone || null,
					bookmarks_per_page: bookmarksPerPage
				})
			);
			pushToast({ type: 'success', title: 'Settings saved' });
		} catch (e) {
			error = userMessageFromError(e, 'Save failed');
		} finally {
			saving = false;
		}
	}
</script>

<main class="mx-auto w-full max-w-3xl p-6">
	<h1 class="text-2xl font-semibold">Settings</h1>

	{#if loading}
		<StateBlock message="Loading..." />
	{:else}
		<form class="mt-4 space-y-4" onsubmit={submit}>
			<div class="space-y-1">
				<label class="text-sm" for="theme">Theme</label>
				<input id="theme" class="w-full rounded-md border border-border bg-background px-3 py-2" bind:value={theme} placeholder="light / dark" />
			</div>
			<div class="space-y-1">
				<label class="text-sm" for="locale">Locale</label>
				<input id="locale" class="w-full rounded-md border border-border bg-background px-3 py-2" bind:value={locale} placeholder="en-US" />
			</div>
			<div class="space-y-1">
				<label class="text-sm" for="timezone">Timezone</label>
				<input id="timezone" class="w-full rounded-md border border-border bg-background px-3 py-2" bind:value={timezone} placeholder="Europe/Moscow" />
			</div>
			<div class="space-y-1">
				<label class="text-sm" for="bpp">Bookmarks per page</label>
				<input id="bpp" class="w-full rounded-md border border-border bg-background px-3 py-2" bind:value={bookmarksPerPage} type="number" min={1} max={100} />
			</div>

			{#if error}<StateBlock kind="error" message={error} />{/if}

			<Button type="submit" disabled={saving}>{saving ? 'Saving...' : 'Save'}</Button>
		</form>
	{/if}
</main>
