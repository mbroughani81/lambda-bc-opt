package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDBV2("10.10.0.1", "8080")

func init() {
	log.SetOutput(os.Stdout)
	log.Printf("thread : %d", runtime.GOMAXPROCS(-1))
}

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for key := range jobs {
		result, err := rdb.Get(key)
		if err != nil {
			log.Printf("Worker %d: error => %v", id, err)
			results <- fmt.Sprintf("Worker %d: error => %v", id, err)
			continue
		}

		log.Printf("Worker %d: result => %s", id, result)
		results <- result
	}
}

func Main(args map[string]interface{}) map[string]interface{} {
	key := "cnt"
	numJobs := 10    // Number of rdb.Get operations
	numWorkers := 10 // Number of workers
	var wg sync.WaitGroup

	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 0; j < numJobs; j++ {
		jobs <- key
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var lastResult string
	for res := range results {
		lastResult = res
	}

	log.Println("All workers completed.")
	return map[string]interface{}{
		"statusCode": 200,
		"body":       fmt.Sprintf("Last result: %s", lastResult),
	}
}
