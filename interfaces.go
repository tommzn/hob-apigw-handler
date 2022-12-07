package main

import (
	"github.com/aws/aws-lambda-go/events"
)

// Handler is used to process request forwarded by AWS API Gateway.
type Handler interface {

	// Process will handle capture requests for time tracking.
	Process(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
