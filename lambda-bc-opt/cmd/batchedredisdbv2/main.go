package main

import (
	"lambda-bc-opt/db"
	"log"
)

func main() {
	batchedRedisDB := db.ConsBatchedRedisDBV2("10.10.0.1:8080")
	key := "cnt"
	value, _ := batchedRedisDB.Get(key)
	log.Printf("key: %s, value: %s", key, value)
}
