package db

import (
	"context"

	"github.com/pafthang/orm"
)

func MigrateUp(ctx context.Context, database orm.DB, dir string) error {
	migrations, err := orm.LoadMigrationsDir(dir)
	if err != nil {
		return err
	}
	runner := orm.NewMigrationRunner(database)
	return runner.MigrateUp(ctx, migrations)
}
