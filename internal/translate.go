package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Translate struct {
	AudioPath   string
	LanguageTag string
	APIEndpoint string
}

func NewTranslate(audioPath string, languageTag string, apiEndpoint string) *Translate {
	return &Translate{
		AudioPath:   audioPath,
		LanguageTag: languageTag,
		APIEndpoint: apiEndpoint,
	}
}

// Translate TODO: translate records file to text using Own API (Whisper)
func (t *Translate) Translate(fileName string) (string, error) {
	log.Printf("Translate audio file %s to text\n", fileName)

	// Define query parameters for the API
	queryParams := map[string]string{
		"task":     "transcribe",
		"language": "en",
		//"initial_prompt":  "translate from " + t.LanguageTag + " to english",
		"encode": "true",
		"output": "txt",
		//"word_timestamps": "false",
	}

	// Create a buffer to store the multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add query parameters to the request
	for key, val := range queryParams {
		_ = writer.WriteField(key, val)
	}

	// Add the audio file to the multipart form
	audioFile, err := os.Open(t.AudioPath + "/" + fileName)
	if err != nil {
		return "", fmt.Errorf("error opening audio file: %v", err)
	}
	defer func(audioFile *os.File) {
		err := audioFile.Close()
		if err != nil {
			log.Printf("Error closing audio file: %v\n", err)
			return
		}
	}(audioFile)

	part, err := writer.CreateFormFile("audio_file", fileName)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(part, audioFile)
	if err != nil {
		return "", fmt.Errorf("error copying audio file data: %v", err)
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %v", err)
	}

	// Create a POST request with the multipart body and content type
	req, err := http.NewRequest("POST", t.APIEndpoint, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read and return the response from the API
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	return string(responseBody), nil
}
