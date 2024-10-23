package main

import (
	"fmt"
	"lambda-bc-opt/db"
	"log/slog"
	"net/http"
	"os"
)

const workerCount int = 1000

type Op struct {
	opType string
	callback chan struct{}
}

const bufferSize int = 1000000
var tasksChan chan Op = make(chan Op, bufferSize)
var rdb db.KeyValueStoreDB

func startWorkers() {
	rdb = db.ConsRedisDB("localhost", "6379")

	for i := 0; i < workerCount; i++ {
		f := func(goroutineId int) {
			for {
				select {
				case op := <-tasksChan: // a task is assigned
					slog.Debug(fmt.Sprintf("opType <%s> - goroutineId %d : Starting", op.opType, goroutineId))
					result, _ := rdb.Get("cnt")
					op.callback <- struct{}{}
					slog.Debug(fmt.Sprintf("opType <%s> - goroutineId %d : Ended - %s", op.opType, goroutineId, result))
				}
				slog.Debug("recurse")
			}
		}
		go f(i)
	}
}


func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Starting Benchmark")


	go func() {
		slog.Info("Running tasks: Starting")
		startWorkers()
	} ()

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		cb := make(chan struct{})
		tasksChan <- Op{
			opType: "Get",
			callback: cb,
		}
		<-cb
	}

	// rdb = db.ConsMockRedisDB()
	http.HandleFunc("/locallambda", httpHandler)
	slog.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("ListenAndServe: ", err)
	}
}
