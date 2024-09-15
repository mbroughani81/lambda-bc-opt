package main

import (
	"context"
	"fmt"
	"log"
	// "strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redis/go-redis/v9"
)

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Change this to your Redis server address
		DB:   0,                // Default DB number
	})
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// rdb := initRedis()

	// // Get the current value of "cnt" from Redis
	// cntVal, err := rdb.Get(ctx, "cnt").Result()
	// if err == redis.Nil {
	// 	// Key does not exist, initialize with 0
	// 	cntVal = "0"
	// } else if err != nil {
	// 	log.Printf("Error fetching 'cnt' from Redis: %v", err)
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       "Internal Server Error",
	// 		StatusCode: 500,
	// 	}, nil
	// }

	// // Convert the current count to an integer, increment it
	// cnt, _ := strconv.Atoi(cntVal)
	// cnt++


	// // Update the "cnt" key in Redis
	// err = rdb.Set(ctx, "cnt", cnt, 0).Err()
	// if err != nil {
	// 	log.Printf("Error updating 'cnt' in Redis: %v", err)
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       "Internal Server Error",
	// 		StatusCode: 500,
	// 	}, nil
	// }

	greeting := fmt.Sprintf("Hello! You are visitor number %d.\n", 1)
	log.Printf("greeting => %s", greeting)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
