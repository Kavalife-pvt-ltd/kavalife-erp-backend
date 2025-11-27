package models

import "time"

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

type NotificationEvent struct {
	ID              int64              `json:"id"`
	POID            *int64             `json:"poId,omitempty"`
	EventType       string             `json:"eventType"`
	RecipientUserID int64              `json:"recipientUserId"`
	Status          NotificationStatus `json:"status"`
	Payload         map[string]any     `json:"payload,omitempty"`
	ErrorMessage    *string            `json:"errorMessage,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	SentAt          *time.Time         `json:"sentAt,omitempty"`
}

// For inserting a new event from handlers
type CreateNotificationEventRequest struct {
	POID            *int64         `json:"poId,omitempty"`
	EventType       string         `json:"eventType"`
	RecipientUserID int64          `json:"recipientUserId"`
	Payload         map[string]any `json:"payload,omitempty"`
}
