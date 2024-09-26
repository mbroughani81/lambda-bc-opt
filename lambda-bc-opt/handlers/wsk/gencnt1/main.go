package main

import (
	"lambda-bc-opt/db"
	"log"
)

func Main(args map[string]interface{}) map[string]interface{} {
	rdb := db.ConsRedisDB()
	log.Println("gooz1")

	// Set the number of goroutines you're going to wait for
	n := 100
	cc := make(chan int, n)
	for i := 0; i < n; i++ {
		go func() {
			rdb.Get("cnt")
			cc <- 1
		}()
	}
	for i := 0; i < n; i++ {
		_ = <-cc
	}

	return map[string]interface{}{
		"statusCode": 200,
		"body":       "salam",
	}
}
