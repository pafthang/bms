package app

import (
	bmsauth "github.com/pafthang/bms/internal/auth"
	"github.com/pafthang/bms/internal/config"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	servicebookmarks "github.com/pafthang/bms/internal/services/bookmarks"
	servicesettings "github.com/pafthang/bms/internal/services/settings"
	servicetags "github.com/pafthang/bms/internal/services/tags"
	serviceusers "github.com/pafthang/bms/internal/services/users"
	serviceworkspaces "github.com/pafthang/bms/internal/services/workspaces"
	"github.com/pafthang/dbx"
)

type Container struct {
	Auth       *serviceauth.Service
	Users      *serviceusers.Service
	Settings   *servicesettings.Service
	Workspaces *serviceworkspaces.Service
	Tags       *servicetags.Service
	Bookmarks  *servicebookmarks.Service
}

func BuildContainer(cfg config.Config, database *dbx.DB) *Container {
	jwtManager := bmsauth.NewJWTManager(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	tagsSvc := servicetags.NewService(database)
	return &Container{
		Auth:       serviceauth.NewService(database, jwtManager),
		Users:      serviceusers.NewService(database),
		Settings:   servicesettings.NewService(database),
		Workspaces: serviceworkspaces.NewService(database),
		Tags:       tagsSvc,
		Bookmarks:  servicebookmarks.NewService(database, tagsSvc),
	}
}
