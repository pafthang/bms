/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DataEnvelope_WorkspaceMembershipResponse } from '../models/DataEnvelope_WorkspaceMembershipResponse';
import type { ListEnvelope_WorkspaceMembershipResponse } from '../models/ListEnvelope_WorkspaceMembershipResponse';
import type { Problem } from '../models/Problem';
import type { WorkspaceUserCreateRequest } from '../models/WorkspaceUserCreateRequest';
import type { WorkspaceUserPatchRequest } from '../models/WorkspaceUserPatchRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class WorkspaceUsersService {
    /**
     * @param workspaceId
     * @param limit
     * @param offset
     * @returns ListEnvelope_WorkspaceMembershipResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspaceUsersList(
        workspaceId: number,
        limit?: number,
        offset?: number,
    ): CancelablePromise<ListEnvelope_WorkspaceMembershipResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/workspaces/{workspaceId}/users',
            path: {
                'workspaceId': workspaceId,
            },
            query: {
                'limit': limit,
                'offset': offset,
            },
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                404: `Not Found`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param requestBody
     * @returns DataEnvelope_WorkspaceMembershipResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspaceUsersAdd(
        workspaceId: number,
        requestBody: WorkspaceUserCreateRequest,
    ): CancelablePromise<DataEnvelope_WorkspaceMembershipResponse | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/workspaces/{workspaceId}/users',
            path: {
                'workspaceId': workspaceId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                404: `Not Found`,
                409: `Conflict`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param userId
     * @returns any OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspaceUsersDelete(
        workspaceId: number,
        userId: number,
    ): CancelablePromise<any | Problem> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/workspaces/{workspaceId}/users/{userId}',
            path: {
                'workspaceId': workspaceId,
                'userId': userId,
            },
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                404: `Not Found`,
                409: `Conflict`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param userId
     * @param requestBody
     * @returns DataEnvelope_WorkspaceMembershipResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspaceUsersPatch(
        workspaceId: number,
        userId: number,
        requestBody: WorkspaceUserPatchRequest,
    ): CancelablePromise<DataEnvelope_WorkspaceMembershipResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/api/v1/workspaces/{workspaceId}/users/{userId}',
            path: {
                'workspaceId': workspaceId,
                'userId': userId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                404: `Not Found`,
                409: `Conflict`,
                422: `Unprocessable Entity`,
            },
        });
    }
}
