package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllVendorsData(c context.Context) ([]models.Vendors, error) {
	var products []models.Vendors
	var err error
	rows, err := db.DB.Query(c, "select id,name,created_at,status,gov_id,type,updated_by,updated_at FROM public.vendors")
	if err != nil {
		return []models.Vendors{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Vendors
		err := rows.Scan(&p.ID, &p.Name, &p.Created_at, &p.Status, &p.GovId, &p.Type, &p.UpdatedBy, &p.Updated_at) // Adjust based on your struct
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		products = append(products, p)
	}
	utils.SortByID(products, true)
	if err = rows.Err(); err != nil {
		return []models.Vendors{}, err
	}
	return products, nil
}

func AddVendor(c *gin.Context, vendor models.Vendors, userID int) error {

	// Check if vendor with the same GovId already exists
	var exists bool
	err := db.DB.QueryRow(c, "SELECT EXISTS(SELECT 1 FROM public.vendors WHERE gov_id = $1 and type = $2)", vendor.GovId, vendor.Type).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("vendor with GovId '%s' already exists with type '%s'", vendor.GovId, vendor.Type)
	}

	// Prepare vendor data
	now := time.Now()
	vendor.UpdatedBy = userID
	vendor.Created_at = now
	vendor.Updated_at.Time = now
	vendor.Updated_at.Valid = true

	// Insert vendor
	_, err = db.DB.Exec(c, `
		INSERT INTO public.vendors (name, gov_id, status, type, created_at, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		vendor.Name,
		vendor.GovId,
		vendor.Status,
		vendor.Type,
		vendor.Created_at,
		vendor.UpdatedBy,
		vendor.Updated_at,
	)
	if err != nil {
		return fmt.Errorf("failed to insert vendor: %w", err)
	}

	return nil
}
// func UpdateVendor(ctx context.Context, productID int, quantity float64, userID int) error {
// 	query := `UPDATE public.vendors 
// 	          SET quantity = $1, userid = $2 
// 	          WHERE id = $3`

// 	result, err := db.DB.Exec(ctx, query, quantity, userID, productID)
// 	if err != nil {
// 		return err
// 	}
// 	rowsAffected := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("product not found")
// 	}
// 	return nil
// }
