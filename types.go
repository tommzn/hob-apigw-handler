package main

import (
	"time"

	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

// IotClickType represents a click on an AWS IOT 1-Clickt type.
type IotClickType string

const (
	SINGLE_CLICK IotClickType = "SINGLE"
	DOUBLE_CLICK IotClickType = "DOuBLE"
	LONG_PRESS   IotClickType = "LONG_PRESS"
)

// APIGatewayRequestHandler process and persist captured request for time tracking records.
type APIGatewayRequestHandler struct {
	logger      log.Logger
	timeTracker timetracker.TimeTracker
}

// TimeTrackingReport os a single captured time tracking event.
type TimeTrackingRecord struct {

	// DeviceId is an identifier of a device which captures a time tracking record.
	DeviceId string `json:"deviceid`

	// Type of a time tracking event.
	ClickType IotClickType `json:"clicktyoe"`

	// Timestamp is the point in time a time tracking event has occured.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}
