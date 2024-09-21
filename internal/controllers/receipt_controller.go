package controllers

import (
	"EasySplit/internal/services"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func UploadMedia(ctx echo.Context) error {
	email := ctx.Get("user_email").(string)
	var imageURL, audioURL, detectedText, audioTranscription string

	if imageFile, err := ctx.FormFile("image"); err == nil {
		src, err := imageFile.Open()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Could not open image file",
				"status":  "error",
			})
		}
		defer src.Close()

		imageURL, err = processAndUploadImageToS3(imageFile.Filename, src)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Error uploading image to S3",
				"status":  "error",
			})
		}

		detectedText, err = DetectTextFromImage(imageURL)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Error detecting text from image",
				"status":  "error",
			})
		}
	}

	if audioFile, err := ctx.FormFile("audio"); err == nil {
		src, err := audioFile.Open()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Could not open audio file",
				"status":  "error",
			})
		}
		defer src.Close()

		audioData, err := io.ReadAll(src)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Could not read audio file",
				"status":  "error",
			})
		}

		audioURL, err = uploadToS3(audioFile.Filename, bytes.NewReader(audioData))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Error uploading audio to S3",
				"status":  "error",
			})
		}

		audioTranscription, err = WhisperTranscriptionFromAudio(audioURL)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Error transcribing audio",
				"status":  "error",
			})
		}
	}

	if err := services.StoreMediaURL(email, imageURL, audioURL); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error storing media URLs",
			"status":  "error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Media uploaded successfully",
		"imageURL":      imageURL,
		"audioURL":      audioURL,
		"detectedText":  detectedText,
		"transcription": audioTranscription,
		"status":        "success",
	})
}

func DetectTextFromImage(imageURL string) (string, error) {
	imageData, err := fetchImageFromURL(imageURL)
	if err != nil {
		return "", fmt.Errorf("error fetching image from URL: %v", err)
	}

	detectedText, err := detectTextFromImage(imageData)
	if err != nil {
		return "", fmt.Errorf("error detecting text from image: %v", err)
	}

	return detectedText, nil
}

func WhisperTranscriptionFromAudio(s3Url string) (string, error) {
	audioFile, err := fetchAudioFile(s3Url)
	if err != nil {
		return "", fmt.Errorf("error fetching audio file: %v", err)
	}
	defer audioFile.Close()

	transcription, err := transcribeWithWhisper(audioFile)
	if err != nil {
		return "", fmt.Errorf("error transcribing audio: %v", err)
	}

	return transcription, nil
}

func processAndUploadImageToS3(fileName string, file multipart.File) (string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Unable to decode image: %v", err)
		return "", err
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		log.Printf("Unable to encode image: %v", err)
		return "", err
	}

	return uploadToS3(fileName, bytes.NewReader(buf.Bytes()))
}

func uploadToS3(fileName string, file *bytes.Reader) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("Unable to load SDK config: %v", err)
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
