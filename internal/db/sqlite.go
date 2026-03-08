package db

import (
	_ "modernc.org/sqlite"

	"github.com/pafthang/bms/internal/config"
	"github.com/pafthang/dbx"
)

func Open(cfg config.Config) (*dbx.DB, error) {
	return dbx.MustOpen("sqlite", cfg.DBDSN)
}
