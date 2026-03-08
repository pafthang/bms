package middleware

import (
	"net/http"
	"strings"

	"github.com/pafthang/arc"
)

// CORS applies a minimal CORS policy for UI clients.
func CORS(allowedOrigins []string) arc.Middleware {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		allowed[origin] = struct{}{}
	}

	return func(next arc.Handler) arc.Handler {
		return func(rc *arc.RequestContext) error {
			origin := strings.TrimSpace(rc.Request.Header.Get("Origin"))
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					h := rc.Writer.Header()
					h.Set("Access-Control-Allow-Origin", origin)
					h.Set("Vary", "Origin")
					h.Set("Access-Control-Allow-Credentials", "true")
					h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
					h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				}
			}

			if rc.Request.Method == http.MethodOptions {
				rc.Writer.WriteHeader(http.StatusNoContent)
				return nil
			}
			return next(rc)
		}
	}
}
