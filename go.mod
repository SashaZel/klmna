module klmna/main

go 1.23.3

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	klmna/dbs v0.0.0-00010101000000-000000000000
)

replace klmna/dbs => ./src/dbs
