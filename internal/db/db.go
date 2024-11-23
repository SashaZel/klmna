package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func StartDB() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file provided")
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresSslMode := os.Getenv("POSTGRES_SSL_MODE")

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		postgresUser,
		postgresPassword,
		postgresHost,
		postgresPort,
		postgresDB,
		postgresSslMode,
	)

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
