package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"os"

	"github.com/redis/go-redis/v9"
)

const workerCount int = 1000
const invCount int = 100000

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Starting Benchmark")
	db := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", "localhost", "6379"),
		DB:   0,
		PoolSize: 1,
	})
	slog.Info("DB connected")

	// each of workers listern for the taskChan
	var tasksChan chan string = make(chan string, invCount)
	var done chan struct{} = make(chan struct{})

	// creating the task
	for i := 0; i < invCount; i++ {
		tasksChan <- "Get"
	}

	slog.Info("Tasks are ready in tasksChan")
	time.Sleep(5 * time.Second)

	// creating the workers (each worker is a goroutine)
	slog.Info("Running tasks: Starting")
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(invCount) // Done when goroutine finished the work

	ctx := context.Background()

	for i := 0; i < workerCount; i++ {
		f := func(goroutineId int) {
			for {
				slog.Debug("recurse")
				select {
				case <-done:
					return
				case task := <-tasksChan: // a task is assigned
					slog.Debug(fmt.Sprintf("task <%s> - goroutineId %d : Starting", task, goroutineId))
					result := db.Get(ctx, "cnt")
					slog.Debug(fmt.Sprintf("task <%s> - goroutineId %d : Done - result : %s", task, goroutineId, result))
					wg.Done()
				}
			}
		}
		go f(i)
	}
	wg.Wait()
	// done <- struct{}{}

	slog.Info("Running tasks: Done")
	duration := time.Since(start)
	slog.Info(fmt.Sprintf("Duration => %v", duration))
	// log.SetOutput(io.Discard)
	// http.HandleFunc("/getterNaive", getterHandler)

	// log.Println("Starting server on :8080")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	// }
}
