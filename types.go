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
	DOUBLE_CLICK IotClickType = "DOUBLE"
	LONG_PRESS   IotClickType = "LONG"
)

// RequestedResource is a resiurce used in API Gateway requests.
type RequestedResource string

// RequestRouter is used to distribute requests to suitable handler bases on used resource.
type RequestRouter struct {

	// Logger, to log errors or any kind of other information.
	logger log.Logger

	// CaptureHandler is used to process time tracking event requests.
	captureHandler Handler

	// ReportHandler will handle request for report generation.
	reportHandler Handler

	// Routes is a map of handlers asssigned to specific resources.
	routes map[RequestedResource]Handler
}

// CaptureRequestHandler process and persist captured request for time tracking records.
type CaptureRequestHandler struct {
	logger      log.Logger
	timeTracker timetracker.TimeTracker
}

// ReportGenerateRequestHandler will process request to generate and publish monthly time tracking reports.
type ReportGenerateRequestHandler struct {
	logger log.Logger
}

// TimeTrackingReport os a single captured time tracking event.
type TimeTrackingRecord struct {

	// DeviceId is an identifier of a device which captures a time tracking record.
	DeviceId string `json:"deviceid`

	// Type of a time tracking event.
	ClickType IotClickType `json:"clicktype"`

	// Timestamp is the point in time a time tracking event has occured.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// ReportGenerateRequest is used to triiger report generation for a specific year and month.
type ReportGenerateRequest struct {

	// Type of report which hould be generated. Atm monthly reports are supported, only.
	Type string `json:"type`

	// Year a monthly report should be generated for.
	Year int `json:"year`

	// Month a monthly report should be generated for.
	Month int `json:"month`
}
