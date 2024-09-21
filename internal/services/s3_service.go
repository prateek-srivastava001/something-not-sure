package services

import (
	"EasySplit/internal/database"
)

func StoreImageURL(email, imageURL string) error {
	query := `INSERT INTO user_images (user_email, image_url) VALUES ($1, $2)`
	_, err := database.DB.Exec(query, email, imageURL)
	return err
}
