package main

import (
	"context"
	"lambda-bc-opt/db"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func test1(n int, sleep int) {
	var redisDB db.KeyValueStoreDB = db.ConsRedisDB("10.10.0.1", "6379")
	for {
		start := time.Now()

		for i := 1; i < n; i++ {
			key := "var_" + strconv.Itoa(i)
			redisDB.Get(key)
		}
		duration := time.Since(start)
		averageTimePerQuery := float64(duration.Microseconds()) / float64(n) / 1000

		log.Printf("duration : %v ", duration)
		log.Printf("average duration : %v ", averageTimePerQuery)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

func test2(n int, sleep int) {
	var redisDB *db.RedisDB = db.ConsRedisDB("10.10.0.1", "6379")
	rc := redisDB.Getrc()
	// Using pipe
	result := [10000]*redis.StringCmd{}

	for {
		start := time.Now()
		ctx := context.Background()
		pipe := rc.Pipeline()
		for i := 0; i < n; i++ {
			// key := "var_" + strconv.Itoa(i)
			key := "cnt"
			result[i] = pipe.Get(ctx, key)
		}
		pipe.Exec(ctx)
		log.Println(result[0])

		duration := time.Since(start)
		averageTimePerQuery := float64(duration.Microseconds()) / float64(n) / 1000

		log.Printf("duration : %v ", duration)
		log.Printf("average duration : %v ", averageTimePerQuery)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

func main() {
	var n int = 10000
	test2(n, 1000)
	// test1(n, 1000)
}
