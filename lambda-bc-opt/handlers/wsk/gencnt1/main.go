package main

import (
	"fmt"
	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDBV2("10.10.0.1:8080")

func Main(args map[string]interface{}) map[string]interface{} {
	// n := 100
	// cc := make(chan int, n)
	// for i := 0; i < n; i++ {
	//	go func() {
	//		rdb.Get("cnt")
	//		cc <- 1
	//	}()
	// }
	// for i := 0; i < n; i++ {
	//	<-cc
	// }
	result, _ := rdb.Get("cnt")
	body := fmt.Sprintf("gooz result %s", result)

	return map[string]interface{}{
		"statusCode": 200,
		"body":       body,
	}
}
