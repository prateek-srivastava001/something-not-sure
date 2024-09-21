package services

import (
	"EasySplit/internal/database"
)

func StoreMediaURL(email, imageURL, audioURL string) error {
	query := `INSERT INTO user_media (user_email, image_url, audio_url) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, email, imageURL, audioURL)
	return err
}
