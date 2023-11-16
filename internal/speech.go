package internal

import (
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

type Speech struct {
	AudioPath   string
	LanguageTag string
	APIEndpoint string
}

func NewSpeech(audioPath string, languageTag string, apiEndpoint string) *Speech {
	return &Speech{
		AudioPath:   audioPath,
		LanguageTag: languageTag,
		APIEndpoint: apiEndpoint,
	}
}

func (s *Speech) Speech(sentence string) error {
	speech := htgotts.Speech{Folder: s.AudioPath, Language: voices.English, Handler: &handlers.Native{}}
	err := speech.Speak(sentence)
	if err != nil {
		return err
	}
	return nil
}
