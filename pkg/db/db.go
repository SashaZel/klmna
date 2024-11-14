package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"

	"klmna/pkg/migrations"
)

func StartDB() *bun.DB {
	var connectionString string
	if os.Getenv("ENV") == "PROD" {
		connectionString = os.Getenv("DATABASE_URL")
	} else {
		connectionString = "postgres://klmna-user:pwdfrdby@127.0.0.1:5432/klmna-db?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(),
	))

	migrations.Init()

	ctx := context.Background()
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	migrator.Init(ctx)
	if err := migrator.Lock(ctx); err != nil {
		log.Fatal(err)
	}
	defer migrator.Unlock(ctx) //nolint:errcheck

	group, err := migrator.Migrate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if group.IsZero() {
		fmt.Printf("there are no new migrations to run (database is up to date)\n")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}

	return db
}
