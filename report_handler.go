package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
)

// NewReportGenerateRequestHandler returna handler to generate time tracking reports.
func newReportGenerateRequestHandler(logger log.Logger) *ReportGenerateRequestHandler {
	return &ReportGenerateRequestHandler{logger: logger}
}

// Process will generate and publish time tracking report for passed year/month.
func (handler *ReportGenerateRequestHandler) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	reportGenerateRequest, err := toReportGenerateRequest(request.Body)
	if err != nil {
		handler.logger.Error(err)
		return errorResponse(err), err
	}
	handler.logger.Statusf("Report requested. type: %s, year: %d, month: %d", reportGenerateRequest.Type, reportGenerateRequest.Year, reportGenerateRequest.Month)

	return successfulResponse(), nil
}

// ToReportGenerateRequest try to convert passed request body to a report generate request.
func toReportGenerateRequest(requestBody string) (ReportGenerateRequest, error) {
	var reportGenerateRequest ReportGenerateRequest
	err := json.Unmarshal([]byte(requestBody), &reportGenerateRequest)
	return reportGenerateRequest, err
}
