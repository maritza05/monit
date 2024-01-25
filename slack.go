package main

import (
	"os"
	"sync"

	"github.com/slack-go/slack"
)

type Notifier interface {
	Notify(string)
}

type FileNotifier struct {
	mu       sync.Mutex
	filepath string
}

func NewFileNotifier(filepath string) *FileNotifier {
	return &FileNotifier{
		filepath: filepath,
	}
}

func (n *FileNotifier) Notify(message string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	file, err := os.Create(n.filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(message)
}

type SlackNotifier struct {
	tokenID   string
	channelID string
	client    *slack.Client
}

func NewSlackNotifier(tokenID string, channelID string) *SlackNotifier {
	return &SlackNotifier{
		tokenID:   tokenID,
		channelID: channelID,
		client:    slack.New(tokenID),
	}
}

func (s *SlackNotifier) Notify(message string) {
	_, _, err := s.client.PostMessage(s.channelID, slack.MsgOptionText(message, false))
	if err != nil {
		panic(err)
	}
}
