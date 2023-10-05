package internal

import (
	"log"
	"os"
	"time"
)

// Record do FFmpeg
type Record struct {
	Command        *Command
	AudioPath      string
	FileName       string
	SilenceTimeout time.Duration
	Threshold      int
}

func NewRecord(ffmpegPath string, audioPath string, fileName string, silenceTimeout time.Duration, threshold int) *Record {
	return &Record{
		Command:        NewCommand(ffmpegPath),
		AudioPath:      audioPath,
		FileName:       fileName,
		SilenceTimeout: silenceTimeout,
		Threshold:      threshold,
	}
}

// Start TODO: start recording records
func (r *Record) Start() error {
	log.Printf("Start recording records to %s\n", r.AudioPath)

	err := r.CreateRecordsFolder()
	if err != nil {
		return err
	}

	cmd := r.Command.ExecuteFFCommand("ffmpeg", []string{"-f", "dshow", "-i", "audio=" + "FrontMic (Realtek(R) Audio)", "-ac", "1", "-ar", "44100", "-af", "silencedetect=noise=-30dB:d=0.5", "-y", r.AudioPath + "/" + r.FileName})
	err = cmd.Start()
	if err != nil {
		log.Println("Error starting records recorder:", err)
		return err
	}
	return nil
}

// Stop TODO: stop recording records
func (r *Record) Stop() error {
	log.Printf("Stop recording records to %s\n", r.AudioPath)

	err := r.Command.cmd.Process.Kill()
	if err != nil {
		log.Printf("Error stopping records recorder: %v\n", err)
		return err
	}
	return nil
}

// Silence TODO: detect silence in stream records (silenceTimeout) and save records in a file (audioName)_(timestamp).mp3
func (r *Record) Silence() (bool, error) {
	log.Printf("Detect silence in records stream\n")

	cmd := r.Command.ExecuteFFCommand("ffmpeg", []string{})
	err := cmd.Run()
	if err != nil {
		log.Printf("Error detecting silence in records stream: %v\n", err)
		return false, err
	}
	return true, nil
}

// Save TODO: save records in a file
func (r *Record) Save() error {
	log.Println("Save records to", r.AudioPath)

	cmd := r.Command.ExecuteFFCommand("ffmpeg", []string{})
	err := cmd.Run()
	if err != nil {
		log.Printf("Error saving records: %v\n", err)
		return err
	}
	return nil
}

// Delete TODO: delete records file
func (r *Record) Delete() error {
	log.Printf("Delete records file %s\n", r.AudioPath)

	err := os.Remove(r.AudioPath)
	if err != nil {
		log.Fatalf("Error deleting records file: %v\n", err)
		return err
	}

	log.Printf("Audio file %s deleted\n", r.AudioPath)
	return nil
}

// CreateRecordsFolder TODO: create records folder
func (r *Record) CreateRecordsFolder() error {
	log.Printf("Create records folder %s\n", r.AudioPath)
	if _, err := os.Stat(r.AudioPath); os.IsNotExist(err) {
		err := os.Mkdir(r.AudioPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
