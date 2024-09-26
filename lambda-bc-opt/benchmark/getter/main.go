package main

import (
	"fmt"
	"log"
	"net/http"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDB()

// var rdb db.KeyValueStoreDB = db.ConsRedisDB()

func getterHandler(w http.ResponseWriter, r *http.Request) {
	// Call the getter function
	log.Println("gooz1")

	n := 100
	cc := make(chan int, n)
	for i := 0; i < n; i++ {
		go func() {
			result, _ := rdb.Get("cnt")
			log.Printf("result => %s", result)
			cc <- 1
		}()
	}
	for i := 0; i < n; i++ {
		_ = <-cc
	}

	// Send a response back to the client
	fmt.Fprintf(w, "Getter function executed successfully")
}

func main() {
	http.HandleFunc("/getter", getterHandler)

	// Start the HTTP server on port 8080
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
