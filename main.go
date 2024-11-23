package main

import (
	"klmna/internal/api"
	"klmna/internal/db"
	"log"
	"net/http"
)

func main() {
	log.Print("server is starting")

	pgdb := db.StartDB()

	router := api.StartAPI(pgdb)

	defer pgdb.Close()

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("error from router %v\n", err)
	}
}
