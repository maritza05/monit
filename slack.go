package main

import (
	"github.com/slack-go/slack"
)

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
