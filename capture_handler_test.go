package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/suite"
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

	res1, err1 := handler.Process(suite.requestForTest(record))
	suite.Nil(err1)
	suite.Equal(200, res1.StatusCode)

	now := time.Now()
	record.Timestamp = &now
	res2, err2 := handler.Process(suite.requestForTest(record))
	suite.Nil(err2)
	suite.Equal(200, res2.StatusCode)
}

func (suite *HandlerTestSuite) TestConvertClickType() {

	suite.Equal(timetracker.WORKDAY, toTimeTrackingRecordType(SINGLE_CLICK))
	suite.Equal(timetracker.ILLNESS, toTimeTrackingRecordType(DOUBLE_CLICK))
	suite.Equal(timetracker.VACATION, toTimeTrackingRecordType(LONG_PRESS))
}

func (suite *HandlerTestSuite) requestForTest(record TimeTrackingRecord) events.APIGatewayProxyRequest {
	content, err := json.Marshal(record)
	suite.Nil(err)
	return events.APIGatewayProxyRequest{Body: string(content)}
}

func handlerForTest() *CaptureRequestHandler {
	return newCaptureRequestHandler(timetracker.NewLocaLRepository(), loggerForTest())
}
