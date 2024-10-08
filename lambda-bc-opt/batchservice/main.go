package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lambda-bc-opt/db"
	"lambda-bc-opt/utility"
	"log"
	"net/http"
)

func getHandler(rdb db.AKeyValueStoreDB) func(http.ResponseWriter, *http.Request) {
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

		ch := make(chan string)
		rdb.AGet(getOp.K, ch)
		result := <-ch
		log.Printf("Received key: %s => %s\n", getOp.K, result)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

func main() {
	// DB
	// log.SetOutput(io.Discard)
	// rdb := db.ConsBatchedRedisDB()

	redisHost := utility.GetEnv("REDIS_HOST", "10.10.0.1")
	redisPort := "6379"
	rdb := db.ConsBatchedRedisDB(redisHost, redisPort)
	// rdb := db.ConsMockRedisDB()

	// API
	host := utility.GetEnv("APP_HOST", "127.0.0.1")
	port := utility.GetEnv("APP_PORT", "8080")
	address := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("/get", getHandler(rdb))
	fmt.Printf("Server listening onnn %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
