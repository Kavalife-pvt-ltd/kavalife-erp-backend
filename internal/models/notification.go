package models

import "time"

type NotificationType string

const (
	NotificationNewPOCreated       NotificationType = "new_po_created"
	NotificationAdminFirstApproved NotificationType = "admin_first_approved"
	NotificationAdminRejected      NotificationType = "admin_rejected"
	NotificationClientApproved     NotificationType = "client_approved"
	NotificationFinalAdminApproved NotificationType = "final_admin_approved"
	NotificationRoutedToPurchase   NotificationType = "routed_to_purchase"
	NotificationRoutedToProduction NotificationType = "routed_to_production"
)

// NotificationEvent represents a row in notification_events.
type NotificationEvent struct {
	ID int64 `db:"id" json:"id"`

	POID int64            `db:"po_id" json:"poId"`
	Type NotificationType `db:"type" json:"type"`

	RecipientEmail string `db:"recipient_email" json:"recipientEmail"`
	Subject        string `db:"subject" json:"subject"`
	Body           string `db:"body" json:"body"`

	SentAt time.Time `db:"sent_at" json:"sentAt"`
}

// LogNotificationInput can be used by handlers/services to log a sent mail.
type LogNotificationInput struct {
	POID           int64            `json:"poId"`
	Type           NotificationType `json:"type"`
	RecipientEmail string           `json:"recipientEmail"`
	Subject        string           `json:"subject"`
	Body           string           `json:"body"`
}
