package fcm

import (
	"context"
	"encoding/json"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	enabled bool
	client  *messaging.Client
}

type PushMessage struct {
	Token   string
	Title   string
	Body    string
	OrderID string
	Event   string
}

func New(ctx context.Context, enabled bool, projectID string, credentialsJSON string) (*Client, error) {
	if !enabled {
		return &Client{enabled: false}, nil
	}

	if projectID == "" {
		return nil, errors.New("FCM_PROJECT_ID is required when FCM_ENABLED=true")
	}

	if credentialsJSON == "" {
		return nil, errors.New("FCM_CREDENTIALS_JSON is required when FCM_ENABLED=true")
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(credentialsJSON), &raw); err != nil {
		return nil, errors.New("invalid FCM_CREDENTIALS_JSON")
	}

	app, err := firebase.NewApp(
		ctx,
		&firebase.Config{ProjectID: projectID},
		option.WithCredentialsJSON([]byte(credentialsJSON)),
	)
	if err != nil {
		return nil, err
	}

	msgClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		enabled: true,
		client:  msgClient,
	}, nil
}

func (c *Client) SendToToken(ctx context.Context, push PushMessage) (string, error) {
	if !c.enabled {
		return "fcm-disabled", nil
	}

	if push.Token == "" {
		return "", errors.New("push token is required")
	}

	message := &messaging.Message{
		Token: push.Token,
		Notification: &messaging.Notification{
			Title: push.Title,
			Body:  push.Body,
		},
		Data: map[string]string{
			"order_id": push.OrderID,
			"event":    push.Event,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	return c.client.Send(ctx, message)
}
