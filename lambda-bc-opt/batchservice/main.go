package main

import (
	"fmt"
	"log/slog"
	"encoding/json"
	"os"

	"lambda-bc-opt/db"
	"lambda-bc-opt/utility"

	"github.com/valyala/fasthttp"
)

var rdb db.KeyValueStoreDB

// func getHandler (w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
//		return
//	}

//	defer r.Body.Close()
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}

//	var getOp db.GetOp
//	err = json.Unmarshal(body, &getOp)
//	if err != nil {
//		http.Error(w, "Invalid JSON", http.StatusBadRequest)
//		return
//	}

//	result, _ := rdb.Get(getOp.K)
//	slog.Debug(fmt.Sprintf("Received key: %s => %s\n", getOp.K, result))

//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(result))
// }
func getHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body()

	var getOp db.GetOp

	err := json.Unmarshal(body, &getOp)
	if err != nil {
		panic(err)
	}

	result, err := rdb.Get(getOp.K)
	if err != nil {
		panic(err)
	}
	slog.Debug(fmt.Sprintf("%s value in DB => %s\n", getOp.K, result))

	ctx.SetBodyString(result)
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

	fmt.Printf("Server listening onnn %s\n", address)
	if err := fasthttp.ListenAndServe(address, getHandler); err != nil {
		slog.Error(fmt.Sprintf("Error in ListenAndServe: %v", err))
	}
}
