package internal

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	LogFilePath string
}

func NewLogger(logFilePath string) *Logger {
	return &Logger{
		LogFilePath: logFilePath,
	}
}

func (l *Logger) WriteSentenceToLogFile(sentence string) error {
	// Open the file in append mode. If the file doesn't exist, create it.
	file, err := os.OpenFile(l.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close log file: %v\n", err)
		}
	}(file)

	// Write the sentence to the file. DateTime_ + sentence + "\n"
	dateTimes := time.Now().Format("2006-01-02 15:04:05")
	_, err = file.WriteString(fmt.Sprintf("%s - %s\n", dateTimes, sentence))
	if err != nil {
		return fmt.Errorf("failed to write sentence to log file: %w", err)
	}
	return nil
}
