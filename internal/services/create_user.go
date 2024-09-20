package services

import (
	"EasySplit/internal/database"
	"EasySplit/internal/models"
)

func InsertUser(user *models.User) error {
	query := `INSERT INTO users (id, first_name, last_name, username, email, phone, password, role, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := database.DB.Exec(query, user.ID, user.FirstName, user.LastName, user.Username, user.Email, user.Phone, user.Password, user.Role, user.CreatedAt)
	return err
}
