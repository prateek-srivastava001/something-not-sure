package services

import (
	"EasySplit/internal/database"
	"EasySplit/internal/models"
)

func FindUserByEmail(emailOrUsername string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, last_name, username, email, phone, password, role 
		FROM users 
		WHERE email = $1 OR username = $2
	`
	err := database.DB.QueryRow(query, emailOrUsername, emailOrUsername).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
