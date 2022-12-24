package main

import (
	"github.com/golang/protobuf/proto"
)

// sqsMock mocks access to AWS SQS for testing.
type sqsMock struct {
	callCount int
}

// newSqsMock creates a new mock for AWS SQS.
func newSqsMock() *sqsMock {
	return &sqsMock{callCount: 0}
}

func (mock *sqsMock) Send(message proto.Message) error {
	mock.callCount++
	return nil
}
