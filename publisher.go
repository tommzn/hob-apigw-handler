package main

import (
	"github.com/golang/protobuf/proto"
	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hob-core"
)

// newSqsPublisher creates a new SQS message publisher with given queue and archive queue.
func newSqsPublisher(conf config.Config, logger log.Logger) *SqsPublisher {
	queue := conf.Get("hob.queue", config.AsStringPtr("de.tsl.hob.unknown"))
	return &SqsPublisher{
		logger:    logger,
		sqsClient: sqs.NewPublisher(conf),
		queue:     *queue,
	}
}

// send will publish passed message to given queues.
func (publisher *SqsPublisher) Send(message proto.Message) error {

	defer publisher.logger.Flush()

	messageString, err := core.SerializeEvent(message)
	if err != nil {
		publisher.logger.Errorf("Failed to encode event, type: %T, reason: %s", message, err)
		return err
	}

	messageId, err := publisher.sqsClient.Send(messageString, publisher.queue)
	if err != nil {
		publisher.logger.Error("Unable to semd event, reason: ", err)
		return err
	}
	publisher.logger.Infof("Event send, type: %T, queue: %s, id: %s", message, publisher.queue, *messageId)
	return nil
}
