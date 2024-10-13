package main

import (
	"io"
	"log"
	"net/http"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsMockRedisDB()

func getterHandler(w http.ResponseWriter, r *http.Request) {
	rdb.Get("cnt")
}

func main() {
	log.SetOutput(io.Discard)
	http.HandleFunc("/getterMock", getterHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
