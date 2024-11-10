package dbs

import (
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func CreateConnection() (*sqlx.DB, error) {
	connStr := "user=klmna-user password=pwdfrdby dbname=klmna-db sslmode=disable"
    connection, err := sqlx.Open("postgres", connStr)
    if err != nil {
       log.Fatal(err)
    }

    err = connection.Ping()
    if err != nil {
        log.Fatal(err)
    }
    
    return connection, nil
}