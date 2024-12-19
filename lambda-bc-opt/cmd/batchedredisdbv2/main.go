package main

import (
	"fmt"
	"lambda-bc-opt/db"
	"log/slog"
	"os"
)

func main() {
	opts := &slog.HandlerOptions{
		// Level: slog.LevelInfo,
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	n := 10
	counter := 0
	batchedRedisDB := db.ConsBatchedRedisDBV2("127.0.0.1", "8090")
	key := "cnt"
	for i := 0; i < n; i++ {
		result, err := batchedRedisDB.Get(key)
		if err != nil {
			slog.Warn(err.Error())
		}
		counter++
		slog.Debug(fmt.Sprintf("counter => %d", counter))
		slog.Debug(fmt.Sprintf("result => %s", result))
	}
}
