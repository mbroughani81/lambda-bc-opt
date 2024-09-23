package main

import (
	"lambda-bc-opt/db"
)

func main() {
	batchedRedisDB := db.ConsBatchedRedisDBV2("10.10.0.1:8080")
	batchedRedisDB.Get("goozal")
}