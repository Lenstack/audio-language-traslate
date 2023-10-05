package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
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
		"task":            "translate",
		"language":        t.LanguageTag,
		"initial_prompt":  "",
		"encode":          "true",
		"output":          "txt",
		"word_timestamps": "false",
	}

	// Add query parameters to the API url
	apiUrl, err := url.Parse(t.APIEndpoint)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %v", err)
	}

	// Add query parameters to the API url
	query := apiUrl.Query()
	for key, val := range queryParams {
		query.Set(key, val)
	}

	// Encode the query parameters and update the API url
	apiUrl.RawQuery = query.Encode()

	// Create a buffer to store the multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the audio file to the multipart form
	audioFile, err := os.Open(t.AudioPath + "/" + fileName)
	if err != nil {
		return "", fmt.Errorf("error opening audio file: %v", err)
	}

	// Close the audio file when the function returns
	defer func(audioFile *os.File) {
		err := audioFile.Close()
		if err != nil {
			log.Printf("Error closing audio file: %v\n", err)
			return
		}
	}(audioFile)

	// Create a form field for the audio file
	part, err := writer.CreateFormFile("audio_file", fileName)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}

	// Copy the audio file data to the form field
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
	req, err := http.NewRequest("POST", apiUrl.String(), body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type header, as well as the multipart boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	// Close the response body when the function returns
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read and return the response from the API endpoint as a string
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	return string(responseBody), nil
}
