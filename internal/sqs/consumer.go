package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Event struct {
	EventType string `json:"event_type"`
	OrderID   string `json:"order_id"`
	DriverID  string `json:"driver_id"`
	Status    string `json:"status,omitempty"`
}

type snsEnvelope struct {
	Message string `json:"Message"`
}

type Handler func(ctx context.Context, event Event) error

type Consumer struct {
	client      *sqs.Client
	queueURL    string
	waitSeconds int32
	maxMessages int32
	handler     Handler
}

func New(client *sqs.Client, queueURL string, waitSeconds int32, maxMessages int32, handler Handler) *Consumer {
	return &Consumer{
		client:      client,
		queueURL:    queueURL,
		waitSeconds: waitSeconds,
		maxMessages: maxMessages,
		handler:     handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if c.queueURL == "" {
		return errors.New("SQS_QUEUE_URL is required")
	}

	for {
		out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			WaitTimeSeconds:     c.waitSeconds,
			MaxNumberOfMessages: c.maxMessages,
		})
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		for _, msg := range out.Messages {
			event, err := parseEvent(*msg.Body)
			if err != nil {
				continue
			}

			if event.EventType == "order.created" {
				event.EventType = "NEW_ORDER"
			}

			if err := c.handler(ctx, event); err != nil {
				continue
			}

			_, _ = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
}

func parseEvent(body string) (Event, error) {
	var event Event

	if err := json.Unmarshal([]byte(body), &event); err == nil && event.EventType != "" {
		return event, nil
	}

	var envelope snsEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return event, err
	}

	if envelope.Message == "" {
		return event, errors.New("empty sns message")
	}

	if err := json.Unmarshal([]byte(envelope.Message), &event); err != nil {
		return event, err
	}

	return event, nil
}
