package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
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

	timeRangeStart, timeRangeEnd := reportTimeRange(reportGenerateRequest.Year, reportGenerateRequest.Month)
	var timeTrackingRecords []timetracker.TimeTrackingRecord
	for _, deviceId := range handler.deviceIds {
		deviceRecords, err := handler.timeTracker.ListRecords(deviceId, timeRangeStart, timeRangeEnd)
		if err != nil {
			handler.logger.Error(err)
			return errorResponse(err), err
		}
		timeTrackingRecords = append(timeTrackingRecords, deviceRecords...)
	}
	handler.calculator.WithTimeTrackingRecords(timeTrackingRecords)

	if handler.calendar != nil {
		if holidays, err := handler.calendar.GetHolidays(reportGenerateRequest.Year, reportGenerateRequest.Month); err == nil {
			handler.formatter.WithHolidays(holidays)
		}
	}

	monthlyReport, err := handler.calculator.MonthlyReport(reportGenerateRequest.Year, reportGenerateRequest.Month, timetracker.WORKDAY)
	if err != nil {
		handler.logger.Error(err)
		return errorResponse(err), err
	}

	reportBuffer, err := handler.formatter.WriteMonthlyReportToBuffer(monthlyReport)
	if err != nil {
		handler.logger.Error(err)
		return errorResponse(err), err
	}

	reportFileName := "TimeTrackingReport_" + timeRangeStart.Format("200601") + handler.formatter.FileExtension()
	err = handler.publisher.Send(reportBuffer.Bytes(), reportFileName)
	if err != nil {
		handler.logger.Error(err)
		return errorResponse(err), err
	}

	return successfulResponse(), nil
}

// ToReportGenerateRequest try to convert passed request body to a report generate request.
func toReportGenerateRequest(requestBody string) (ReportGenerateRequest, error) {
	var reportGenerateRequest ReportGenerateRequest
	err := json.Unmarshal([]byte(requestBody), &reportGenerateRequest)
	return reportGenerateRequest, err
}

// ReportTimeRange generates first amd last day for report time range.
func reportTimeRange(year, month int) (time.Time, time.Time) {
	firstOfThisMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return firstOfThisMonth, firstOfThisMonth.AddDate(0, 1, 0).Add(-1 * time.Second)
}
