package main

import (
	"github.com/lenstack/audio-language-translate/internal"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	ffmpegPath     = "ffmpeg/bin"
	audioPath      = "records"
	silenceTimeout = 2 * time.Second
	threshold      = -30
	languageTag    = "en"
	apiEndpoint    = "http://192.168.31.20:9001/asr"
	fileName       = "records.wav"
)

func main() {
	recorder := internal.NewRecord(ffmpegPath, audioPath, fileName, silenceTimeout, threshold)

	// Start recording audio
	err := recorder.Start()
	if err != nil {
		log.Fatalf("Error starting records recorder: %v", err)
	}

	log.Println("Recording started. Press Ctrl+C to stop.")

	// Wait for interruption (Ctrl+C) to stop recording
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Stop recording
	err = recorder.Stop()
	if err != nil {
		log.Fatalf("Error stopping records recorder: %v", err)
	}

	log.Println("Recording stopped.")

	// Convert audio to text using Own API
	translate := internal.NewTranslate(audioPath, languageTag, apiEndpoint)
	text, err := translate.Translate(fileName)
	if err != nil {
		log.Fatalf("Error translating records file: %v", err)
	}

	log.Printf("Translated text: %s\n", text)
}
