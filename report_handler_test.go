package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/suite"
	timetracker "github.com/tommzn/hob-timetracker"
)

type ReportHandlerTestSuite struct {
	suite.Suite
}

func TestReportHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ReportHandlerTestSuite))
}

func (suite *ReportHandlerTestSuite) TestProcessRequests() {

	handler := reportHandlerForTest()
	reportGenerateRequest := ReportGenerateRequest{Type: "monthly", Year: 2022, Month: 1}

	res1, err1 := handler.Process(suite.requestForTest(reportGenerateRequest))
	suite.Nil(err1)
	suite.Equal(200, res1.StatusCode)
}

func (suite *ReportHandlerTestSuite) requestForTest(reportGenerateRequest ReportGenerateRequest) events.APIGatewayProxyRequest {
	content, err := json.Marshal(reportGenerateRequest)
	suite.Nil(err)
	return events.APIGatewayProxyRequest{Body: string(content)}
}

func reportHandlerForTest() *ReportGenerateRequestHandler {

	conf := configForTest()
	locale := newLocale(conf)
	logger := loggerForTest()

	deviceIds := deviceIds(conf)
	formatter := newReportFormatter()
	calculator := newReportCalulator(locale)

	return &ReportGenerateRequestHandler{
		logger:      logger,
		deviceIds:   deviceIds,
		timeTracker: timeTrackeForTest(),
		calculator:  calculator,
		formatter:   formatter,
		publisher:   timetracker.NewFilePublisher(),
	}
}

func timeTrackeForTest() timetracker.TimeTracker {

	device := "Device01"
	tracker := timetracker.NewLocaLRepository()

	now := time.Now()
	firstOfThisMonth := time.Date(2022, time.Month(1), 1, 8, 0, 0, 0, now.Location())
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)
	tracker.Captured(device, timetracker.WORKDAY, firstOfLastMonth)
	tracker.Captured(device, timetracker.WORKDAY, firstOfLastMonth.Add(7*time.Hour))
	return tracker
}
