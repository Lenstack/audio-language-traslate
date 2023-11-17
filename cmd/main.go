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
	logFilePath    = "logs/log.txt"
	silenceTimeout = 2 * time.Second
	threshold      = -30 // dB
	languageTag    = "en"
	apiEndpoint    = "http://192.168.31.20:9000/asr"
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

	// Convert audio to text using Own API (Whisper)
	translate := internal.NewTranslate(audioPath, languageTag, apiEndpoint)
	sentence, err := translate.Translate(fileName)
	if err != nil {
		log.Fatalf("Error translating records file: %v", err)
	}

	log.Println("Speaking sentence: ", sentence)

	// Speak translated sentence using Own API (Bark)
	speech := internal.NewSpeech(audioPath, languageTag, apiEndpoint)
	err = speech.Speech(sentence)
	if err != nil {
		log.Fatalf("Error speaking sentence: %v", err)
	}

	log.Println("Writing translated sentence to log file.")

	// Write translated sentence to log file (for debugging)
	logger := internal.NewLogger(logFilePath)
	err = logger.WriteSentenceToLogFile(sentence)
	if err != nil {
		log.Fatalf("Error writing translated sentence to log file: %v", err)
	}
}
