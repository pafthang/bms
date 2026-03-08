package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/pafthang/arc"
	bmsauth "github.com/pafthang/bms/internal/auth"
	"github.com/pafthang/bms/internal/http/dto"
	httpmw "github.com/pafthang/bms/internal/http/middleware"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	"github.com/pafthang/orm"
)

type authHandler struct {
	svc *serviceauth.Service
}

func RegisterAuth(api *arc.Group, svc *serviceauth.Service) {
	h := &authHandler{svc: svc}
	auth := api.Group("/auth").WithTags("Auth")

	arc.HandleGroup(auth, http.MethodPost, "/register", "auth_register", h.register,
		arc.WithRequestExamples(map[string]any{
			"default": map[string]any{
				"email":    "user@example.com",
				"password": "password123",
			},
		}),
		arc.WithResponseExamples(map[string]any{
			"default": map[string]any{
				"data": map[string]any{
					"user": map[string]any{
						"id":            1,
						"email":         "user@example.com",
						"is_superadmin": false,
						"status":        "active",
					},
					"tokens": map[string]any{
						"access_token":  "eyJhbGciOi...",
						"refresh_token": "eyJhbGciOi...",
					},
				},
			},
		}),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusConflict:            {Code: "conflict_error", Detail: "resource conflict"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(auth, http.MethodPost, "/login", "auth_login", h.login,
		arc.WithRequestExamples(map[string]any{
			"default": map[string]any{
				"email":    "user@example.com",
				"password": "password123",
			},
		}),
		arc.WithResponseExamples(map[string]any{
			"default": map[string]any{
				"data": map[string]any{
					"user": map[string]any{
						"id":            1,
						"email":         "user@example.com",
						"is_superadmin": false,
						"status":        "active",
					},
					"tokens": map[string]any{
						"access_token":  "eyJhbGciOi...",
						"refresh_token": "eyJhbGciOi...",
					},
				},
			},
		}),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_invalid_credentials", Detail: "invalid email or password"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(auth, http.MethodPost, "/refresh", "auth_refresh", h.refresh,
		arc.WithRequestExamples(map[string]any{
			"default": map[string]any{
				"refresh_token": "eyJhbGciOi...",
			},
		}),
		arc.WithResponseExamples(map[string]any{
			"default": map[string]any{
				"data": map[string]any{
					"tokens": map[string]any{
						"access_token":  "eyJhbGciOi...",
						"refresh_token": "eyJhbGciOi...",
					},
				},
			},
		}),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_token_invalid", Detail: "invalid token"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(auth.Group("", httpmw.RequireAuth(svc)), http.MethodGet, "/me", "auth_me", h.me,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusNotFound:     {Code: "not_found", Detail: "resource not found"},
		}),
	)
}

func (h *authHandler) register(ctx context.Context, in *dto.AuthRegisterRequest) (*arc.Response[dto.DataEnvelope[dto.AuthPayload]], error) {
	user, tokens, err := h.svc.Register(ctx, in.Email, in.Password)
	if err != nil {
		return nil, mapAuthError(err)
	}
	return arc.Created(dto.DataEnvelope[dto.AuthPayload]{
		Data: dto.AuthPayload{User: user, Tokens: tokens},
	}), nil
}

func (h *authHandler) login(ctx context.Context, in *dto.AuthLoginRequest) (*arc.Response[dto.DataEnvelope[dto.AuthPayload]], error) {
	user, tokens, err := h.svc.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, mapAuthError(err)
	}
	return arc.OK(dto.DataEnvelope[dto.AuthPayload]{
		Data: dto.AuthPayload{User: user, Tokens: tokens},
	}), nil
}

func (h *authHandler) refresh(ctx context.Context, in *dto.AuthRefreshRequest) (*arc.Response[dto.DataEnvelope[dto.RefreshPayload]], error) {
	tokens, err := h.svc.Refresh(ctx, in.RefreshToken)
	if err != nil {
		return nil, mapAuthError(err)
	}
	return arc.OK(dto.DataEnvelope[dto.RefreshPayload]{
		Data: dto.RefreshPayload{Tokens: tokens},
	}), nil
}

func (h *authHandler) me(ctx context.Context, in *struct{}) (*arc.Response[dto.DataEnvelope[*serviceauth.UserView]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	user, err := h.svc.Me(ctx, identity.UserID)
	if err != nil {
		return nil, mapAuthError(err)
	}
	return arc.OK(dto.DataEnvelope[*serviceauth.UserView]{Data: user}), nil
}

func mapAuthError(err error) error {
	if err == nil {
		return nil
	}
	if orm.HasCode(err, orm.CodeConflict) {
		return &arc.APIError{Status: http.StatusConflict, Code: "conflict_error", Message: "resource conflict"}
	}
	if orm.HasCode(err, orm.CodeNotFound) {
		return &arc.APIError{Status: http.StatusNotFound, Code: "not_found", Message: "resource not found"}
	}
	if errors.Is(err, bmsauth.ErrExpiredToken) {
		return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_token_expired", Message: "token expired"}
	}
	if errors.Is(err, bmsauth.ErrInvalidToken) {
		return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_token_invalid", Message: "invalid token"}
	}
	if errors.Is(err, serviceauth.ErrInvalidCredentials) {
		return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_invalid_credentials", Message: "invalid email or password"}
	}
	if errors.Is(err, serviceauth.ErrNotConfigured) {
		return &arc.APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "service is not configured"}
	}
	return err
}
