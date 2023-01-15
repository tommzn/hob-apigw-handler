package main

import (
	"strings"
	"time"
)

var dateFormatList = []string{time.RFC3339, "2006-01-02T15:04:05-0700"}

type APITime struct {
	time.Time
}

func (apiTime *APITime) UnmarshalJSON(b []byte) error {

	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	var parseError error
	for _, format := range dateFormatList {
		parsedTime, err := time.Parse(format, value)
		if err == nil {
			*apiTime = APITime{Time: parsedTime}
			return nil
		}
		parseError = err
	}
	return parseError

}

func (apiTime *APITime) AsTime() time.Time {
	return apiTime.Time
}
