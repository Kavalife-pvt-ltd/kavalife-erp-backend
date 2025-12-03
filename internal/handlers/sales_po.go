package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

//
// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────
//

// generate next PO number like "PO-012025-001"
func getNextPONumber(ctx context.Context, ts time.Time) (string, error) {
	monthYear := ts.Format("012006") // MMYYYY
	prefix := "PO-" + monthYear + "-"

	var lastPONumber string

	query := `
		SELECT po_number FROM sales_po
		WHERE po_number LIKE $1
		ORDER BY id DESC
		LIMIT 1
	`
	err := db.DB.QueryRow(ctx, query, prefix+"%").Scan(&lastPONumber)
	if err != nil && err.Error() != "no rows in result set" {
		return "", err
	}

	var next int
	if lastPONumber == "" {
		next = 1
	} else {
		parts := strings.Split(lastPONumber, "-")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid PO number format: %s", lastPONumber)
		}
		n, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("invalid number in PO: %w", err)
		}
		next = n + 1
	}

	return fmt.Sprintf("%s%03d", prefix, next), nil
}

func CreateSalesPO(ctx context.Context, req models.CreateSalesPORequest) (*models.SalesPO, error) {
	// basic sanity
	if req.ProductID == 0 {
		return nil, errors.New("productId is required")
	}
	if req.CompanyName == "" || req.CompanyAddress == "" {
		return nil, errors.New("companyName and companyAddress are required")
	}
	if req.RequestType != models.RequestTypeSample && req.RequestType != models.RequestTypePurchase {
		return nil, errors.New("invalid requestType")
	}

	now := time.Now()
	requestDate := now
	if req.RequestDate != nil {
		requestDate = *req.RequestDate
	}

	poNumber, err := getNextPONumber(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PO number: %v", err)
	}

	var id int64

	query := `
	INSERT INTO sales_po (
		po_number,
		product_id,
		company_name,
		company_address,
		coa_url,
		company_contact_name,
		company_contact_number,
		company_contact_email,
		purity,
		grade,
		request_type,
		quantity,
		quantity_unit,
		asking_price,
		comments,
		expected_delivery_date,
		request_date,
		status,
		sales_rep_id,
		send_to
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8,
		$9, $10, $11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20
	)
	RETURNING id
	`

	err = db.DB.QueryRow(
		ctx,
		query,
		poNumber,
		req.ProductID,
		req.CompanyName,
		req.CompanyAddress,
		req.COAURL,
		req.CompanyContactName,
		req.CompanyContactNumber,
		req.CompanyContactEmail,
		req.Purity,
		req.Grade,
		req.RequestType,
		req.Quantity,
		req.QuantityUnit,
		req.AskingPrice,
		req.Comments,
		req.ExpectedDeliveryDate,
		requestDate,
		models.StatusQuoteRequested,
		req.SalesRepID,
		"admin", // new POs go first to admin
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Fetch the full PO struct
	po, err := GetSalesPOByID(ctx, id)
	if err != nil {
		return nil, err
	}

	adminRows, err := db.DB.Query(ctx, `SELECT id, email FROM public.users WHERE role = 'admin'`)
	if err == nil {
		defer adminRows.Close()

		// decide PO number to display
		poNum := poNumber
		if po.PONumber != nil {
			poNum = *po.PONumber
		}

		// Prepare template data
		quantityUnit := ""
		if po.QuantityUnit != nil {
			quantityUnit = *po.QuantityUnit
		}
		purity := ""
		if po.Purity != nil {
			purity = *po.Purity
		}
		grade := ""
		if po.Grade != nil {
			grade = *po.Grade
		}
		expectedDelivery := ""
		if po.ExpectedDeliveryDate != nil {
			expectedDelivery = po.ExpectedDeliveryDate.Format("02 Jan 2006")
		}
		createdAtStr := po.CreatedAt.Format("02 Jan 2006 15:04")

		additionalComment := ""
		if po.Comments != nil {
			additionalComment = *po.Comments
		}

		templateData := utils.NewPONotificationData{
			PONumber:          poNum,
			CompanyName:       po.CompanyName,
			CompanyAddress:    po.CompanyAddress,
			RequestType:       string(po.RequestType),
			Quantity:          po.Quantity,
			QuantityUnit:      quantityUnit,
			Purity:            purity,
			Grade:             grade,
			ExpectedDelivery:  expectedDelivery,
			CreatedAt:         createdAtStr,
			AdditionalComment: additionalComment,
		}

		for adminRows.Next() {
			var adminID int64
			var adminEmail string

			if scanErr := adminRows.Scan(&adminID, &adminEmail); scanErr != nil {
				continue
			}

			// Build template — subject + HTML body
			subject, body, tmplErr := utils.BuildNewPONotificationEmail(templateData)
			if tmplErr != nil {
				// don't block PO creation
				continue
			}

			// Log notification event in DB
			_, _ = InsertNotificationEvent(ctx, models.CreateNotificationEventRequest{
				POID:            &po.ID,
				EventType:       "new_po_created",
				RecipientUserID: adminID,
				Payload: map[string]any{
					"poNumber": poNum,
					"company":  po.CompanyName,
					"status":   po.Status,
				},
			})

			// Send email (best-effort)
			if adminEmail != "" {
				_ = utils.SendEmail(adminEmail, subject, body)
			}
		}
	}

	return po, nil
}

// GetSalesPOByID fetches a single PO row by id.
func GetSalesPOByID(ctx context.Context, poID int64) (*models.SalesPO, error) {
	var po models.SalesPO

	var (
		poNumberNS, coaURLNS,
		contactNameNS, contactNumberNS, contactEmailNS,
		purityNS, gradeNS,
		quantityUnitNS, commentsNS,
		rejectionReasonNS,
		fulfillmentTypeNS,
		deliveryCodeNS sql.NullString
	)

	var (
		salesRepIDNS, approvedByIDNS, rejectedByIDNS,
		purchaseOrderIDNS, productionBatchIDNS,
		packedByIDNS sql.NullInt64
	)

	var (
		expectedDateNT, approvedAtNT, packedAtNT sql.NullTime
	)

	var askingPriceNF sql.NullFloat64

	query := `
	SELECT
		id,
		po_number,
		product_id,
		company_name,
		company_address,
		coa_url,
		company_contact_name,
		company_contact_number,
		company_contact_email,
		purity,
		grade,
		request_type,
		quantity,
		quantity_unit,
		asking_price,
		comments,
		expected_delivery_date,
		request_date,
		status,
		sales_rep_id,
		approved_by_id,
		approved_at,
		rejected_by_id,
		rejection_reason,
		fulfillment_type,
		purchase_order_id,
		production_batch_id,
		packed_by_id,
		packed_at,
		delivery_code,
		send_to,
		created_at,
		updated_at
	FROM sales_po
	WHERE id = $1
	`
	var sendTo string
	err := db.DB.QueryRow(ctx, query, poID).Scan(
		&po.ID,
		&poNumberNS,
		&po.ProductID,
		&po.CompanyName,
		&po.CompanyAddress,
		&coaURLNS,
		&contactNameNS,
		&contactNumberNS,
		&contactEmailNS,
		&purityNS,
		&gradeNS,
		&po.RequestType,
		&po.Quantity,
		&quantityUnitNS,
		&askingPriceNF,
		&commentsNS,
		&expectedDateNT,
		&po.RequestDate,
		&po.Status,
		&salesRepIDNS,
		&approvedByIDNS,
		&approvedAtNT,
		&rejectedByIDNS,
		&rejectionReasonNS,
		&fulfillmentTypeNS,
		&purchaseOrderIDNS,
		&productionBatchIDNS,
		&packedByIDNS,
		&packedAtNT,
		&deliveryCodeNS,
		&sendTo,
		&po.CreatedAt,
		&po.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if poNumberNS.Valid {
		po.PONumber = &poNumberNS.String
	}
	if coaURLNS.Valid {
		po.COAURL = &coaURLNS.String
	}
	if contactNameNS.Valid {
		po.CompanyContactName = &contactNameNS.String
	}
	if contactNumberNS.Valid {
		po.CompanyContactNumber = &contactNumberNS.String
	}
	if contactEmailNS.Valid {
		po.CompanyContactEmail = &contactEmailNS.String
	}
	if purityNS.Valid {
		po.Purity = &purityNS.String
	}
	if gradeNS.Valid {
		po.Grade = &gradeNS.String
	}
	if quantityUnitNS.Valid {
		po.QuantityUnit = &quantityUnitNS.String
	}
	if askingPriceNF.Valid {
		v := askingPriceNF.Float64
		po.AskingPrice = &v
	}
	if commentsNS.Valid {
		po.Comments = &commentsNS.String
	}
	if expectedDateNT.Valid {
		t := expectedDateNT.Time
		po.ExpectedDeliveryDate = &t
	}
	if salesRepIDNS.Valid {
		v := salesRepIDNS.Int64
		po.SalesRepID = &v
	}
	if approvedByIDNS.Valid {
		v := approvedByIDNS.Int64
		po.ApprovedByID = &v
	}
	if approvedAtNT.Valid {
		t := approvedAtNT.Time
		po.ApprovedAt = &t
	}
	if rejectedByIDNS.Valid {
		v := rejectedByIDNS.Int64
		po.RejectedByID = &v
	}
	if rejectionReasonNS.Valid {
		po.RejectionReason = &rejectionReasonNS.String
	}
	if fulfillmentTypeNS.Valid {
		ft := models.SalesPOFulfillmentType(fulfillmentTypeNS.String)
		po.FulfillmentType = &ft
	}
	if purchaseOrderIDNS.Valid {
		v := purchaseOrderIDNS.Int64
		po.PurchaseOrderID = &v
	}
	if productionBatchIDNS.Valid {
		v := productionBatchIDNS.Int64
		po.ProductionBatchID = &v
	}
	if packedByIDNS.Valid {
		v := packedByIDNS.Int64
		po.PackedByID = &v
	}
	if packedAtNT.Valid {
		t := packedAtNT.Time
		po.PackedAt = &t
	}
	if deliveryCodeNS.Valid {
		po.DeliveryCode = &deliveryCodeNS.String
	}

	po.SendTo = sendTo

	return &po, nil
}

// SalesPOFilter is used for ListSalesPO dynamic querying.
type SalesPOFilter struct {
	Status     *string
	SalesRepID *int64
	ProductID  *int64
	SendTo     *string
	// later you can add date ranges etc.
}

// handlers/sales_po.go
func ListSalesPO(ctx context.Context, filter SalesPOFilter) ([]models.SalesPO, error) {
	base := `
	SELECT
		id,
		po_number,
		product_id,
		company_name,
		company_address,
		coa_url,
		company_contact_name,
		company_contact_number,
		company_contact_email,
		purity,
		grade,
		request_type,
		quantity,
		quantity_unit,
		asking_price,
		comments,
		expected_delivery_date,
		request_date,
		status,
		sales_rep_id,
		approved_by_id,
		approved_at,
		rejected_by_id,
		rejection_reason,
		fulfillment_type,
		purchase_order_id,
		production_batch_id,
		packed_by_id,
		packed_at,
		delivery_code,
		send_to,
		created_at,
		updated_at
	FROM sales_po
	`

	var (
		where []string
		args  []interface{}
	)

	if filter.Status != nil && *filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *filter.Status)
	}
	if filter.SalesRepID != nil && *filter.SalesRepID > 0 {
		where = append(where, fmt.Sprintf("sales_rep_id = $%d", len(args)+1))
		args = append(args, *filter.SalesRepID)
	}
	if filter.ProductID != nil && *filter.ProductID > 0 {
		where = append(where, fmt.Sprintf("product_id = $%d", len(args)+1))
		args = append(args, *filter.ProductID)
	}
	if filter.SendTo != nil && *filter.SendTo != "" {
		where = append(where, fmt.Sprintf("send_to = $%d", len(args)+1))
		args = append(args, *filter.SendTo)
	}

	query := base
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.SalesPO

	for rows.Next() {
		var po models.SalesPO

		var (
			poNumberNS, coaURLNS,
			contactNameNS, contactNumberNS, contactEmailNS,
			purityNS, gradeNS,
			quantityUnitNS, commentsNS,
			rejectionReasonNS,
			fulfillmentTypeNS,
			deliveryCodeNS sql.NullString
		)
		var (
			salesRepIDNS, approvedByIDNS, rejectedByIDNS,
			purchaseOrderIDNS, productionBatchIDNS,
			packedByIDNS sql.NullInt64
		)
		var (
			expectedDateNT, approvedAtNT, packedAtNT sql.NullTime
		)
		var askingPriceNF sql.NullFloat64
		var sendTo string

		if err := rows.Scan(
			&po.ID,
			&poNumberNS,
			&po.ProductID,
			&po.CompanyName,
			&po.CompanyAddress,
			&coaURLNS,
			&contactNameNS,
			&contactNumberNS,
			&contactEmailNS,
			&purityNS,
			&gradeNS,
			&po.RequestType,
			&po.Quantity,
			&quantityUnitNS,
			&askingPriceNF,
			&commentsNS,
			&expectedDateNT,
			&po.RequestDate,
			&po.Status,
			&salesRepIDNS,
			&approvedByIDNS,
			&approvedAtNT,
			&rejectedByIDNS,
			&rejectionReasonNS,
			&fulfillmentTypeNS,
			&purchaseOrderIDNS,
			&productionBatchIDNS,
			&packedByIDNS,
			&packedAtNT,
			&deliveryCodeNS,
			&sendTo,
			&po.CreatedAt,
			&po.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if poNumberNS.Valid {
			po.PONumber = &poNumberNS.String
		}
		if coaURLNS.Valid {
			po.COAURL = &coaURLNS.String
		}
		if contactNameNS.Valid {
			po.CompanyContactName = &contactNameNS.String
		}
		if contactNumberNS.Valid {
			po.CompanyContactNumber = &contactNumberNS.String
		}
		if contactEmailNS.Valid {
			po.CompanyContactEmail = &contactEmailNS.String
		}
		if purityNS.Valid {
			po.Purity = &purityNS.String
		}
		if gradeNS.Valid {
			po.Grade = &gradeNS.String
		}
		if quantityUnitNS.Valid {
			po.QuantityUnit = &quantityUnitNS.String
		}
		if askingPriceNF.Valid {
			v := askingPriceNF.Float64
			po.AskingPrice = &v
		}
		if commentsNS.Valid {
			po.Comments = &commentsNS.String
		}
		if expectedDateNT.Valid {
			t := expectedDateNT.Time
			po.ExpectedDeliveryDate = &t
		}
		if salesRepIDNS.Valid {
			v := salesRepIDNS.Int64
			po.SalesRepID = &v
		}
		if approvedByIDNS.Valid {
			v := approvedByIDNS.Int64
			po.ApprovedByID = &v
		}
		if approvedAtNT.Valid {
			t := approvedAtNT.Time
			po.ApprovedAt = &t
		}
		if rejectedByIDNS.Valid {
			v := rejectedByIDNS.Int64
			po.RejectedByID = &v
		}
		if rejectionReasonNS.Valid {
			po.RejectionReason = &rejectionReasonNS.String
		}
		if fulfillmentTypeNS.Valid {
			ft := models.SalesPOFulfillmentType(fulfillmentTypeNS.String)
			po.FulfillmentType = &ft
		}
		if purchaseOrderIDNS.Valid {
			v := purchaseOrderIDNS.Int64
			po.PurchaseOrderID = &v
		}
		if productionBatchIDNS.Valid {
			v := productionBatchIDNS.Int64
			po.ProductionBatchID = &v
		}
		if packedByIDNS.Valid {
			v := packedByIDNS.Int64
			po.PackedByID = &v
		}
		if packedAtNT.Valid {
			t := packedAtNT.Time
			po.PackedAt = &t
		}
		if deliveryCodeNS.Valid {
			po.DeliveryCode = &deliveryCodeNS.String
		}

		po.SendTo = sendTo

		result = append(result, po)
	}

	return result, nil
}

func isValidStatusTransition(from, to models.SalesPOStatus) bool {
	// Map of allowed "from → to" transitions
	allowed := map[models.SalesPOStatus][]models.SalesPOStatus{
		models.StatusQuoteRequested: {
			models.StatusQuoteAdminApproved,
			models.StatusAdminRejected,
			models.StatusCancelled,
		},

		models.StatusQuoteAdminApproved: {
			models.StatusQuoteSentToClient,
			models.StatusClientNegotiation,
			models.StatusCancelled,
		},
		models.StatusQuoteSentToClient: {
			models.StatusClientNegotiation,
			models.StatusClientApproved,
			models.StatusClientRejected,
			models.StatusCancelled,
		},
		models.StatusClientNegotiation: {
			models.StatusQuoteAdminApproved,
			models.StatusQuoteSentToClient,
			models.StatusClientApproved,
			models.StatusClientRejected,
			models.StatusCancelled,
		},
		models.StatusClientApproved: {
			models.StatusRoutedToPurchase,
			models.StatusRoutedToProduction,
			models.StatusAdminRejected,
			models.StatusCancelled,
		},
		models.StatusRoutedToPurchase: {
			models.StatusPurchaseCompleted,
			models.StatusCancelled,
		},

		models.StatusRoutedToProduction: {
			models.StatusProductionCompleted,
			models.StatusCancelled,
		},
		models.StatusPurchaseCompleted: {
			models.StatusFinalAdminApproved,
			models.StatusCancelled,
		},
		models.StatusProductionCompleted: {
			models.StatusFinalAdminApproved,
			models.StatusCancelled,
		},
		models.StatusFinalAdminApproved: {},
		models.StatusAdminRejected: {
			models.StatusQuoteRequested,
			models.StatusCancelled,
		},
		models.StatusClientRejected: {
			models.StatusQuoteRequested,
			models.StatusCancelled,
		},
		models.StatusCancelled: {},
	}

	nextList, ok := allowed[from]
	if !ok {
		return false
	}
	for _, n := range nextList {
		if n == to {
			return true
		}
	}
	return false
}

func insertStatusLog(ctx context.Context, input models.CreateStatusLogInput) error {
	query := `
	INSERT INTO sales_po_status_log (po_id, from_status, to_status, changed_by, note)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.DB.Exec(
		ctx,
		query,
		input.POID,
		input.FromStatus,
		input.ToStatus,
		input.ChangedBy,
		input.Note,
	)
	return err
}

// UpdateSalesPOStatus updates status + related fields and logs the transition.
func UpdateSalesPOStatus(ctx context.Context, req models.UpdateSalesPOStatusRequest, actorID int64) (*models.SalesPO, error) {
	// 1. Load current PO
	current, err := GetSalesPOByID(ctx, req.POID)
	if err != nil {
		return nil, err
	}

	fromStatus := current.Status
	toStatus := req.ToStatus

	if !isValidStatusTransition(fromStatus, toStatus) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", fromStatus, toStatus)
	}

	now := time.Now()
	setParts := []string{"status = $1", "updated_at = $2"}
	args := []interface{}{toStatus, now}
	argPos := 3

	// negotiation / editable fields
	if req.NewQuantity != nil {
		setParts = append(setParts, fmt.Sprintf("quantity = $%d", argPos))
		args = append(args, *req.NewQuantity)
		argPos++
	}
	if req.NewAskingPrice != nil {
		setParts = append(setParts, fmt.Sprintf("asking_price = $%d", argPos))
		args = append(args, *req.NewAskingPrice)
		argPos++
	}
	if req.NewComments != nil {
		setParts = append(setParts, fmt.Sprintf("comments = $%d", argPos))
		args = append(args, *req.NewComments)
		argPos++
	}

	var noteParts []string

	switch toStatus {
	case models.StatusQuoteAdminApproved,
		models.StatusFinalAdminApproved:
		setParts = append(setParts,
			fmt.Sprintf("approved_by_id = $%d", argPos),
			fmt.Sprintf("approved_at = $%d", argPos+1),
		)
		args = append(args, actorID, now)
		argPos += 2

	case models.StatusAdminRejected:
		setParts = append(setParts,
			fmt.Sprintf("rejected_by_id = $%d", argPos),
			fmt.Sprintf("rejection_reason = $%d", argPos+1),
		)
		args = append(args, actorID, req.RejectionReason)
		argPos += 2

		if req.RejectionReason != nil {
			noteParts = append(noteParts, "Rejection: "+*req.RejectionReason)
		}
	}

	if toStatus == models.StatusRoutedToPurchase || toStatus == models.StatusRoutedToProduction {
		var ft models.SalesPOFulfillmentType
		if req.FulfillmentType != nil {
			ft = *req.FulfillmentType
		} else {
			if toStatus == models.StatusRoutedToPurchase {
				ft = models.FulfillmentPurchase
			} else {
				ft = models.FulfillmentProduction
			}
		}
		setParts = append(setParts, fmt.Sprintf("fulfillment_type = $%d", argPos))
		args = append(args, ft)
		argPos++
	}

	// completion by purchase / production
	if toStatus == models.StatusPurchaseCompleted {
		noteParts = append(noteParts, "Marked completed by purchase team")
	}
	if toStatus == models.StatusProductionCompleted {
		noteParts = append(noteParts, "Marked completed by production team")
	}
	if toStatus == models.StatusClosed {
		noteParts = append(noteParts, "PO closed by admin")
	}

	// delivery code (optional)
	if req.DeliveryCode != nil {
		setParts = append(setParts, fmt.Sprintf("delivery_code = $%d", argPos))
		args = append(args, *req.DeliveryCode)
		argPos++
		noteParts = append(noteParts, "DeliveryCode set/updated")
	}

	// send_to logic
	// 1) explicit sendTo from request (frontend-driven)
	if req.SendTo != nil && *req.SendTo != "" {
		setParts = append(setParts, fmt.Sprintf("send_to = $%d", argPos))
		args = append(args, *req.SendTo)
		argPos++
	} else {
		// 2) automatic defaults based on toStatus if not provided
		var autoSendTo string

		switch toStatus {
		case models.StatusQuoteRequested:
			// Sales is (re)requesting quote → goes to admin
			autoSendTo = "admin"
		case models.StatusAdminRejected:
			// Goes back to sales
			autoSendTo = "sales"
		case models.StatusRoutedToPurchase:
			autoSendTo = "purchase"
		case models.StatusRoutedToProduction:
			autoSendTo = "production"
		case models.StatusPurchaseCompleted, models.StatusProductionCompleted, models.StatusFinalAdminApproved:
			autoSendTo = "admin"
		}

		if autoSendTo != "" {
			setParts = append(setParts, fmt.Sprintf("send_to = $%d", argPos))
			args = append(args, autoSendTo)
			argPos++
		}
	}

	// Final UPDATE
	query := fmt.Sprintf(`
        UPDATE sales_po
        SET %s
        WHERE id = $%d
    `, strings.Join(setParts, ", "), argPos)

	args = append(args, req.POID)

	tag, err := db.DB.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("no sales_po found with id %d", req.POID)
	}

	var note *string
	if len(noteParts) > 0 {
		joined := strings.Join(noteParts, " | ")
		note = &joined
	}
	logInput := models.CreateStatusLogInput{
		POID:       req.POID,
		FromStatus: &fromStatus,
		ToStatus:   toStatus,
		ChangedBy:  &actorID,
		Note:       note,
	}
	_ = insertStatusLog(ctx, logInput)

	// 4. Return updated PO
	return GetSalesPOByID(ctx, req.POID)
}

// GetSalesPOStatusLog fetches all status transitions for a PO.
func GetSalesPOStatusLog(ctx context.Context, poID int64) ([]models.SalesPOStatusLog, error) {
	rows, err := db.DB.Query(ctx, `
		SELECT
			id,
			po_id,
			from_status,
			to_status,
			changed_by,
			note,
			changed_at
		FROM sales_po_status_log
		WHERE po_id = $1
		ORDER BY changed_at ASC
	`, poID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.SalesPOStatusLog

	for rows.Next() {
		var log models.SalesPOStatusLog
		var fromStatusNS sql.NullString
		var changedByNS sql.NullInt64
		var noteNS sql.NullString

		if err := rows.Scan(
			&log.ID,
			&log.POID,
			&fromStatusNS,
			&log.ToStatus,
			&changedByNS,
			&noteNS,
			&log.ChangedAt,
		); err != nil {
			return nil, err
		}

		if fromStatusNS.Valid {
			fs := models.SalesPOStatus(fromStatusNS.String)
			log.FromStatus = &fs
		}
		if changedByNS.Valid {
			v := changedByNS.Int64
			log.ChangedBy = &v
		}
		if noteNS.Valid {
			log.Note = &noteNS.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}
