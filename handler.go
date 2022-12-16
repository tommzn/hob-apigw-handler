package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

// NewRequestHandler create a handler to process API Gateway requests.
func newRequestHandler(timeTracker timetracker.TimeTracker, logger log.Logger) *APIGatewayRequestHandler {
	return &APIGatewayRequestHandler{
		logger:      logger,
		timeTracker: timeTracker,
	}
}

// Process will process time tracking request and persist it using time tracker repository.
func (handler *APIGatewayRequestHandler) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	defer handler.logger.Flush()

	timeTrackingRecord, err := toTimeTrackingRecord(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Unable to decode request body",
		}, err
	}

	handler.logger.Debugf("TimeTrackingRecord: %+v", timeTrackingRecord)
	handler.logger.Statusf("Receive capture request (%s) from %s at %s", timeTrackingRecord.ClickType, timeTrackingRecord.DeviceId, timeTrackingRecord.Timestamp)

	recordType := toTimeTrackingRecordType(timeTrackingRecord.ClickType)
	if timeTrackingRecord.Timestamp == nil {
		err = handler.timeTracker.Capture(timeTrackingRecord.DeviceId, recordType)
	} else {
		err = handler.timeTracker.Captured(timeTrackingRecord.DeviceId, recordType, *timeTrackingRecord.Timestamp)
	}

	response := events.APIGatewayProxyResponse{StatusCode: 200}
	if err != nil {
		handler.logger.Error("Unable to capture time tracking recoed, reason: ", err)
		response.StatusCode = 500
		response.Body = err.Error()
	}
	return response, err
}

// ToTimeTrackingRecordType converts a AWS IOT click type to a time tracking record type.
func toTimeTrackingRecordType(clickType IotClickType) timetracker.RecordType {
	switch clickType {
	case SINGLE_CLICK:
		return timetracker.WORKDAY
	case DOUBLE_CLICK:
		return timetracker.ILLNESS
	case LONG_PRESS:
		return timetracker.VACATION
	default:
		return timetracker.WORKDAY
	}
}

func toTimeTrackingRecord(requestBody string) (TimeTrackingRecord, error) {
	var timeTrackingRecord TimeTrackingRecord
	err := json.Unmarshal([]byte(requestBody), &timeTrackingRecord)
	return timeTrackingRecord, err
}
