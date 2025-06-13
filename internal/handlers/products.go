package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllProductsData(c context.Context) ([]models.Product, error) {
	var products []models.Product
	var err error
	rows, err := db.DB.Query(c, "select id, name, quantity, userid FROM public.products")
	if err != nil {
		return []models.Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.UserId) // Adjust based on your struct
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return []models.Product{}, err
	}
	return products, nil
}

func AddProduct(c context.Context, req models.Product) error {
	// Normalize name to lowercase with no spaces
	nameHash := utils.ToLowerNoSpaces(req.Name)

	// Efficient duplicate check using SQL EXISTS
	checkQuery := `SELECT EXISTS (
		SELECT 1 FROM public.products WHERE namehash = $1
	)`
	var exists bool
	if err := db.DB.QueryRow(c, checkQuery, nameHash).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return errors.New("product already exists in the database")
	}

	// Insert product
	insertQuery := `INSERT INTO public.products (name, quantity, userid, namehash) VALUES ($1, $2, $3, $4)`
	_, err := db.DB.Exec(c, insertQuery, req.Name, req.Quantity, req.UserId, nameHash)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProductQuantityAndUser(ctx context.Context, productID int, quantity float64, userID int) error {
	query := `UPDATE public.products 
	          SET quantity = $1, userid = $2 
	          WHERE id = $3`

	result, err := db.DB.Exec(ctx, query, quantity, userID, productID)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
