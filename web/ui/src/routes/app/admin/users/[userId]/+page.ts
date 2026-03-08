import { error } from '@sveltejs/kit';

export const ssr = false;

export function load({ params }: { params: { userId: string } }) {
	const userId = Number(params.userId);
	if (!Number.isInteger(userId) || userId <= 0) {
		throw error(400, 'invalid user id');
	}
	return { userId };
}
