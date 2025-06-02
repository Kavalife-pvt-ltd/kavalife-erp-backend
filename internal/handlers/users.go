package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func AllUsers(c *gin.Context) {
	rows, err := db.DB.Query(c, "select * FROM users")
	if err != nil {
		c.JSON(200, map[string]any{
			"error": err.Error(),
		})
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Printf("Failed to get values: %v", err)
			continue
		}
		fieldDescriptions := rows.FieldDescriptions()
		rowMap := make(map[string]interface{}, len(values))

		for i, val := range values {
			colName := string(fieldDescriptions[i].Name)
			rowMap[colName] = val
		}

		var user models.User
		if err := mapstructure.Decode(rowMap, &user); err != nil {
			log.Printf("Failed to decode: %v", err)
			continue
		}

		users = append(users, user)
	}
	c.JSON(200, users)

}
