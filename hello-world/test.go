package main

import (
	"context"

	"hello-world/handlers/insertrow"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	req := events.APIGatewayProxyRequest{
		Body: "Test Request",
	}
	ctx := context.Background()
	resp, _ := insertrow.Handler(ctx, req)
	println(resp.Body)
}