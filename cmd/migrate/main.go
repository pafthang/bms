package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/pafthang/bms/internal/config"
	"github.com/pafthang/bms/internal/db"
	"github.com/pafthang/orm"
)

func main() {
	var (
		action string
		dir    string
		steps  int
	)
	flag.StringVar(&action, "action", "up", "migration action: up|down|status|validate")
	flag.StringVar(&dir, "dir", "migrations", "migrations directory")
	flag.IntVar(&steps, "steps", 1, "down steps count")
	flag.Parse()

	cfg := config.Load()
	database, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	migrationsDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("resolve migrations dir: %v", err)
	}
	migrations, err := orm.LoadMigrationsDir(migrationsDir)
	if err != nil {
		log.Fatalf("load migrations: %v", err)
	}

	runner := orm.NewMigrationRunner(database)
	ctx := context.Background()

	switch action {
	case "up":
		err = runner.MigrateUp(ctx, migrations)
	case "down":
		err = runner.MigrateDown(ctx, migrations, steps)
	case "validate":
		err = runner.Validate(migrations)
	case "status":
		var st *orm.MigrationStatus
		st, err = runner.Status(ctx, migrations)
		if err == nil {
			fmt.Printf("current=%d applied=%d pending=%d dirty=%v\n", st.CurrentVersion, len(st.Applied), len(st.Pending), st.Dirty)
		}
	default:
		log.Fatalf("unsupported action %q", action)
	}

	if err != nil {
		log.Fatalf("migration %s failed: %v", action, err)
	}
	log.Printf("migration %s completed", action)
}
