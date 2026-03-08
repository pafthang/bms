/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type WorkspaceUserCreateRequest = {
    email?: (string | null);
    role: WorkspaceUserCreateRequest.role;
    user_id?: (number | null);
    workspace_id: number;
};
export namespace WorkspaceUserCreateRequest {
    export enum role {
        ADMIN = 'admin',
        EDITOR = 'editor',
        VIEWER = 'viewer',
    }
}

