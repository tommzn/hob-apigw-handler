package main

import (
	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
)

// NewReportGenerateRequestHandler returna handler to maintina, add and delete, time tracking records.
func newTimeTrackingRecordHandler(logger log.Logger) *TimeTrackingRecordHandler {
	return &TimeTrackingRecordHandler{
		logger: logger,
	}
}

// Process will generate and publish time tracking report for passed year/month.
func (handler *TimeTrackingRecordHandler) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	handler.logger.Statusf("[TimeTrackingManager] Request received: %s %s", request.HTTPMethod, request.Path)
	handler.logger.Statusf("[TimeTrackingManager] Query: %+v", request.QueryStringParameters)
	handler.logger.Statusf("[TimeTrackingManager] PathParams: %+v", request.PathParameters)

	return successfulResponse(), nil
}
