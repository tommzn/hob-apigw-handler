package main

import (
	"github.com/aws/aws-lambda-go/events"
)

// SuccessfulResponse returns a response with status code 200.
func successfulResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{StatusCode: 200}
}

// SuccessfulResponse returns a response with status code 200.
func responseWithContent(content string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{StatusCode: statusCode, Body: content}
}

// ErrorResponse returns a response with status code 500 and passed error as body.
func errorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}
}

// errorResponseWithStatus returns a response with given status code and passed error as body.
func errorResponseWithStatus(err error, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       err.Error(),
	}
}
