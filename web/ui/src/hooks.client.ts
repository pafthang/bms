import type { Reroute } from '@sveltejs/kit';
import { loginPath, normalizeNextPath } from '$lib/auth/navigation';
import { getSessionSnapshot } from '$lib/auth/session';

const GUEST_ONLY_PATHS = new Set(['/login', '/register']);

export const reroute: Reroute = ({ url }) => {
	const { pathname, search } = url;
	const hasSession = Boolean(getSessionSnapshot());

	if (pathname.startsWith('/app') && !hasSession) {
		const next = `${pathname}${search}`;
		return loginPath(next);
	}

	if (GUEST_ONLY_PATHS.has(pathname) && hasSession) {
		return normalizeNextPath(url.searchParams.get('next'));
	}
};
