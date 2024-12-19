// internal/service/webhook.go

package service

import (
	"context"
	"errors"
	"fmt"
)

type WebhookService interface {
	ProcessNotification(ctx context.Context, payload map[string]interface{}) error
}

type webhookService struct{}

func NewWebhookService() WebhookService {
	return &webhookService{}
}

func (s *webhookService) ProcessNotification(ctx context.Context, payload map[string]interface{}) error {
	// Log the incoming payload (optional)
	fmt.Println("Webhook received:", payload)

	// Extract necessary fields
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return errors.New("invalid or missing order_id in payload")
	}

	status, ok := payload["transaction_status"].(string)
	if !ok {
		return errors.New("invalid or missing transaction_status in payload")
	}

	// Example: Handle the order and status
	fmt.Printf("Processing order %s with status %s\n", orderID, status)

	// Add logic to update your database or perform necessary actions based on the status
	return nil
}
