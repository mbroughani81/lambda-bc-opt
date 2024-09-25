package main

import (
	"log"
	"sync"

	"lambda-bc-opt/db"
)

func Main(args map[string]interface{}) map[string]interface{} {
	rdb := db.ConsRedisDB()
	log.Println("gooz1")

	var wg sync.WaitGroup

	// Set the number of goroutines you're going to wait for
	n := 10
	wg.Add(10)
	for i := 0; i < n; i++ {
		go func() {
			rdb.Get("cnt")
			wg.Done()
		}()
	}
	wg.Wait()

	return map[string]interface{}{
		"statusCode": 200,
		"body":       "salam",
	}
}
