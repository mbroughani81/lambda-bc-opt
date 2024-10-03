package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lambda-bc-opt/db"
	"log"
	"net/http"
)

func getHandler(rdb db.KeyValueStoreDB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var getOp db.GetOp
		err = json.Unmarshal(body, &getOp)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		log.Printf("Received key: %s\n", getOp.K)
		result, _ := rdb.Get(getOp.K)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

func main() {
	// DB
	// log.SetOutput(io.Discard)
	rdb := db.ConsBatchedRedisDB()

	// API
	http.HandleFunc("/get", getHandler(rdb))
	fmt.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe("10.10.0.1:8080", nil))
}
