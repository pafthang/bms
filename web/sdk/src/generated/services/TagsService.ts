/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DataEnvelope_TagResponse } from '../models/DataEnvelope_TagResponse';
import type { ListEnvelope_TagResponse } from '../models/ListEnvelope_TagResponse';
import type { Problem } from '../models/Problem';
import type { TagCreateRequest } from '../models/TagCreateRequest';
import type { TagPatchRequest } from '../models/TagPatchRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class TagsService {
    /**
     * @param workspaceId
     * @param limit
     * @param offset
     * @returns ListEnvelope_TagResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static tagsList(
        workspaceId: number,
        limit?: number,
        offset?: number,
    ): CancelablePromise<ListEnvelope_TagResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/workspaces/{workspaceId}/tags',
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
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param requestBody
     * @returns DataEnvelope_TagResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static tagsCreate(
        workspaceId: number,
        requestBody: TagCreateRequest,
    ): CancelablePromise<DataEnvelope_TagResponse | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/workspaces/{workspaceId}/tags',
            path: {
                'workspaceId': workspaceId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                409: `Conflict`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param tagId
     * @returns any OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static tagsDelete(
        workspaceId: number,
        tagId: number,
    ): CancelablePromise<any | Problem> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/workspaces/{workspaceId}/tags/{tagId}',
            path: {
                'workspaceId': workspaceId,
                'tagId': tagId,
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
     * @param tagId
     * @param requestBody
     * @returns DataEnvelope_TagResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static tagsPatch(
        workspaceId: number,
        tagId: number,
        requestBody: TagPatchRequest,
    ): CancelablePromise<DataEnvelope_TagResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/api/v1/workspaces/{workspaceId}/tags/{tagId}',
            path: {
                'workspaceId': workspaceId,
                'tagId': tagId,
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
