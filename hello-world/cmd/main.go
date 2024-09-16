package main

import (
	"hello-world/db"
	"log"
	"sync"
)

func main() {
	batchedRedisDB := db.ConsBatchedRedisDB()
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {		
		result, _ := batchedRedisDB.Get("g1")
		log.Println("result-1 => %s", result)
		wg.Done()
	}()

	go func() {		
		result, _ := batchedRedisDB.Get("g2")
		log.Println("result-2 => %s", result)
		wg.Done()
	}()

	go func() {		
		result, _ := batchedRedisDB.Get("g3")
		log.Println("result-3 => %s", result)
		wg.Done()
	}()

	wg.Wait()
}