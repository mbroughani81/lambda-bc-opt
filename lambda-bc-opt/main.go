package main

import (
	"fmt"
	"log"
	"strconv"

	"lambda-bc-opt/db"
)

func Main(args map[string]interface{}) map[string]interface{} {
	rdb := db.ConsMockRedisDB()
	cntVal, err := rdb.Get("cnt")

	if err != nil {
		return map[string]interface{}{
			"statusCode": 500,
			"body":       "Internal server error",
		}
	}
	cnt, _ := strconv.Atoi(cntVal)
	cnt++

	err = rdb.Set("cnt", strconv.Itoa(cnt))
	if err != nil {
		log.Printf("Error updating 'cnt' in Redis: %v", err)
		return map[string]interface{}{
			"statusCode": 500,
			"body":       "Internal Server Error",
		}
	}

	// Create the greeting message
	greeting := fmt.Sprintf("Hello! You are visitor number %d.\n", cnt)
	log.Printf("greeting => %s", greeting)

	// Return the response in OpenWhisk format
	return map[string]interface{}{
		"statusCode": 200,
		"body":       greeting,
	}
}
