package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
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
		deviceIds := deviceIdsFromRequest(request)
		if len(deviceIds) == 0 {
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
		handler.logger.Debugf("Receive GET for DeviceId: %s, Date: %s", strings.Join(deviceIds, ","), dateStr)

		timeRangeStart, timeRangeEnd := handler.timeRangeForDate("2006-01-02", dateStr)
		if timeRangeStart == nil || timeRangeEnd == nil {
			err := errors.New("Unable to determin time rage for date: " + dateStr)
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		handler.logger.Debugf("Looking for records in range %s/%s", timeRangeStart.Format(time.RFC3339), timeRangeEnd.Format(time.RFC3339))

		records := []timetracker.TimeTrackingRecord{}
		for _, deviceId := range deviceIds {
			recordsForDevice, err := handler.timeTracker.ListRecords(deviceId, *timeRangeStart, *timeRangeEnd)
			if err != nil {
				handler.logger.Error(err)
				return errorResponseWithStatus(err, http.StatusInternalServerError), err
			}
			records = append(records, recordsForDevice...)
		}

		if len(records) == 0 {
			handler.logger.Errorf("No time tracking records found. (%s&%s)", strings.Join(deviceIds, ","), dateStr)
			return responseWithStatus(http.StatusNotFound), nil
		}

		for idx, _ := range records {
			records[idx].Key = queryExcapeKey(records[idx].Key)
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

		id, ok := request.QueryStringParameters["id"]
		if !ok {
			err := errors.New("Missing time tracking record id.")
			handler.logger.Error(err)
			return errorResponseWithStatus(err, http.StatusBadRequest), err
		}
		decodedId := queryUnexcapeKey(id)
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

func queryExcapeKey(key string) string {
	return url.QueryEscape(key)
}

func queryUnexcapeKey(key string) string {
	unEscapedKey, _ := url.QueryUnescape(key)
	return unEscapedKey
}

func deviceIdsFromRequest(request events.APIGatewayProxyRequest) []string {

	listOfDeviceIds := []string{}
	if deviceId, ok := request.QueryStringParameters["deviceid"]; ok {
		listOfDeviceIds = append(listOfDeviceIds, deviceId)
	}
	if deviceIdStr, ok := request.QueryStringParameters["deviceids"]; ok {
		if excapedDeviceIds, err := url.QueryUnescape(deviceIdStr); err == nil {
			deviceIds := strings.Split(excapedDeviceIds, ",")
			listOfDeviceIds = append(listOfDeviceIds, deviceIds...)
		}
	}
	return listOfDeviceIds
}
