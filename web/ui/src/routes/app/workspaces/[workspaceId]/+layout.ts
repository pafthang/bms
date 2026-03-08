import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ params }) => {
	return {
		workspaceId: Number(params.workspaceId)
	};
};
