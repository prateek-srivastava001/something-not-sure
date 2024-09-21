package controllers

import (
	"EasySplit/internal/services"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

// In UploadImage function
func UploadImage(ctx echo.Context) error {
	email := ctx.Get("user_email").(string)
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "File is required",
			"status":  "fail",
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Could not open file",
			"status":  "error",
		})
	}
	defer src.Close()

	s3URL, err := uploadToS3(file.Filename, src)
	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error uploading to S3",
			"status":  "error",
		})
	}

	// Check if audio URL is already stored
	audioURL, err := services.GetAudioURL(email)
	if err != nil {
		log.Printf("Error fetching audio URL from DB: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error fetching audio URL",
			"status":  "error",
		})
	}

	if err := services.StoreMediaURL(email, s3URL, audioURL); err != nil { // Store image URL
		log.Printf("Error storing image URL in DB: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error storing image URL",
			"status":  "error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Image uploaded successfully",
		"imageURL": s3URL,
		"status":   "success",
	})
}

// In UploadAudio function
func UploadAudio(ctx echo.Context) error {
	email := ctx.Get("user_email").(string)
	file, err := ctx.FormFile("audio")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "File is required",
			"status":  "fail",
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Could not open file",
			"status":  "error",
		})
	}
	defer src.Close()

	s3URL, err := uploadToS3(file.Filename, src)
	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error uploading to S3",
			"status":  "error",
		})
	}

	// Check if image URL is already stored
	imageURL, err := services.GetImageURL(email)
	if err != nil {
		log.Printf("Error fetching image URL from DB: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error fetching image URL",
			"status":  "error",
		})
	}

	if err := services.StoreMediaURL(email, imageURL, s3URL); err != nil {
		log.Printf("Error storing audio URL in DB: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error storing audio URL",
			"status":  "error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Audio uploaded successfully",
		"audioURL": s3URL,
		"status":   "success",
	})
}

func uploadToS3(fileName string, file multipart.File) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return "", err
	}

	s3Client := s3.NewFromConfig(cfg)

	newFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), fileName)

	bucketName := os.Getenv("S3_BUCKET_NAME")
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(newFileName),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("AWS_REGION"), newFileName)
	return s3URL, nil
}
