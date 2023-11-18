package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Speech struct {
	AudioPath   string
	LanguageTag string
	APIEndpoint string
	Voice       string
}

func NewSpeech(audioPath string, languageTag string, apiEndpoint string, voice string) *Speech {
	return &Speech{
		AudioPath:   audioPath,
		LanguageTag: languageTag,
		APIEndpoint: apiEndpoint,
		Voice:       voice,
	}
}

/*
func (s *Speech) Speech(sentence string) error {
	speech := htgotts.Speech{Folder: s.AudioPath, Language: voices.English, Handler: &handlers.Native{}}
	err := speech.Speak(sentence)
	if err != nil {
		return err
	}
	return nil
}
*/

// Speech TODO: speak translated sentence using Own API (Piper)
func (s *Speech) Speech(sentence string) error {

	// Prepare the payload to send to the API
	payload := map[string]string{
		"text":  sentence,
		"voice": s.Voice,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.APIEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "audio/wav" {
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	audioData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Write the audio data to a file
	err = ioutil.WriteFile(s.AudioPath+"/speech.wav", audioData, 0644)
	if err != nil {
		return err
	}
	return nil
}
