package main

import (
	"klmna/pkg/api"
	"klmna/pkg/db"
	"log"
	"net/http"
)

func main() {
	log.Print("server is starting")

	pgdb, err := db.StartDB()
	if err != nil {
		log.Printf("error starting the db %v", err)
	}

	router := api.StartAPI(pgdb)

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Printf("error from router %v\n", err)
	}
}
