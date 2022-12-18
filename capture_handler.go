package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

// NewRequestHandler create a handler to process API Gateway requests.
func newCaptureRequestHandler(timeTracker timetracker.TimeTracker, logger log.Logger) *CaptureRequestHandler {
	return &CaptureRequestHandler{
		logger:      logger,
		timeTracker: timeTracker,
	}
}

// Process will process time tracking request and persist it using time tracker repository.
func (handler *CaptureRequestHandler) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	timeTrackingRecord, err := toTimeTrackingRecord(request.Body)
	if err != nil {
		handler.logger.Error(err)
		return errorResponse(err), err
	}

	handler.logger.Debugf("TimeTrackingRecord: %+v", timeTrackingRecord)
	handler.logger.Statusf("Receive capture request (%s) from %s at %s", timeTrackingRecord.ClickType, timeTrackingRecord.DeviceId, timeTrackingRecord.Timestamp)

	recordType := toTimeTrackingRecordType(timeTrackingRecord.ClickType)
	if timeTrackingRecord.Timestamp == nil {
		err = handler.timeTracker.Capture(timeTrackingRecord.DeviceId, recordType)
	} else {
		err = handler.timeTracker.Captured(timeTrackingRecord.DeviceId, recordType, *timeTrackingRecord.Timestamp)
	}

	if err != nil {
		handler.logger.Error("Unable to capture time tracking recoed, reason: ", err)
		return errorResponse(err), err
	}
	return successfulResponse(), nil
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

// ToTimeTrackingRecord try to convert passed request body to a time tracking record.
func toTimeTrackingRecord(requestBody string) (TimeTrackingRecord, error) {
	var timeTrackingRecord TimeTrackingRecord
	err := json.Unmarshal([]byte(requestBody), &timeTrackingRecord)
	return timeTrackingRecord, err
}
