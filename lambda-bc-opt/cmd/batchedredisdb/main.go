package main

import (
	"lambda-bc-opt/db"
	"log"
	"runtime"
	"sync"
	"time"
)

func main() {
	batchedRedisDB := db.ConsBatchedRedisDB()
	// db.ConsBatchedRedisDB()

	// batchedRedisDB := db.ConsRedisDB()
	var wg sync.WaitGroup
	wg.Add(3)

	println(runtime.GOMAXPROCS(0))

	go func() {
		result := batchedRedisDB.Set("g1", "value1")
		log.Printf("result1 %e", result)
		wg.Done()
	}()

	go func() {
		result, _ := batchedRedisDB.Get("g2")
		log.Printf("result2 %s", result)
		wg.Done()
	}()

	go func() {
		log.Println("DONE")
		time.Sleep(20 * time.Second)
		wg.Done()
	}()

	wg.Wait()
}
