package main

import (
    "fmt"
    "net/http"
    "log"
    "klmna/internal/dbs"
)

func main() {
    
    fmt.Println(dbs.Sum(2, 2))
    
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
