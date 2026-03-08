/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BookmarkCreateRequest } from '../models/BookmarkCreateRequest';
import type { BookmarkPatchRequest } from '../models/BookmarkPatchRequest';
import type { DataEnvelope_BookmarkResponse } from '../models/DataEnvelope_BookmarkResponse';
import type { ListEnvelope_BookmarkResponse } from '../models/ListEnvelope_BookmarkResponse';
import type { Problem } from '../models/Problem';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class BookmarksService {
    /**
     * @param workspaceId
     * @param limit
     * @param offset
     * @param search
     * @param tag
     * @param tagId
     * @param isArchived
     * @param sort
     * @param order
     * @returns ListEnvelope_BookmarkResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static bookmarksList(
        workspaceId: number,
        limit?: number,
        offset?: number,
        search?: string,
        tag?: string,
        tagId?: (number | null),
        isArchived?: (boolean | null),
        sort?: string,
        order?: string,
    ): CancelablePromise<ListEnvelope_BookmarkResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/workspaces/{workspaceId}/bookmarks',
            path: {
                'workspaceId': workspaceId,
            },
            query: {
                'limit': limit,
                'offset': offset,
                'search': search,
                'tag': tag,
                'tag_id': tagId,
                'is_archived': isArchived,
                'sort': sort,
                'order': order,
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
     * @returns DataEnvelope_BookmarkResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static bookmarksCreate(
        workspaceId: number,
        requestBody: BookmarkCreateRequest,
    ): CancelablePromise<DataEnvelope_BookmarkResponse | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/workspaces/{workspaceId}/bookmarks',
            path: {
                'workspaceId': workspaceId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param workspaceId
     * @param bookmarkId
     * @returns any OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static bookmarksDelete(
        workspaceId: number,
        bookmarkId: number,
    ): CancelablePromise<any | Problem> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/workspaces/{workspaceId}/bookmarks/{bookmarkId}',
            path: {
                'workspaceId': workspaceId,
                'bookmarkId': bookmarkId,
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
     * @param bookmarkId
     * @returns DataEnvelope_BookmarkResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static bookmarksGet(
        workspaceId: number,
        bookmarkId: number,
    ): CancelablePromise<DataEnvelope_BookmarkResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/workspaces/{workspaceId}/bookmarks/{bookmarkId}',
            path: {
                'workspaceId': workspaceId,
                'bookmarkId': bookmarkId,
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
     * @param bookmarkId
     * @param requestBody
     * @returns DataEnvelope_BookmarkResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static bookmarksPatch(
        workspaceId: number,
        bookmarkId: number,
        requestBody: BookmarkPatchRequest,
    ): CancelablePromise<DataEnvelope_BookmarkResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/api/v1/workspaces/{workspaceId}/bookmarks/{bookmarkId}',
            path: {
                'workspaceId': workspaceId,
                'bookmarkId': bookmarkId,
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
