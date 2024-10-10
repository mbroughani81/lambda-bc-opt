package main

import (
	"log"
	"net/http"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsRedisDB("10.10.0.1", "6379")

func getterHandler(w http.ResponseWriter, r *http.Request) {
	rdb.Get("cnt")
}

func main() {
	http.HandleFunc("/getter", getterHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
