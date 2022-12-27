package main

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hob-core"
)

// NewReportGenerateRequestHandler returna handler to generate time tracking reports.
func newReportGenerateRequestHandler(logger log.Logger, publisher Publisher) *ReportGenerateRequestHandler {
	return &ReportGenerateRequestHandler{
		logger:    logger,
		publisher: publisher,
	}
}

// Process will generate and publish time tracking report for passed year/month.
func (handler *ReportGenerateRequestHandler) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	reportGenerateRequest, err := toReportGenerateRequest(request.Body)
	if err != nil {
		handler.logger.Error(err)
		return errorResponse(err), err
	}
	handler.logger.Statusf("Report requested. type: %s, year: %d, month: %d", reportGenerateRequest.Type, reportGenerateRequest.Year, reportGenerateRequest.Month)

	event := &core.GenerateReportRequest{
		Format:      core.ReportFormat_EXCEL,
		Type:        toReportType(reportGenerateRequest.Type),
		Year:        int64(reportGenerateRequest.Year),
		Month:       int64(reportGenerateRequest.Month),
		NamePattern: "TimeTrackingReport_200601",
		Delivery: &core.ReportDelivery{
			S3: &core.S3Target{},
		},
	}

	publishErr := handler.publisher.Send(event)
	if publishErr != nil {
		handler.logger.Error(publishErr)
		return errorResponse(publishErr), publishErr
	}

	return successfulResponse(), nil
}

// ToReportGenerateRequest try to convert passed request body to a report generate request.
func toReportGenerateRequest(requestBody string) (ReportGenerateRequest, error) {
	var reportGenerateRequest ReportGenerateRequest
	err := json.Unmarshal([]byte(requestBody), &reportGenerateRequest)
	return reportGenerateRequest, err
}

func toReportType(reportType string) core.ReportType {

	switch strings.ToLower(reportType) {
	case "monthly":
		return core.ReportType_MONTHLY_REPORT
	default:
		return core.ReportType_NO_TYPE
	}
}

func cloudWatchTrigger() (string, error) {
	event := &core.GenerateReportRequest{
		Format:      core.ReportFormat_EXCEL,
		Type:        core.ReportType_MONTHLY_REPORT,
		NamePattern: "TimeTrackingReport_200601",
		Delivery: &core.ReportDelivery{
			S3: &core.S3Target{},
		},
	}
	return core.SerializeEvent(event)
}
