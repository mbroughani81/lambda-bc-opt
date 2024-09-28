package main

import (
	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDB()

func Main(args map[string]interface{}) map[string]interface{} {
	n := 100
	cc := make(chan int, n)
	for i := 0; i < n; i++ {
		go func() {
			rdb.Get("cnt")
			cc <- 1
		}()
	}
	for i := 0; i < n; i++ {
		<-cc
	}

	return map[string]interface{}{
		"statusCode": 200,
		"body":       "salam",
	}
}
