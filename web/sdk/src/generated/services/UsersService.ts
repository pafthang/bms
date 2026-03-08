/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DataEnvelope_UserMeResponse } from '../models/DataEnvelope_UserMeResponse';
import type { PatchUserMeRequest } from '../models/PatchUserMeRequest';
import type { Problem } from '../models/Problem';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UsersService {
    /**
     * @returns DataEnvelope_UserMeResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static usersMeGet(): CancelablePromise<DataEnvelope_UserMeResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/users/me',
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
            },
        });
    }
    /**
     * @param requestBody
     * @returns DataEnvelope_UserMeResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static usersMePatch(
        requestBody: PatchUserMeRequest,
    ): CancelablePromise<DataEnvelope_UserMeResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/api/v1/users/me',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                422: `Unprocessable Entity`,
            },
        });
    }
}
