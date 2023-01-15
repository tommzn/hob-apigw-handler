package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

// NewReportGenerateRequestHandler returna handler to maintina, add and delete, time tracking records.
func newTimeTrackingRecordHandler(manager timetracker.TimeTrackingRecordManager, timeTracker timetracker.TimeTracker, logger log.Logger) *TimeTrackingRecordHandler {
	return &TimeTrackingRecordHandler{
		logger:              logger,
		timeTrackingManager: manager,
		timeTracker:         timeTracker,
	}
}

// Process will generate and publish time tracking report for passed year/month.
func (handler *TimeTrackingRecordHandler) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	handler.logger.Debugf("Request received: %s %s, Query: %+v, PathParams: %+v", request.HTTPMethod, request.Path, request.QueryStringParameters, request.PathParameters)

	switch request.HTTPMethod {

	case http.MethodGet:
		deviceId, ok1 := request.QueryStringParameters["deviceid"]
		if !ok1 {
			err := errors.New("Missing device id.")
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		dateStr, ok2 := request.QueryStringParameters["date"]
		if !ok2 {
			err := errors.New("Missing date.")
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		handler.logger.Debugf("Receive GET for DeviceId: %s, Date: %s", deviceId, dateStr)

		timeRangeStart, timeRangeEnd := handler.timeRangeForDate("2006-01-02", dateStr)
		if timeRangeStart == nil || timeRangeEnd == nil {
			err := errors.New("Unable to determin time rage for date: " + dateStr)
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		handler.logger.Debugf("Looking for records in range %s/%s", timeRangeStart.Format(time.RFC3339), timeRangeEnd.Format(time.RFC3339))

		records, err := handler.timeTracker.ListRecords(deviceId, *timeRangeStart, *timeRangeEnd)
		handler.logger.Debug(err)
		if err != nil {
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusInternalServerError), err
		}
		if len(records) == 0 {
			handler.logger.Errorf("No time tracking records found. (%s&%s)", deviceId, dateStr)
			return responseWithStatus(http.StatusNotFound), nil
		}

		for idx, _ := range records {
			records[idx].Key = encodeKey(records[idx].Key)
		}
		responseContent, err := json.Marshal(records)
		if err != nil {
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusInternalServerError), err
		}
		return responseWithContent(string(responseContent), http.StatusOK), nil

	case http.MethodPost:

		var record timetracker.TimeTrackingRecord
		if err := json.Unmarshal([]byte(request.Body), &record); err != nil {
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}

		if record.DeviceId == "" ||
			record.Type == "" ||
			record.Timestamp.Before(time.Now().Add(-2*365*24*time.Hour)) ||
			record.Timestamp.After(time.Now().Add(1*365*24*time.Hour)) {
			err := errors.New("Invalid time tracking record.")
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		handler.logger.Debugf("Receive new time tracking record: %+v", record)

		newRecord, err := handler.timeTrackingManager.Add(record)
		if err != nil {
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusInternalServerError), err
		}

		responseContent, err := json.Marshal(newRecord)
		if err != nil {
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusInternalServerError), err
		}
		return responseWithContent(string(responseContent), http.StatusCreated), nil

	case http.MethodDelete:

		id, ok := request.PathParameters["id"]
		if !ok {
			err := errors.New("Missing time tracking record id.")
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		decodedId := decodeKey(id)
		handler.logger.Debug("Receive time tracking record delete for id: ", decodedId)

		err := handler.timeTrackingManager.Delete(decodedId)
		if err != nil {
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusInternalServerError), err
		}
		return responseWithContent("", http.StatusNoContent), nil

	default:
		err := errors.New("Unsupported HTTP method.")
		handler.logger.Error(err)
		return errorResponseWithStatus(err, http.StatusMethodNotAllowed), err
	}

	err := errors.New("Invalid request.")
	handler.logger.Error(err)
	return errorResponseWithStatus(err, http.StatusBadRequest), err
}

func (handler *TimeTrackingRecordHandler) timeRangeForDate(layout, dateValue string) (*time.Time, *time.Time) {

	date, err := time.Parse(layout, dateValue)
	if err != nil {
		handler.logger.Error(err)
		return nil, nil
	}

	location := time.Now().Location()
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	return timePtr(start), timePtr(start.AddDate(0, 0, 1).Add(-1 * time.Second))
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func encodeKey(key string) string {
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func decodeKey(key string) string {
	decoded, _ := base64.StdEncoding.DecodeString(key)
	return string(decoded)
}
