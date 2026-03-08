/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DataEnvelope_UserSettingsResponse } from '../models/DataEnvelope_UserSettingsResponse';
import type { PatchUserSettingsRequest } from '../models/PatchUserSettingsRequest';
import type { Problem } from '../models/Problem';
import type { PutUserSettingsRequest } from '../models/PutUserSettingsRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UserSettingsService {
    /**
     * @returns DataEnvelope_UserSettingsResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static usersMeSettingsGet(): CancelablePromise<DataEnvelope_UserSettingsResponse | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/users/me/settings',
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
            },
        });
    }
    /**
     * @param requestBody
     * @returns DataEnvelope_UserSettingsResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static usersMeSettingsPatch(
        requestBody: PatchUserSettingsRequest,
    ): CancelablePromise<DataEnvelope_UserSettingsResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/api/v1/users/me/settings',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @param requestBody
     * @returns DataEnvelope_UserSettingsResponse OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static usersMeSettingsPut(
        requestBody: PutUserSettingsRequest,
    ): CancelablePromise<DataEnvelope_UserSettingsResponse | Problem> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/users/me/settings',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                422: `Unprocessable Entity`,
            },
        });
    }
}
