package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	db := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", "10.10.0.1", "6379"),
		DB:   0,
	})

	ctx := context.Background()
	n := 10000
	sleep := 2000

	for {
		start := time.Now()
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			go func() {
				wg.Add(1)

				db.Get(ctx, "cnt").Result()

				fmt.Printf("wg: %v\n", wg.)
				wg.Done()
				// log.Printf("result => %s", result)
			}()
		}

		duration := time.Since(start)
		averageTimePerQuery := float64(duration.Microseconds()) / float64(n) / 1000

		log.Printf("duration : %v ", duration)
		log.Printf("average duration : %v ", averageTimePerQuery)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

}
