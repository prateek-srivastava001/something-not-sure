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

func UploadMedia(ctx echo.Context) error {
	email := ctx.Get("user_email").(string)
	var imageURL, audioURL string

	// Handle image upload
	if imageFile, err := ctx.FormFile("image"); err == nil {
		src, err := imageFile.Open()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Could not open image file",
				"status":  "error",
			})
		}
		defer src.Close()

		imageURL, err = uploadToS3(imageFile.Filename, src)
		if err != nil {
			log.Printf("Error uploading image to S3: %v", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Error uploading image to S3",
				"status":  "error",
			})
		}
	}

	// Handle audio upload
	if audioFile, err := ctx.FormFile("audio"); err == nil {
		src, err := audioFile.Open()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Could not open audio file",
				"status":  "error",
			})
		}
		defer src.Close()

		audioURL, err = uploadToS3(audioFile.Filename, src)
		if err != nil {
			log.Printf("Error uploading audio to S3: %v", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Error uploading audio to S3",
				"status":  "error",
			})
		}
	}

	// Store the media URLs in the database
	if err := services.StoreMediaURL(email, imageURL, audioURL); err != nil {
		log.Printf("Error storing media URLs in DB: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error storing media URLs",
			"status":  "error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Media uploaded successfully",
		"imageURL": imageURL,
		"audioURL": audioURL,
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
