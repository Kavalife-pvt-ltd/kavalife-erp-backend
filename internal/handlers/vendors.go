package handlers

import (
	"context"
	"log"

	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func AllVendorsData(c context.Context) ([]models.Vendors, error) {
	var products []models.Vendors
	var err error
	rows, err := db.DB.Query(c, "select * FROM public.vendors")
	if err != nil {
		return []models.Vendors{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Vendors
		err := rows.Scan(&p.ID, &p.Created_at, &p.GovId, &p.Name, &p.Status, &p.Type, &p.UpdatedBy, &p.Updated_at) // Adjust based on your struct
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return []models.Vendors{}, err
	}
	return products, nil
}
