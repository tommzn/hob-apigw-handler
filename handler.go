package main

import (
	"context"

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
func (handler *APIGatewayRequestHandler) Process(ctx context.Context, timeTrackingRecord TimeTrackingRecord) error {

	defer handler.logger.Flush()

	handler.logger.WithContext(ctx)
	handler.logger.Debugf("TimeTrackingRecord: %+v", timeTrackingRecord)
	handler.logger.Statusf("Receive capture request (%s) from %s at %s", timeTrackingRecord.ClickType, timeTrackingRecord.DeviceId, timeTrackingRecord.Timestamp)

	var err error
	recordType := toTimeTrackingRecordType(timeTrackingRecord.ClickType)
	if timeTrackingRecord.Timestamp == nil {
		err = handler.timeTracker.Capture(timeTrackingRecord.DeviceId, recordType)
	} else {
		err = handler.timeTracker.Captured(timeTrackingRecord.DeviceId, recordType, *timeTrackingRecord.Timestamp)
	}
	if err != nil {
		handler.logger.Error("Unable to capture time tracking recoed, reason: ", err)
	}
	return err
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
