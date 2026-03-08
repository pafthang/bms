/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type WorkspaceUserPatchRequest = {
    role: WorkspaceUserPatchRequest.role;
    user_id: number;
    workspace_id: number;
};
export namespace WorkspaceUserPatchRequest {
    export enum role {
        ADMIN = 'admin',
        EDITOR = 'editor',
        VIEWER = 'viewer',
    }
}

