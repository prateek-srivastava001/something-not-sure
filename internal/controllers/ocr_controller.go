package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type ImageURLRequest struct {
	ImageURL string `json:"image_url"`
}

type TextResponse struct {
	Text string `json:"text"`
}

func DetectText(c echo.Context) error {
	var req ImageURLRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.ImageURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No image URL provided"})
	}

	imageData, err := fetchImageFromURL(req.ImageURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	detectedText, err := detectTextFromImage(imageData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, TextResponse{Text: detectedText})
}

func fetchImageFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image, status code: %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}
	fmt.Println("fetched image")
	return imageData, nil
}

func detectTextFromImage(imageData []byte) (string, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return "", fmt.Errorf("failed to create vision client: %v", err)
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to create image: %v", err)
	}

	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return "", fmt.Errorf("failed to detect text: %v", err)
	}

	if annotation == nil || annotation.Text == "" {
		return "No text detected.", nil
	}

	fullText := annotation.Text
	pattern := regexp.MustCompile(`(o{3,}|O{3,}|0{3,})`)
	cleanedText := pattern.ReplaceAllString(fullText, "\n")

	return cleanedText, nil
}
