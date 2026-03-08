const FALLBACK_PATH = '/app';

export function normalizeNextPath(input: string | null | undefined): string {
	if (!input) return FALLBACK_PATH;
	if (!input.startsWith('/')) return FALLBACK_PATH;
	if (input.startsWith('//')) return FALLBACK_PATH;
	if (input.startsWith('/login') || input.startsWith('/register')) return FALLBACK_PATH;
	return input;
}

export function loginPath(next?: string): string {
	const target = normalizeNextPath(next);
	return `/login?next=${encodeURIComponent(target)}`;
}
