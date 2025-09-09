package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func GetVendorIDByName(ctx context.Context, name string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(ctx, `SELECT id FROM vendors WHERE name = $1`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetProductIDByName(ctx context.Context, name string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(ctx, `SELECT id FROM products WHERE name = $1`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getNextVIRNumber(ctx context.Context, createdAt time.Time) (string, error) {
	monthYear := createdAt.Format("012006") // MMYYYY

	prefix := "VIR-" + monthYear + "-"

	var lastVIRNumber string

	query := `
		SELECT vir_number FROM vir
		WHERE vir_number LIKE $1
		ORDER BY id DESC LIMIT 1
	`

	err := db.DB.QueryRow(ctx, query, prefix+"%").Scan(&lastVIRNumber)
	if err != nil && err.Error() != "no rows in result set" {
		return "", err
	}

	var nextNumber int
	if lastVIRNumber == "" {
		nextNumber = 1
	} else {
		parts := strings.Split(lastVIRNumber, "-")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid VIR number format: %s", lastVIRNumber)
		}
		num, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("invalid number in VIR: %w", err)
		}
		nextNumber = num + 1
	}

	return fmt.Sprintf("%s%03d", prefix, nextNumber), nil
}

func InsertVIRHandler(ctx context.Context, req models.CreateVIRRequest) (*models.VIR, error) {
	vendorID, err := GetVendorIDByName(ctx, req.Vendor)
	if err != nil {
		return nil, errors.New("invalid vendor")
	}

	productID, err := GetProductIDByName(ctx, req.Product)
	if err != nil {
		return nil, errors.New("invalid product")
	}

	createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
	if err != nil {
		return nil, errors.New("invalid createdAt timestamp")
	}

	checklistJSONBytes, err := json.Marshal(req.Checklist)
	if err != nil {
		return nil, errors.New("invalid checklist format")
	}
	checklistJSON := string(checklistJSONBytes)

	virNumber, err := getNextVIRNumber(ctx, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate VIR number: %v", err)
	}

	vir := &models.VIR{
		CreatedAt: createdAt,
		VendorID:  vendorID,
		ProductID: productID,
		Checklist: req.Checklist,
		Remarks:   req.Remarks,
		CreatedBy: req.CreatedBy.Data.ID,
		Status:    "in-progress",
		VIRNumber: virNumber,
	}

	query := `
	INSERT INTO vir (created_at, vendor_id, product_id, checklist, remarks, created_by, status, vir_number)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`

	err = db.DB.QueryRow(ctx, query,
		vir.CreatedAt,
		vir.VendorID,
		vir.ProductID,
		checklistJSON,
		vir.Remarks,
		vir.CreatedBy,
		vir.Status,
		vir.VIRNumber,
	).Scan(&vir.ID)

	if err != nil {
		return nil, err
	}

	return vir, nil
}

func GetAllVIR(ctx context.Context) ([]models.VIR, error) {
	rows, err := db.DB.Query(ctx, `
	SELECT
		vi.id,
		vi.created_at,
		vi.vendor_id,
		vi.product_id,
		vi.checklist,
		vi.remarks,
		vi.created_by,
		vi.checked_by,
		vi.checked_at,
		vi.status,
		vi.vir_number,
		ven.name  AS vendor_name,
		prod.name AS product_name
	FROM vir vi
	LEFT JOIN vendors  ven  ON ven.id  = vi.vendor_id
	LEFT JOIN products prod ON prod.id = vi.product_id
	ORDER BY vi.created_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var virs []models.VIR

	for rows.Next() {
		var vir models.VIR
		var checklistJSON string
		var vendorNameNS, productNameNS sql.NullString

		err := rows.Scan(&vir.ID, &vir.CreatedAt, &vir.VendorID, &vir.ProductID, &checklistJSON,
			&vir.Remarks, &vir.CreatedBy, &vir.CheckedBy, &vir.CheckedAt, &vir.Status, &vir.VIRNumber, &vendorNameNS,
			&productNameNS)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(checklistJSON), &vir.Checklist); err != nil {
			return nil, errors.New("failed to parse checklist")
		}

		if vendorNameNS.Valid {
			vir.VendorName = &vendorNameNS.String
		}
		if productNameNS.Valid {
			vir.ProductName = &productNameNS.String
		}
		virs = append(virs, vir)
	}

	return virs, nil
}

func GetVIRByNumber(ctx context.Context, virNumber string) (*models.VIR, error) {
	var vir models.VIR
	var checklistJSON string
	var vendorNameNS, productNameNS sql.NullString

	query := `
		SELECT
			vi.id,
			vi.created_at,
			vi.vendor_id,
			vi.product_id,
			vi.checklist,
			vi.remarks,
			vi.created_by,
			vi.checked_by,
			vi.checked_at,
			vi.status,
			vi.vir_number,
			ven.name  AS vendor_name,
			prod.name AS product_name
		FROM vir vi
		LEFT JOIN vendors  ven  ON ven.id  = vi.vendor_id
		LEFT JOIN products prod ON prod.id = vi.product_id
		WHERE vi.vir_number = $1
	`

	err := db.DB.QueryRow(ctx, query, virNumber).Scan(
		&vir.ID, &vir.CreatedAt, &vir.VendorID, &vir.ProductID, &checklistJSON,
		&vir.Remarks, &vir.CreatedBy, &vir.CheckedBy, &vir.CheckedAt, &vir.Status, &vir.VIRNumber, &vendorNameNS,
		&productNameNS,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(checklistJSON), &vir.Checklist); err != nil {
		return nil, errors.New("failed to parse checklist")
	}

	if vendorNameNS.Valid {
		vir.VendorName = &vendorNameNS.String
	}
	if productNameNS.Valid {
		vir.ProductName = &productNameNS.String
	}

	return &vir, nil
}

func UpdateVIRVerification(ctx context.Context, virNumber string, req models.VerifyVIRRequest) error {
	checkedAt, err := time.Parse(time.RFC3339, req.CheckedAt)
	if err != nil {
		return errors.New("invalid checkedAt timestamp")
	}

	query := `
		UPDATE vir
		SET checked_by = $1,
			checked_at = $2,
			status = $3
		WHERE vir_number = $4
	`

	tag, err := db.DB.Exec(ctx, query, req.CheckedBy.Data.ID, checkedAt, "completed", virNumber)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no VIR found with number %s", virNumber)
	}

	return nil
}
