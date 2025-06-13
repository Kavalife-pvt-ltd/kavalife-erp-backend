package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func AllUsers(c *gin.Context) ([]models.User, error) {
	var users []models.User
	var err error
	rows, err := db.DB.Query(c, "select id, username, created_at, role, phone_num FROM public.users")
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.User
		err := rows.Scan(&p.ID, &p.Username, &p.Created_at, &p.Role, &p.Phone_Num)
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		users = append(users, p)

	}
	if err = rows.Err(); err != nil {
		return []models.User{}, err
	}
	return users, err
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, role, password, phone_num FROM public.users WHERE username = $1`

	var user models.User
	err := db.DB.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Role, &user.Password, &user.Phone_Num,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no rows found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetLoggedUserByUsername(ctx context.Context, username string) (*models.LoggedUserDetails, error) {
	var user models.LoggedUserDetails

	query := `SELECT id, username, role FROM public.users WHERE username = $1`
	err := db.DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Role)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
