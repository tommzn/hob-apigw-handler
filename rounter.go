package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
)

// NewRequestRouter returns a router with passed routes.
func newRequestRouter(routes map[RequestedResource]Handler, logger log.Logger) *RequestRouter {
	return &RequestRouter{logger: logger, routes: routes}
}

// Process will pick up a handler from internal routes for passed resource and forward current request to it.
// If there's no suitable route for passed request it returns with an error.
func (router *RequestRouter) Process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	defer router.logger.Flush()
	router.logger.Debugf("Requested resource: %s, path: %s", request.Resource, request.Path)

	if handler, ok := router.routes[resourceFromRequest(request)]; ok {
		return handler.Process(request)
	}
	err := fmt.Errorf("No route matching: %s", request.Resource)
	router.logger.Error(err)
	return errorResponse(err), err
}

// ResourceFromRequest extracts requested resource from passed API Gateway request.
func resourceFromRequest(request events.APIGatewayProxyRequest) RequestedResource {
	return RequestedResource(request.Resource)
}
