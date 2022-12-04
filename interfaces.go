package main

import "context"

// Handler is used to process request forwarded by AWS API Gateway.
type Handler interface {

	// Process will handle capture requests for time tracking.
	Process(context.Context, TimeTrackingRecord) error
}
