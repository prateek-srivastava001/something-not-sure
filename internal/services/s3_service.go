package services

import (
	"EasySplit/internal/database"
)

func StoreMediaURL(email, imageURL, audioURL, billText, audioText string) error {
	query := `INSERT INTO user_media (user_email, image_url, audio_url, image_parsed, audio_parsed) VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query, email, imageURL, audioURL, billText, audioText)
	return err
}
