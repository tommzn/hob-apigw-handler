package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/protobuf/proto"
)

// Handler is used to process request forwarded by AWS API Gateway.
type Handler interface {

	// Process will handle capture requests for time tracking.
	Process(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

// Publisher is used to send messages to one or multiple queues.
type Publisher interface {

	// Send will publish passed message to given queues.
	Send(message proto.Message) error
}
