package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func getNextGRNNumber(ctx context.Context, createdAt time.Time) (string, error) {
	monthYear := createdAt.Format("012006") // MMYYYY
	prefix := "GRN-" + monthYear + "-"

	var lastGRN string
	query := `SELECT grn_number FROM grns WHERE grn_number LIKE $1 ORDER BY id DESC LIMIT 1`
	err := db.DB.QueryRow(ctx, query, prefix+"%").Scan(&lastGRN)

	var nextNumber int
	if err != nil || lastGRN == "" {
		nextNumber = 1
	} else {
		parts := strings.Split(lastGRN, "-")
		num, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("invalid GRN number format")
		}
		nextNumber = num + 1
	}

	return fmt.Sprintf("%s%03d", prefix, nextNumber), nil
}

func InsertGRNHandler(ctx context.Context, req models.CreateGRNRequest) (*models.GRN, error) {
	createdAt := time.Now().UTC()
	grnNumber, err := getNextGRNNumber(ctx, createdAt)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO grns (
			grn_number, vir_number, container_qty,
			quantity, invoice, invoice_date, invoice_img, packaging_status,
			created_by, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var id int64
	err = db.DB.QueryRow(ctx, query,
		grnNumber, req.VIRNumber, req.ContainerQty,
		req.Quantity, req.Invoice, req.InvoiceDate, req.InvoiceImg, req.PackagingStatus,
		req.CreatedBy, createdAt,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &models.GRN{
		ID:        id,
		GRNNumber: grnNumber,
		CreatedAt: createdAt.Format(time.RFC3339),
	}, nil
}

func ViewGRNsHandler(ctx context.Context, grnNo string) ([]models.GRNResponse, error) {
	query := `
		SELECT 
			grns.id, grns.grn_number, grns.created_at, grns.container_qty, grns.quantity, grns.invoice, 
			grns.invoice_date, grns.invoice_img, grns.packaging_status, grns.created_by, 
			vir.vir_number, products.name AS product_name, vendors.name AS vendor_name
		FROM grns
		JOIN vir ON grns.vir_number = vir.vir_number
		JOIN products ON vir.product_id = products.id
		JOIN vendors ON vir.vendor_id = vendors.id
	`
	var rows pgx.Rows
	var err error

	if grnNo != "" {
		query += " WHERE grns.grn_number = $1 ORDER BY grns.created_at DESC"
		rows, err = db.DB.Query(ctx, query, grnNo)
	} else {
		query += " ORDER BY grns.created_at DESC"
		rows, err = db.DB.Query(ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grns []models.GRNResponse

	for rows.Next() {
		var grn models.GRNResponse
		err := rows.Scan(
			&grn.ID, &grn.GRNNumber, &grn.CreatedAt, &grn.ContainerQty, &grn.Quantity, &grn.Invoice,
			&grn.InvoiceDate, &grn.InvoiceImg, &grn.PackagingStatus, &grn.CreatedBy,
			&grn.VIRNumber, &grn.ProductName, &grn.VendorName,
		)
		if err != nil {
			return nil, err
		}
		grns = append(grns, grn)
	}

	return grns, nil
}
