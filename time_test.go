package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type APITimeTestSuite struct {
	suite.Suite
}

type testStruct struct {
	TimeStamp APITime `json:"timestamp"`
}

func TestAPITimeTestSuite(t *testing.T) {
	suite.Run(t, new(APITimeTestSuite))
}

func (suite *APITimeTestSuite) TestUnmarshalTime() {

	val1 := testStruct{}
	str1 := "{\"timestamp\":\"2023-01-02T01:02:03+0100\"}"
	suite.Nil(json.Unmarshal([]byte(str1), &val1))

	val2 := testStruct{}
	str2 := "{\"timestamp\":\"2023-01-02T01:02:03+01:00\"}"
	suite.Nil(json.Unmarshal([]byte(str2), &val2))

	suite.Equal(val1.TimeStamp.UTC().Unix(), val2.TimeStamp.UTC().Unix())

	suite.NotNil(json.Unmarshal([]byte("{\"timestamp\":\"2023-01-02T01:02:03.441Z01:00\"}"), &testStruct{}))
}
