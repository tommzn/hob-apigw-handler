package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

type HandlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestProcessRequests() {

	handler := handlerForTest()
	record := TimeTrackingRecord{DeviceId: "Device01", ClickType: SINGLE_CLICK}

	suite.Nil(handler.Process(record))
	now := time.Now()
	record.Timestamp = &now
	suite.Nil(handler.Process(record))
}
func (suite *HandlerTestSuite) TestConvertClickType() {

	suite.Equal(timetracker.WORKDAY, toTimeTrackingRecordType(SINGLE_CLICK))
	suite.Equal(timetracker.ILLNESS, toTimeTrackingRecordType(DOUBLE_CLICK))
	suite.Equal(timetracker.VACATION, toTimeTrackingRecordType(LONG_PRESS))
}

func handlerForTest() *APIGatewayRequestHandler {
	return newRequestHandler(timetracker.NewLocaLRepository(), loggerForTest())
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}
