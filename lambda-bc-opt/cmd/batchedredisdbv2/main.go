package main

import (
	"lambda-bc-opt/db"
	"log"
	"sync"
)

func main() {
	// batchedRedisDB := db.ConsBatchedRedisDBV2("10.10.0.1:8080")
	// key := "cnt"
	// value, _ := batchedRedisDB.Get(key)
	// log.Printf("key: %s, value: %s", key, value)

	n := 7
	counter := 0
	batchedRedisDB := db.ConsBatchedRedisDBV2("10.10.0.1:8080")
	key := "cnt"
	var wg sync.WaitGroup
	wg.Add(n)
	result := "???"
	for i := 0; i < n; i++ {
		go func() {
			result, _ = batchedRedisDB.Get(key)
			counter++
			log.Printf("counter => %d", counter)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println(result)
}
