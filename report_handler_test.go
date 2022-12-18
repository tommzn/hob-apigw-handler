package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/suite"
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
	return newReportGenerateRequestHandler(loggerForTest())
}
