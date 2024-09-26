package main

import (
	"log"

	"lambda-bc-opt/db"
)

func Main(args map[string]interface{}) map[string]interface{} {
	rdb := db.ConsRedisDB()
	log.Println("gooz2")

	// Set the number of goroutines you're going to wait for
	for i := 0; i < 100; i++ {
		rdb.Get("cnt")
	}

	return map[string]interface{}{
		"statusCode": 200,
		"body":       "salam",
	}

}
