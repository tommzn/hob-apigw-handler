package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

func (suite *RouterTestSuite) TestProcessRequests() {

	router := routerForTest()

	res1, err1 := router.Process(emptyRequestForResource("/success"))
	suite.Nil(err1)
	suite.Equal(200, res1.StatusCode)

	res2, err2 := router.Process(emptyRequestForResource("/error"))
	suite.NotNil(err2)
	suite.Equal(500, res2.StatusCode)

	res3, err3 := router.Process(emptyRequestForResource("/xxx"))
	suite.NotNil(err3)
	suite.Equal(500, res3.StatusCode)
}

func routerForTest() *RequestRouter {
	routes := make(map[RequestedResource]Handler)
	routes["/success"] = newHandlerMockForTest(false)
	routes["/error"] = newHandlerMockForTest(true)
	return newRequestRouter(routes, loggerForTest())
}

func emptyRequestForResource(requestedResource string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{Resource: requestedResource}
}
