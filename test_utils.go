package main

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

type handlerMock struct {
	shouldReturnError bool
}

func newHandlerMockForTest(shouldReturnError bool) *handlerMock {
	return &handlerMock{shouldReturnError: shouldReturnError}
}

func (mock *handlerMock) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if mock.shouldReturnError {
		err := errors.New("An error has occurred!")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func configForTest() config.Config {
	configFile := "fixtures/testconfig.yml"
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

func emptyConfigForTest() config.Config {
	configFile := "fixtures/testconfig-empty.yml"
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

// skipCI returns true if env variable CI is set
func skipCI(t *testing.T) {
	if _, isSet := os.LookupEnv("CI"); isSet {
		t.Skip("Skipping testing in CI environment")
	}
}
