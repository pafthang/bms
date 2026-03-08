/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DataEnvelope_AdminUserResponse } from '../models/DataEnvelope_AdminUserResponse';
import type { ListEnvelope_AdminUserResponse } from '../models/ListEnvelope_AdminUserResponse';
import type { Problem } from '../models/Problem';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AdminService {
    /**
     * @param limit
     * @param offset
     * @param search
     * @returns ListEnvelope_AdminUserResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static adminUsersList(
        limit?: number,
        offset?: number,
        search?: string,
    ): CancelablePromise<ListEnvelope_AdminUserResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/admin/users',
            query: {
                'limit': limit,
                'offset': offset,
                'search': search,
            },
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param userId
     * @returns DataEnvelope_AdminUserResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static adminUsersGet(
        userId: number,
    ): CancelablePromise<DataEnvelope_AdminUserResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/admin/users/{userId}',
            path: {
                'userId': userId,
            },
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
                404: `Not Found`,
            },
        });
    }
}
