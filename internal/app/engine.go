package app

import (
	"log"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/config"
	"github.com/pafthang/bms/internal/http/encoding"
	bmsmw "github.com/pafthang/bms/internal/http/middleware"
	"github.com/pafthang/bms/internal/http/routes"
	"github.com/pafthang/dbx"
)

// BuildOptions controls optional non-business route registration.
type BuildOptions struct {
	IncludeSystemRoutes bool
	IncludeHealthRoutes bool
}

// BuildEngine creates configured arc engine with global middleware and routes.
func BuildEngine(cfg config.Config, database *dbx.DB, opts BuildOptions) *arc.Engine {
	e := arc.New()
	e.Use(arc.Recovery(), arc.Logger(log.Default()))
	e.Use(bmsmw.CORS(cfg.AllowedOrigins))
	e.RegisterEncoder("application/problem+json", encoding.ProblemJSONEncoder{})
	e.SetOpenAPIServers([]map[string]any{
		{
			"url":         "http://localhost:8080/api/v1",
			"description": "dev",
		},
	})
	e.RegisterOpenAPISecurityScheme("BearerAuth", map[string]any{
		"type":         "http",
		"scheme":       "bearer",
		"bearerFormat": "JWT",
	})

	container := BuildContainer(cfg, database)
	routes.Register(e, routes.Dependencies{
		Auth:       container.Auth,
		Users:      container.Users,
		Settings:   container.Settings,
		Workspaces: container.Workspaces,
		Tags:       container.Tags,
		Bookmarks:  container.Bookmarks,
	})
	if opts.IncludeSystemRoutes && cfg.IncludeDocsRoute {
		e.RegisterSystemRoutes("/openapi.json", "/docs")
	}
	if opts.IncludeHealthRoutes {
		e.RegisterHealthRoutes()
	}
	return e
}
