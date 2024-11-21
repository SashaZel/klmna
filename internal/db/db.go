package db

import (
	"context"
	"database/sql"
	"embed"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func StartDB() *sql.DB {
	var connectionString string
	if os.Getenv("ENV") == "PROD" {
		connectionString = os.Getenv("DATABASE_URL")
	} else {
		connectionString = "postgres://klmna-user:pwdfrdby@127.0.0.1:5432/klmna-db?sslmode=disable"
	}

	pool, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("error open db %#w", err)
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	if err := pool.PingContext(context.Background()); err != nil {
		log.Fatal("error ping db %#w", err)
	}

	// migrations

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(pool, "migrations"); err != nil {
		panic(err)
	}

	return pool
}
