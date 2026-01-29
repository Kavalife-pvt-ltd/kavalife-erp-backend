package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllUsers(c context.Context) ([]models.User, error) {
	var users []models.User
	var err error
	rows, err := db.DB.Query(c, "select id, username, created_at, role, department, phone_num, name, email FROM public.users")
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.User
		err := rows.Scan(&p.ID, &p.Username, &p.Created_at, &p.Role, &p.Department, &p.Phone_Num, &p.Name, &p.Email)
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		users = append(users, p)
	}
	utils.SortByID(users, true)
	if err = rows.Err(); err != nil {
		return []models.User{}, err
	}
	return users, err
}
func AllNewUser(c context.Context) ([]models.NewUser, error) {
	var users []models.NewUser
	var err error
	rows, err := db.DB.Query(c, "select id, username, mob_number,email, name FROM new_user")
	if err != nil {
		return []models.NewUser{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.NewUser
		err := rows.Scan(&p.ID, &p.Username, &p.Mobile, &p.Email, &p.Name)
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		users = append(users, p)
	}
	utils.SortByID(users, true)
	if err = rows.Err(); err != nil {
		return []models.NewUser{}, err
	}
	return users, err
}
func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, role, department, department_role, password, phone_num FROM public.users WHERE username = $1`

	var user models.User
	err := db.DB.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Role, &user.Department,
		&user.DepartmentRole, &user.Password, &user.Phone_Num,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no rows found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserByPhoneNum(ctx context.Context, phoneNo string) (*models.User, error) {
	query := `SELECT id, username, role, department, department_role, password, phone_num FROM public.users WHERE phone_num = $1`

	var user models.User
	err := db.DB.QueryRow(ctx, query, phoneNo).Scan(
		&user.ID, &user.Username, &user.Role, &user.Department,
		&user.DepartmentRole, &user.Password, &user.Phone_Num,
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

	query := `SELECT id, username, name, role , department FROM public.users WHERE username = $1`
	err := db.DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Name, &user.Role, &user.Department)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUsers(ctx context.Context, newUser models.NewUser, role string, department string, adminId int) (int, error) {
	createdAt := time.Now().UTC()
	var id int
	query := `INSERT into users (username, email, password, phone_num, name, role, department, created_at,approved_by)
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9) 
	returning id`
	err := db.DB.QueryRow(ctx, query, newUser.Username, newUser.Email, newUser.Password, newUser.Mobile, newUser.Name, role, department, createdAt, adminId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func FindNewUserByID(ctx context.Context, id int) (*models.NewUser, error) {
	var newUser models.NewUser
	query := `SELECT id, username, email, password, mob_number, name FROM new_user WHERE id = $1`
	err := db.DB.QueryRow(ctx, query, id).Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.Mobile, &newUser.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &newUser, nil
}
func InsertNewUser(ctx context.Context, newUser models.NewUser) (int, error) {
	var id int

	query := `INSERT into new_user (username, email, password, mob_number, name)
	values ($1,$2,$3,$4,$5) 
	returning id`
	err := db.DB.QueryRow(ctx, query, newUser.Username, newUser.Email, newUser.Password, newUser.Mobile, newUser.Name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func DeleteNewUser(ctx context.Context, id int) error {
	query := `DELETE FROM new_user WHERE id = $1`

	cmdTag, err := db.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
