package services

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

// need discusion over auth
// this is auth Users table with 2 email
func GetAuthUsers(c *gin.Context) {

	rows, err := db.DB.Query(c, "select * FROM auth.Users")
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	defer rows.Close()
	var users []map[string]any
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

		var user map[string]any
		if err := mapstructure.Decode(rowMap, &user); err != nil {
			log.Printf("Failed to decode: %v", err)
			continue
		}

		users = append(users, user)
	}
	utils.SuccessWithData(c, users)

}
