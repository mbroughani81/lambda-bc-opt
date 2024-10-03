package main

import (
	"fmt"
	"log"
	"sync"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDBV2("10.10.0.1:8080")

func Main(args map[string]interface{}) map[string]interface{} {
	n := 10
	key := "cnt"

	var wg sync.WaitGroup
	var result string

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			result, _ = rdb.Get(key)
			log.Printf("id = %d, resultt => %s", id, result)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("TAMAM")

	return map[string]interface{}{
		"statusCode": 200,
		"body":       fmt.Sprintf("key = %s, value = %s", key, result),
	}
}
