package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"lambda-bc-opt/db"
)

type Op struct {
	opType string
	callback chan struct{}
}

const workerCount int = 1000
const bufferSize int = 1000000

var tasksChan chan Op = make(chan Op, bufferSize)
var rdbArray [workerCount]db.KeyValueStoreDB

func startWorkers() {
	// start work
	for i := 0; i < workerCount; i++ {
		f := func(workerId int) {
			for {
				select {
				case op := <-tasksChan: // a task is assigned
					slog.Debug(fmt.Sprintf("opType <%s> - workerId %d : Starting", op.opType, workerId))
					result, _ := rdbArray[workerId].Get("cnt")
					op.callback <- struct{}{}
					slog.Debug(fmt.Sprintf("opType <%s> - workerId %d : Ended - %s", op.opType, workerId, result))
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


	// 1: Start workers
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

	// 2: Create Worker's db connection
	for i := 0; i < workerCount; i++ {
		// rdbArray[i] = db.ConsRedisDB("localhost", "6379")
		// rdbArray[i] = db.ConsMockRedisDB()
		rdbArray[i] = db.ConsBatchedRedisDBV2("127.0.0.1", "8090")
	}

	// 3: Start endpoint
	http.HandleFunc("/locallambda", httpHandler)
	slog.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("ListenAndServe: ", err)
	}
}
