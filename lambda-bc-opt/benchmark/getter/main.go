package main

import (
	"log"
	"net/http"
	"sync"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDB()

// var rdb db.KeyValueStoreDB = db.ConsRedisDB()

func getterHandler(w http.ResponseWriter, r *http.Request) {
	// Call the getter function
	n := 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			log.Println("gooz1")
			result, _ := rdb.Get("cnt")
			log.Printf("This is result => %s", result)
			wg.Done()
		}()
	}
	wg.Wait()
	// Send a response back to the client
	// fmt.Fprintf(w, "Getter function executed successfully")
}

func main() {
	// log.SetOutput(io.Discard)
	http.HandleFunc("/getter", getterHandler)

	// Start the HTTP server on port 8080
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
