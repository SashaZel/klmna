package db

import (
	"context"
	// "fmt"
	"log"
	"os"

	// "database/sql"
	// "github.com/uptrace/bun"
	// "github.com/uptrace/bun/dialect/pgdialect"
	// "github.com/uptrace/bun/driver/pgdriver"
	// "github.com/uptrace/bun/extra/bundebug"
	// "github.com/uptrace/bun/migrate"

	// "klmna/internal/migrations"

	// "sync"
	_ "github.com/jackc/pgx/v5/stdlib"
	// "github.com/jackc/pgx/v5/pgxpool"
	"database/sql"
)

// type postgres struct {
// 	db *pgxpool.Pool
// }

// var (
// 	pgInstance *postgres
// 	pgOnce     sync.Once
// )

// func NewPG(ctx context.Context, connString string) (*postgres, error) {
// 	pgOnce.Do(func() {
// 		db, err := pgxpool.New(ctx, connString)
// 		if err != nil {
// 			log.Fatal("unable to create connection pool: %w", err)
// 		}
// 		pgInstance = &postgres{db}
// 	})

// 	return pgInstance, nil
// }

// func (pg *postgres) Ping(ctx context.Context) error {
// 	return pg.db.Ping(ctx)
// }

// func (pg *postgres) Close() {
// 	pg.db.Close()
// }

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

	// defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	// ctx, stop := context.WithCancel(context.Background())
	// defer stop()


	if err := pool.PingContext(context.Background()); err != nil {
		log.Fatal("error ping db %#w", err)
	}

	// pg, err := NewPG(context.Background(), connectionString)
	// if err != nil {
	// 	log.Fatal("error create db %#w", err)
	// }
	// err = pg.Ping(context.Background());
	// if err != nil {
	// 	log.Fatal("error ping db %#w", err)
	// }

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	// conn, err := pgx.Connect(context.Background(), connectionString)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(context.Background())

	// var name string
	// var weight int64
	// err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(name, weight)

	// sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))

	// db := bun.NewDB(sqldb, pgdialect.New())

	// db.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithEnabled(false),
	// 	bundebug.FromEnv(),
	// ))

	// migrations.Init()

	// ctx := context.Background()
	// migrator := migrate.NewMigrator(db, migrations.Migrations)
	// migrator.Init(ctx)
	// if err := migrator.Lock(ctx); err != nil {
	// 	log.Fatal(err)
	// }
	// defer migrator.Unlock(ctx) //nolint:errcheck

	// group, err := migrator.Migrate(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if group.IsZero() {
	// 	fmt.Printf("there are no new migrations to run (database is up to date)\n")
	// } else {
	// 	fmt.Printf("migrated to %s\n", group)
	// }

	return pool
}
