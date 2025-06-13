package services

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllProducts(c *gin.Context) {
	rows, err := db.DB.Query(c, "select id, name, quantity, userid FROM public.products")
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	defer rows.Close()
	var products []models.Product
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
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, products)
}

func InsertProduct(c *gin.Context) {
	var req models.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	id, _ := c.Get("id")
	req.UserId = id.(int)
	err := db.DB.QueryRow(c, `insert into public.products(name, quantity, userid) Values ($1,$2,$3) returning id`, req.Name, req.Quantity, req.UserId).Scan(&req.ID)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	log.Println("--data inserted--", req.ID)
}
