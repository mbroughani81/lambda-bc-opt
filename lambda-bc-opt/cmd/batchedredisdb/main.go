package main

import (
	"lambda-bc-opt/db"
	"log"
	"runtime"
	"sync"
)

func main() {
	batchedRedisDB := db.ConsBatchedRedisDB()
	var wg sync.WaitGroup
	wg.Add(2)

	println(runtime.GOMAXPROCS(0))

	go func() {
		batchedRedisDB.Set("g1", "value1")
		result, _ := batchedRedisDB.Get("g1")
		log.Println("#1 value1 => %s", result)
		wg.Done()
	}()

	go func() {
		result, _ := batchedRedisDB.Get("g2")
		log.Println("#2 => %s", result)
		batchedRedisDB.Set("g2", "value2")
		result, _ = batchedRedisDB.Get("g2")
		log.Println("#2 value2 => %s", result)
		wg.Done()
	}()
	wg.Wait()
}
