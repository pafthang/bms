package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/pafthang/arc"
	bmsauth "github.com/pafthang/bms/internal/auth"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
)

type identityCtxKey struct{}

type Identity struct {
	UserID       int64
	Email        string
	IsSuperadmin bool
}

func IdentityFromContext(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityCtxKey{}).(Identity)
	return id, ok
}

func RequireAuth(svc *serviceauth.Service) arc.Middleware {
	return func(next arc.Handler) arc.Handler {
		return func(rc *arc.RequestContext) error {
			if svc == nil {
				return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "authentication is not configured"}
			}
			authz := strings.TrimSpace(rc.Request.Header.Get("Authorization"))
			if authz == "" {
				return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing bearer token"}
			}
			parts := strings.SplitN(authz, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "invalid authorization header"}
			}

			claims, err := svc.VerifyAccess(parts[1])
			if err != nil {
				if errors.Is(err, bmsauth.ErrExpiredToken) {
					return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_token_expired", Message: "token expired"}
				}
				return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_token_invalid", Message: "invalid token"}
			}
			id := Identity{UserID: claims.Subject, Email: claims.Email, IsSuperadmin: claims.IsSuperadmin}
			rc.Ctx = context.WithValue(rc.Ctx, identityCtxKey{}, id)
			rc.Request = rc.Request.WithContext(rc.Ctx)
			return next(rc)
		}
	}
}
