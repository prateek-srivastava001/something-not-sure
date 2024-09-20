package services

import (
	"EasySplit/internal/database"
	"EasySplit/internal/models"
)

func GetUserByID(userID string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, last_name, username, email, phone, password, created_at 
		FROM users 
		WHERE id = $1
	`
	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, username = $3, email = $4, phone = $5, password = $6 
		WHERE id = $7
	`
	_, err := database.DB.Exec(query, user.FirstName, user.LastName, user.Username, user.Email, user.Phone, user.Password, user.ID)
	return err
}
