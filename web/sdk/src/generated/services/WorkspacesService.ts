/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DataEnvelope_WorkspaceCreateResponse } from '../models/DataEnvelope_WorkspaceCreateResponse';
import type { DataEnvelope_WorkspaceResponse } from '../models/DataEnvelope_WorkspaceResponse';
import type { ListEnvelope_WorkspaceResponse } from '../models/ListEnvelope_WorkspaceResponse';
import type { Problem } from '../models/Problem';
import type { WorkspaceCreateRequest } from '../models/WorkspaceCreateRequest';
import type { WorkspacePatchByIDRequest } from '../models/WorkspacePatchByIDRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class WorkspacesService {
    /**
     * @param limit
     * @param offset
     * @param search
     * @param all
     * @returns ListEnvelope_WorkspaceResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspacesList(
        limit?: number,
        offset?: number,
        search?: string,
        all?: boolean,
    ): CancelablePromise<ListEnvelope_WorkspaceResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/workspaces',
            query: {
                'limit': limit,
                'offset': offset,
                'search': search,
                'all': all,
            },
            errors: {
                401: `Unauthorized`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param requestBody
     * @returns DataEnvelope_WorkspaceCreateResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspacesCreate(
        requestBody: WorkspaceCreateRequest,
    ): CancelablePromise<DataEnvelope_WorkspaceCreateResponse | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/workspaces',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param workspaceId
     * @returns any OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspacesDelete(
        workspaceId: number,
    ): CancelablePromise<any | Problem> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/workspaces/{workspaceId}',
            path: {
                'workspaceId': workspaceId,
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
     * @returns DataEnvelope_WorkspaceResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspacesGet(
        workspaceId: number,
    ): CancelablePromise<DataEnvelope_WorkspaceResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/workspaces/{workspaceId}',
            path: {
                'workspaceId': workspaceId,
            },
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param requestBody
     * @returns DataEnvelope_WorkspaceResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static workspacesPatch(
        workspaceId: number,
        requestBody: WorkspacePatchByIDRequest,
    ): CancelablePromise<DataEnvelope_WorkspaceResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/api/v1/workspaces/{workspaceId}',
            path: {
                'workspaceId': workspaceId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                404: `Not Found`,
                422: `Unprocessable Entity`,
            },
        });
    }
}
