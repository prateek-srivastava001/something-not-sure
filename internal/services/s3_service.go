package services

import (
	"EasySplit/internal/database"
)

func StoreMediaURL(mediaID, email, imageURL, audioURL, billText, audioText string) error {
	query := `INSERT INTO user_media (id, user_email, image_url, audio_url, image_parsed, audio_parsed) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := database.DB.Exec(query, mediaID, email, imageURL, audioURL, billText, audioText)
	return err
}
