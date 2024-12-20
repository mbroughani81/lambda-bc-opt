package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"log/slog"
	"net/http"

	"lambda-bc-opt/db"
	"lambda-bc-opt/utility"
)

var rdb db.KeyValueStoreDB

func getHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
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

	result, _ := rdb.Get(getOp.K)
	slog.Debug(fmt.Sprintf("Received key: %s => %s\n", getOp.K, result))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// DB
	redisHost := utility.GetEnv("REDIS_HOST", "127.0.0.1")
	redisPort := "6379"

	// API
	host := utility.GetEnv("APP_HOST", "127.0.0.1")
	port := utility.GetEnv("APP_PORT", "8090")
	address := fmt.Sprintf("%s:%s", host, port)

	// rdb = db.ConsRedisDB(redisHost, redisPort, 1)
	rdb = db.ConsBatchedRedisDB(redisHost, redisPort, 1)
	http.HandleFunc("/get", getHandler)

	fmt.Printf("Server listening onnn %s\n", address)
	err := http.ListenAndServe(address, nil)

	if err != nil {
		slog.Error("ListenAndServe: ", err)
	}
}
