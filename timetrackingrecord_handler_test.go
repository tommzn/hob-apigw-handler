package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/suite"
	timetracker "github.com/tommzn/hob-timetracker"
)

type TimeTrackingRecordHandlerTestSuite struct {
	suite.Suite
}

func TestTimeTrackingRecordHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TimeTrackingRecordHandlerTestSuite))
}

func (suite *TimeTrackingRecordHandlerTestSuite) TestAddTimeTrackingRecord() {

	handler := timeTrackingRecordHandlerForTest()

	request1 := suite.requestForTest("/timetrackingrecords", http.MethodPost)
	timeTrackingRecord1 := timeTrackingRecordForTest()
	content1, err1 := json.Marshal(timeTrackingRecord1)
	suite.Nil(err1)
	request1.Body = string(content1)

	res1, err1 := handler.Process(request1)
	suite.Nil(err1)
	suite.Equal(http.StatusCreated, res1.StatusCode)
	suite.NotEqual("", res1.Body)

	var timeTrackingRecord1_1 TimeTrackingRecord
	suite.Nil(json.Unmarshal([]byte(res1.Body), &timeTrackingRecord1_1))
	suite.NotEqual("", timeTrackingRecord1_1.Key)
}

func (suite *TimeTrackingRecordHandlerTestSuite) TestListTimeTrackingRecords() {

	handler := timeTrackingRecordHandlerForTest()
	prepareForTest(handler.timeTrackingManager)

	request1 := suite.requestForTest("/timetrackingrecords", http.MethodGet)
	request1.QueryStringParameters = map[string]string{"deviceid": "Device01", "date": "2022-01-01"}
	res1, err1 := handler.Process(request1)
	suite.Nil(err1)
	suite.Equal(http.StatusOK, res1.StatusCode)
	suite.NotEqual("", res1.Body)

	var records []TimeTrackingRecord
	suite.Nil(json.Unmarshal([]byte(res1.Body), &records))
	suite.Len(records, 2)

	request2 := suite.requestForTest("/timetrackingrecords", http.MethodGet)
	request2.QueryStringParameters = map[string]string{"deviceid": "Device01", "date": "2021-01-01"}
	res2, err2 := handler.Process(request2)
	suite.Nil(err2)
	suite.Equal(http.StatusNotFound, res2.StatusCode)

	request3_1 := suite.requestForTest("/timetrackingrecords", http.MethodGet)
	request3_1.QueryStringParameters = map[string]string{"date": "2021-01-01"}
	res3_1, err3_1 := handler.Process(request3_1)
	suite.NotNil(err3_1)
	suite.Equal(http.StatusBadRequest, res3_1.StatusCode)

	request3_2 := suite.requestForTest("/timetrackingrecords", http.MethodGet)
	request3_2.QueryStringParameters = map[string]string{"deviceid": "Device01"}
	res3_2, err3_2 := handler.Process(request3_2)
	suite.NotNil(err3_2)
	suite.Equal(http.StatusBadRequest, res3_2.StatusCode)
}

func (suite *TimeTrackingRecordHandlerTestSuite) TestDeleteTimeTrackingRecords() {

	handler := timeTrackingRecordHandlerForTest()
	prepareForTest(handler.timeTrackingManager)

	request1 := suite.requestForTest("/timetrackingrecords", http.MethodGet)
	request1.QueryStringParameters = map[string]string{"deviceid": "Device01", "date": "2022-01-01"}
	res1, err1 := handler.Process(request1)
	suite.Nil(err1)
	suite.Equal(http.StatusOK, res1.StatusCode)
	suite.NotEqual("", res1.Body)

	var records []TimeTrackingRecord
	suite.Nil(json.Unmarshal([]byte(res1.Body), &records))
	suite.Len(records, 2)

	request1_1 := suite.requestForTest("/timetrackingrecords", http.MethodGet)
	request1_1.QueryStringParameters = map[string]string{"deviceids": "Device01,Device02", "date": "2022-01-01"}
	res1_1, err1_1 := handler.Process(request1_1)
	suite.Nil(err1_1)
	suite.Equal(http.StatusOK, res1_1.StatusCode)
	suite.NotEqual("", res1_1.Body)

	var records1_1 []TimeTrackingRecord
	suite.Nil(json.Unmarshal([]byte(res1_1.Body), &records1_1))
	suite.Len(records1_1, 2)

	request2 := suite.requestForTest("/timetrackingrecords", http.MethodDelete)
	request2.QueryStringParameters = map[string]string{"id": records[0].Key}
	res2, err2 := handler.Process(request2)
	suite.Nil(err2)
	suite.Equal(http.StatusNoContent, res2.StatusCode)

	res2_1, err2_1 := handler.Process(request1)
	suite.Nil(err2_1)
	suite.Equal(http.StatusOK, res2_1.StatusCode)
	suite.NotEqual("", res2_1.Body)

	var records2 []TimeTrackingRecord
	suite.Nil(json.Unmarshal([]byte(res2_1.Body), &records2))
	suite.Len(records2, 1)
}

func (suite *TimeTrackingRecordHandlerTestSuite) TestJsonMarshalTime() {

	t1 := time.Now()
	b, _ := json.Marshal(t1)
	fmt.Println(string(b))
}

func (suite *TimeTrackingRecordHandlerTestSuite) TestKeyEncoding() {

	key := "timetracker/P5SJVQ20074C6774/2022/12/24/faa5260b-01d3-41ed-bde9-8eb7bbfe9c0a"

	suite.Equal("timetracker%2FP5SJVQ20074C6774%2F2022%2F12%2F24%2Ffaa5260b-01d3-41ed-bde9-8eb7bbfe9c0a", queryExcapeKey(key))
	suite.Equal(key, queryUnexcapeKey(queryExcapeKey(key)))
}

func (suite *TimeTrackingRecordHandlerTestSuite) requestForTest(resource, httpMethod string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{Resource: resource, HTTPMethod: httpMethod}
}

func timeTrackingRecordHandlerForTest() *TimeTrackingRecordHandler {
	repo := timetracker.NewLocaLRepository()
	return newTimeTrackingRecordHandler(repo, repo, loggerForTest())
}

func prepareForTest(manager timetracker.TimeTrackingRecordManager) {

	deviceId := "Device01"
	location := time.Now().Location()
	records := []timetracker.TimeTrackingRecord{}

	records = append(records, timetracker.TimeTrackingRecord{
		DeviceId:  deviceId,
		Type:      timetracker.WORKDAY,
		Timestamp: time.Date(2022, time.Month(1), 1, 9, 0, 0, 0, location),
	})
	records = append(records, timetracker.TimeTrackingRecord{
		DeviceId:  deviceId,
		Type:      timetracker.WORKDAY,
		Timestamp: time.Date(2022, time.Month(1), 1, 17, 0, 0, 0, location),
	})
	records = append(records, timetracker.TimeTrackingRecord{
		DeviceId:  deviceId,
		Type:      timetracker.WORKDAY,
		Timestamp: time.Date(2021, time.Month(12), 31, 14, 0, 0, 0, location),
	})
	records = append(records, timetracker.TimeTrackingRecord{
		DeviceId:  deviceId,
		Type:      timetracker.WORKDAY,
		Timestamp: time.Date(2022, time.Month(1), 2, 9, 0, 0, 0, location),
	})
	for _, record := range records {
		manager.Add(record)
	}
}

func timeTrackingRecordForTest() timetracker.TimeTrackingRecord {
	return timetracker.TimeTrackingRecord{
		DeviceId:  "Device01",
		Type:      timetracker.WORKDAY,
		Timestamp: time.Now(),
	}
}
