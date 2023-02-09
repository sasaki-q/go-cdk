package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Body string `json:"body"`
}

func handleItemRequest(
	c context.Context,
	e events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Body: %s", e.Body),
	}, nil
}

func main() {
	lambda.Start(handleItemRequest)
}
