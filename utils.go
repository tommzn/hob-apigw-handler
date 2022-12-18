package main 

import (
	"github.com/aws/aws-lambda-go/events"
)

// SuccessfulResponse returns a response with status code 200.
func successfulResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{StatusCode: 200}
}

// ErrorResponse returns a response with status code 500 abd passed error as body.
func errorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body: err.Error(),
	}
}