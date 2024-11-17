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

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Printf("error from router %v\n", err)
	}
}
