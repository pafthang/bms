import {
	AdminService,
	ApiError,
	AuthService,
	BookmarksService,
	OpenAPI,
	TagsService,
	UserSettingsService,
	UsersService,
	WorkspaceUsersService,
	WorkspacesService
} from '@pafthang/bms-sdk';
import { get } from 'svelte/store';
import { session, setSession } from '$lib/auth/session';
import { pushToast } from '$lib/stores/toast';

type ProblemLike = {
	status?: number;
	code?: string;
	title?: string;
	detail?: string;
	requestId?: string;
	request_id?: string;
};

const API_BASE = (import.meta.env.PUBLIC_BMS_API_BASE as string | undefined) ?? 'http://localhost:8080';

OpenAPI.BASE = API_BASE;
OpenAPI.TOKEN = async () => get(session)?.accessToken ?? '';
let refreshPromise: Promise<boolean> | null = null;

export class UIApiError extends Error {
	readonly status: number;
	readonly code: string;
	readonly title: string;
	readonly detail: string;
	readonly requestId?: string;

	constructor(problem: ProblemLike, fallbackStatus: number) {
		super(problem.detail || problem.title || 'Request failed');
		this.status = problem.status ?? fallbackStatus;
		this.code = problem.code ?? 'unknown_error';
		this.title = problem.title ?? 'Error';
		this.detail = problem.detail ?? '';
		this.requestId = problem.requestId ?? problem.request_id;
	}
}

const ERROR_CODE_MESSAGES: Record<string, string> = {
	auth_invalid_credentials: 'Invalid email or password.',
	auth_unauthorized: 'Your session is not valid. Please sign in again.',
	workspace_forbidden: 'You do not have permission for this workspace action.',
	workspace_not_found: 'Workspace not found.',
	workspace_membership_conflict: 'User membership conflict in this workspace.',
	workspace_last_admin_violation: 'This action is blocked: workspace must keep at least one admin.',
	bookmark_not_found: 'Bookmark not found.',
	tag_not_found: 'Tag not found.',
	validation_error: 'Validation failed. Check input fields.',
	conflict_error: 'Conflict: resource already exists or state changed.',
	internal_error: 'Internal server error. Try again later.'
};

export function userMessageFromError(err: unknown, fallback: string): string {
	if (err instanceof UIApiError) {
		return ERROR_CODE_MESSAGES[err.code] || err.detail || err.message || fallback;
	}
	if (err instanceof Error) return err.message;
	return fallback;
}

function extractTokenPair(payload: unknown): { accessToken: string; refreshToken: string } | null {
	const tokens = (payload as { data?: { tokens?: { access_token?: string; refresh_token?: string } } })?.data?.tokens;
	const access = tokens?.access_token ?? '';
	const refresh = tokens?.refresh_token ?? '';
	if (!access || !refresh) return null;
	return { accessToken: access, refreshToken: refresh };
}

async function refreshSessionOnce(): Promise<boolean> {
	if (refreshPromise) return refreshPromise;
	refreshPromise = (async () => {
		const current = get(session);
		if (!current?.refreshToken) {
			setSession(null);
			return false;
		}
		try {
			const res = await AuthService.authRefresh({ refresh_token: current.refreshToken });
			const next = extractTokenPair(res);
			if (!next) {
				setSession(null);
				return false;
			}
			setSession(next);
			return true;
		} catch {
			setSession(null);
			return false;
		} finally {
			refreshPromise = null;
		}
	})();
	return refreshPromise;
}

export async function withApi<T>(fn: () => Promise<T>): Promise<T> {
	try {
		return await fn();
	} catch (err) {
		if (err instanceof ApiError) {
			if (err.status === 401) {
				const refreshed = await refreshSessionOnce();
				if (refreshed) {
					try {
						return await fn();
					} catch (retryErr) {
						if (retryErr instanceof ApiError) {
							const retryProblem = (retryErr.body ?? {}) as ProblemLike;
							if (retryErr.status === 401) setSession(null);
							const mapped = new UIApiError(retryProblem, retryErr.status);
							pushToast({ type: 'error', title: mapped.title, message: mapped.detail || mapped.message });
							throw mapped;
						}
						throw retryErr;
					}
				}
				setSession(null);
			}
			const problem = (err.body ?? {}) as ProblemLike;
			const mapped = new UIApiError(problem, err.status);
			pushToast({ type: 'error', title: mapped.title, message: mapped.detail || mapped.message });
			throw mapped;
		}
		throw err;
	}
}

export const api = {
	auth: AuthService,
	admin: AdminService,
	users: UsersService,
	settings: UserSettingsService,
	workspaces: WorkspacesService,
	workspaceUsers: WorkspaceUsersService,
	bookmarks: BookmarksService,
	tags: TagsService
};
