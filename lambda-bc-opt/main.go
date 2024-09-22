package main

import (
	"fmt"
	"log"
	"strconv"

	"lambda-bc-opt/db"
)

func Main(args map[string]interface{}) map[string]interface{} {
	rdb := db.ConsRedisDB()
	cntVal, err := rdb.Get("cnt")

	if err != nil {
		return map[string]interface{}{
			"statusCode": 500,
			"body":       "Internal server error",
		}
	}
	cnt, _ := strconv.Atoi(cntVal)
	cnt++

	// Create the greeting message
	greeting := fmt.Sprintf("Hello! cnt is %d.\n", cnt)
	log.Printf("greeting => %s", greeting)

	// Return the response in OpenWhisk format
	return map[string]interface{}{
		"statusCode": 200,
		"body":       greeting,
	}
}
