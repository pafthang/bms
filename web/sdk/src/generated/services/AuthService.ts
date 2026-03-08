/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AuthLoginRequest } from '../models/AuthLoginRequest';
import type { AuthRefreshRequest } from '../models/AuthRefreshRequest';
import type { AuthRegisterRequest } from '../models/AuthRegisterRequest';
import type { DataEnvelope_AuthPayload } from '../models/DataEnvelope_AuthPayload';
import type { DataEnvelope_RefreshPayload } from '../models/DataEnvelope_RefreshPayload';
import type { DataEnvelope_UserView } from '../models/DataEnvelope_UserView';
import type { Problem } from '../models/Problem';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AuthService {
    /**
     * @param requestBody
     * @returns DataEnvelope_AuthPayload OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static authLogin(
        requestBody: AuthLoginRequest,
    ): CancelablePromise<DataEnvelope_AuthPayload | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/auth/login',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Unauthorized`,
                422: `Unprocessable Entity`,
            },
        });
    }
    /**
     * @returns DataEnvelope_UserView OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static authMe(): CancelablePromise<DataEnvelope_UserView | Problem> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/auth/me',
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
            },
        });
    }
    /**
     * @param requestBody
     * @returns DataEnvelope_RefreshPayload OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static authRefresh(
        requestBody: AuthRefreshRequest,
    ): CancelablePromise<DataEnvelope_RefreshPayload | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/auth/refresh',
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
     * @returns DataEnvelope_AuthPayload OK
     * @returns Problem Error
     * @throws ApiError
     */
    public static authRegister(
        requestBody: AuthRegisterRequest,
    ): CancelablePromise<DataEnvelope_AuthPayload | Problem> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/auth/register',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                409: `Conflict`,
                422: `Unprocessable Entity`,
            },
        });
    }
}
