package routes

import (
	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/http/handlers"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	servicebookmarks "github.com/pafthang/bms/internal/services/bookmarks"
	servicesettings "github.com/pafthang/bms/internal/services/settings"
	servicetags "github.com/pafthang/bms/internal/services/tags"
	serviceusers "github.com/pafthang/bms/internal/services/users"
	serviceworkspaces "github.com/pafthang/bms/internal/services/workspaces"
)

const apiPrefix = "/api/v1"

type Dependencies struct {
	Auth       *serviceauth.Service
	Users      *serviceusers.Service
	Settings   *servicesettings.Service
	Workspaces *serviceworkspaces.Service
	Tags       *servicetags.Service
	Bookmarks  *servicebookmarks.Service
}

// Register mounts all business routes.
func Register(e *arc.Engine, deps Dependencies) {
	api := e.Group(apiPrefix)
	handlers.RegisterAuth(api, deps.Auth)
	handlers.RegisterUsers(api, deps.Users, deps.Settings, deps.Auth)
	handlers.RegisterWorkspaces(api, deps.Workspaces, deps.Auth)
	handlers.RegisterTags(api, deps.Tags, deps.Auth)
	handlers.RegisterBookmarks(api, deps.Bookmarks, deps.Auth)
}
