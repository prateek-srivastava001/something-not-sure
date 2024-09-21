package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func WhisperTranscription(ctx echo.Context) error {
	s3Url := ctx.FormValue("s3_url")
	if s3Url == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "S3 URL is required",
			"status":  "fail",
		})
	}

	audioFile, err := fetchAudioFile(s3Url)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch audio from S3",
			"status":  "error",
		})
	}
	defer audioFile.Close()

	transcription, err := transcribeWithWhisper(audioFile)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to transcribe audio",
			"status":  "error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Transcription successful",
		"transcription": transcription,
		"status":        "success",
	})
}

func fetchAudioFile(s3Url string) (io.ReadCloser, error) {
	resp, err := http.Get(s3Url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file from S3, status: %s", resp.Status)
	}

	return resp.Body, nil
}

func transcribeWithWhisper(audioFile io.Reader) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API key is not set")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "audio.m4a")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, audioFile)
	if err != nil {
		return "", err
	}

	writer.WriteField("model", "whisper-1")

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Create the Whisper API request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/translations", &buf)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to transcribe audio: %s", string(body))
	}

	return string(body), nil
}
