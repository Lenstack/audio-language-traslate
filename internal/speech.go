package internal

import (
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

type Speech struct {
	AudioPath string
}

func NewSpeech(audioPath string) *Speech {
	return &Speech{AudioPath: audioPath}
}

func (s *Speech) Speech(sentence string) error {
	speech := htgotts.Speech{Folder: s.AudioPath, Language: voices.English, Handler: &handlers.Native{}}
	err := speech.Speak(sentence)
	if err != nil {
		return err
	}
	return nil
}
