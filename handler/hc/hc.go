package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Body string `json:"body"`
}

func healthCheck(
	c context.Context,
) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Status: %s", os.Getenv("Status")),
	}, nil
}

func main() {
	lambda.Start(healthCheck)
}
