import { api, withApi } from '$lib/api/client';
import { setSession } from '$lib/auth/session';

function extractTokens(payload: unknown): { accessToken: string; refreshToken: string } {
	const data = (payload as { data?: { tokens?: { access_token?: string; refresh_token?: string } } })?.data;
	const access = data?.tokens?.access_token ?? '';
	const refresh = data?.tokens?.refresh_token ?? '';
	if (!access || !refresh) {
		throw new Error('token pair is missing in response');
	}
	return { accessToken: access, refreshToken: refresh };
}

export async function login(email: string, password: string) {
	const res = await withApi(() => api.auth.authLogin({ email, password }));
	setSession(extractTokens(res));
	return res;
}

export async function register(email: string, password: string) {
	const res = await withApi(() => api.auth.authRegister({ email, password }));
	setSession(extractTokens(res));
	return res;
}

export function logout() {
	setSession(null);
}
