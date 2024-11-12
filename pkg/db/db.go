package db

import (
	"log"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func StartDB() (*pg.DB, error) {
	var (
		connectionString string
		opts             *pg.Options
		err              error
	)

	if os.Getenv("ENV") == "PROD" {
		connectionString = os.Getenv("DATABASE_URL")
	} else {
		connectionString = "postgres://klmna-user:pwdfrdby@127.0.0.1:5432/klmna-db?sslmode=disable"
	}

	opts, err = pg.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	collection := migrations.NewCollection()
	err = collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}

	_, _, err = collection.Run(db, "init")
	if err != nil {
		return nil, err
	}

	oldVersion, newVersion, err := collection.Run(db, "up")
	if err != nil {
		return nil, err
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	return db, err
}
