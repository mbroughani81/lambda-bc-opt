package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"lambda-bc-opt/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(rdb db.KeyValueStoreDB) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		cntVal, err := rdb.Get("cnt")
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "Internal Server Error",
				StatusCode: 500,
			}, nil
		}

		cnt, _ := strconv.Atoi(cntVal)
		cnt++

		err = rdb.Set("cnt", strconv.Itoa(cnt))
		if err != nil {
			log.Printf("Error updating 'cnt' in Redis: %v", err)
			return events.APIGatewayProxyResponse{
				Body:       "Internal Server Error",
				StatusCode: 500,
			}, nil
		}

		greeting := fmt.Sprintf("Hello! You are visitor number %d.\n", cnt)
		log.Printf("greeting => %s", greeting)
		return events.APIGatewayProxyResponse{
			Body:       greeting,
			StatusCode: 200,
		}, nil
	}
}


func main() {
	rdb := db.ConsMockRedisDB()
	// rdb := db.ConsRedisDB()
	// rdb := db.ConsBatchedRedisDB()
	handler := Handler(rdb)
	lambda.Start(handler)
}
