package main

import (
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	core "github.com/tommzn/hob-core"
	"testing"
)

type PublisherTestSuite struct {
	suite.Suite
	conf config.Config
}

func TestPublisherTestSuite(t *testing.T) {
	suite.Run(t, new(PublisherTestSuite))
}

func (suite *PublisherTestSuite) SetupTest() {
	suite.conf = configForTest()
}

func (suite *PublisherTestSuite) TestPublishMessage() {

	skipCI(suite.T())

	publisher := newSqsPublisher(suite.conf, loggerForTest())
	suite.NotNil(publisher)

	event := eventForTest()
	suite.Nil(publisher.Send(event))

	publisher1 := newSqsPublisher(emptyConfigForTest(), loggerForTest())
	suite.NotNil(publisher1)
	suite.NotNil(publisher1.Send(event))
}

func eventForTest() *core.GenerateReportRequest {
	return &core.GenerateReportRequest{
		Format:      core.ReportFormat_EXCEL,
		Year:        2022,
		Month:       1,
		NamePattern: "TimeTrackingReport_200601",
	}
}
