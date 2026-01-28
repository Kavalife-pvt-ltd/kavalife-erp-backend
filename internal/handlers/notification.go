package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func InsertNotificationEvent(ctx context.Context, req models.CreateNotificationEventRequest) (*models.NotificationEvent, error) {
	var payloadJSON []byte
	var err error

	if req.Payload != nil {
		payloadJSON, err = json.Marshal(req.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal notification payload: %w", err)
		}
	}

	var (
		id int64
	)

	query := `
		INSERT INTO notification_events (
			po_id,
			event_type,
			recipient_user_id,
			status,
			payload
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = db.DB.QueryRow(
		ctx,
		query,
		req.POID,
		req.EventType,
		req.RecipientUserID,
		models.NotificationStatusPending,
		func() []byte {
			if payloadJSON == nil {
				return nil
			}
			return payloadJSON
		}(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetNotificationByID(ctx, id)
}

func GetNotificationByID(ctx context.Context, id int64) (*models.NotificationEvent, error) {
	var n models.NotificationEvent
	var (
		poIDNS        sql.NullInt64
		payloadJSONNS sql.NullString
		errorMsgNS    sql.NullString
		sentAtNT      sql.NullTime
		statusStr     string
	)

	query := `
		SELECT
			id,
			po_id,
			event_type,
			recipient_user_id,
			status,
			payload,
			error_message,
			created_at,
			sent_at
		FROM notification_events
		WHERE id = $1
	`

	err := db.DB.QueryRow(ctx, query, id).Scan(
		&n.ID,
		&poIDNS,
		&n.EventType,
		&n.RecipientUserID,
		&statusStr,
		&payloadJSONNS,
		&errorMsgNS,
		&n.CreatedAt,
		&sentAtNT,
	)
	if err != nil {
		return nil, err
	}

	n.Status = models.NotificationStatus(statusStr)

	if poIDNS.Valid {
		v := poIDNS.Int64
		n.POID = &v
	}

	if payloadJSONNS.Valid {
		var payload map[string]any
		if err := json.Unmarshal([]byte(payloadJSONNS.String), &payload); err == nil {
			n.Payload = payload
		}
	}

	if errorMsgNS.Valid {
		n.ErrorMessage = &errorMsgNS.String
	}

	if sentAtNT.Valid {
		t := sentAtNT.Time
		n.SentAt = &t
	}

	return &n, nil
}

func ListNotificationEvents(ctx context.Context, userID *int64, status *models.NotificationStatus) ([]models.NotificationEvent, error) {
	base := `
		SELECT
			id,
			po_id,
			event_type,
			recipient_user_id,
			status,
			payload,
			error_message,
			created_at,
			sent_at
		FROM notification_events
	`
	var (
		where []string
		args  []any
	)

	if userID != nil {
		where = append(where, fmt.Sprintf("recipient_user_id = $%d", len(args)+1))
		args = append(args, *userID)
	}
	if status != nil {
		where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, string(*status))
	}

	if len(where) > 0 {
		base += " WHERE " + strings.Join(where, " AND ")
	}
	base += " ORDER BY created_at DESC"

	rows, err := db.DB.Query(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.NotificationEvent

	for rows.Next() {
		var n models.NotificationEvent
		var (
			poIDNS        sql.NullInt64
			payloadJSONNS sql.NullString
			errorMsgNS    sql.NullString
			sentAtNT      sql.NullTime
			statusStr     string
		)

		if err := rows.Scan(
			&n.ID,
			&poIDNS,
			&n.EventType,
			&n.RecipientUserID,
			&statusStr,
			&payloadJSONNS,
			&errorMsgNS,
			&n.CreatedAt,
			&sentAtNT,
		); err != nil {
			return nil, err
		}

		n.Status = models.NotificationStatus(statusStr)

		if poIDNS.Valid {
			v := poIDNS.Int64
			n.POID = &v
		}

		if payloadJSONNS.Valid {
			var payload map[string]any
			if err := json.Unmarshal([]byte(payloadJSONNS.String), &payload); err == nil {
				n.Payload = payload
			}
		}

		if errorMsgNS.Valid {
			n.ErrorMessage = &errorMsgNS.String
		}

		if sentAtNT.Valid {
			t := sentAtNT.Time
			n.SentAt = &t
		}

		result = append(result, n)
	}

	return result, nil
}
