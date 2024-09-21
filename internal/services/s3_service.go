package services

import (
	"EasySplit/internal/database"
)

func StoreMediaURL(email, imageURL, audioURL string) error {
	query := `INSERT INTO user_media (user_email, image_url, audio_url) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, email, imageURL, audioURL)
	return err
}

func GetImageURL(email string) (string, error) {
	var imageURL string
	query := `SELECT image_url FROM user_media WHERE user_email = $1`
	err := database.DB.QueryRow(query, email).Scan(&imageURL)
	return imageURL, err
}

func GetAudioURL(email string) (string, error) {
	var audioURL string
	query := `SELECT audio_url FROM user_media WHERE user_email = $1`
	err := database.DB.QueryRow(query, email).Scan(&audioURL)
	return audioURL, err
}
