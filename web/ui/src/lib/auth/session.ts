import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export type Session = {
	accessToken: string;
	refreshToken: string;
} | null;

const STORAGE_KEY = 'bms.session';

function loadSession(): Session {
	if (!browser) return null;
	const raw = localStorage.getItem(STORAGE_KEY);
	if (!raw) return null;
	try {
		const parsed = JSON.parse(raw);
		if (typeof parsed?.accessToken === 'string' && typeof parsed?.refreshToken === 'string') {
			return { accessToken: parsed.accessToken, refreshToken: parsed.refreshToken };
		}
	} catch {
		return null;
	}
	return null;
}

export function getSessionSnapshot(): Session {
	return loadSession();
}

export const session = writable<Session>(loadSession());

export function setSession(next: Session) {
	session.set(next);
	if (!browser) return;
	if (!next) {
		localStorage.removeItem(STORAGE_KEY);
		return;
	}
	localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
}
