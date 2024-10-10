package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"lambda-bc-opt/db"
)

var rdb db.KeyValueStoreDB = db.ConsMockRedisDB()

func init() {
	log.SetOutput(os.Stdout)
	log.Printf("thread : %d", runtime.GOMAXPROCS(-1))
}

func Main(args map[string]interface{}) map[string]interface{} {
	result, _ := rdb.Get("cnt")
	return map[string]interface{}{
		"statusCode": 200,
		"body":       fmt.Sprintf("Last result: %s", result),
	}
}
