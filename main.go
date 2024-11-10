package main

import (
    "fmt"
    "net/http"
    dbs "klmna/dbs"
    "log"
)

func main() {
    
    db, err := dbs.CreateConnection()
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
    })
	fmt.Println("Server starts!")

    http.ListenAndServe(":80", nil)
}