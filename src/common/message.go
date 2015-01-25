package common

import (
	"time"
)

type Message struct {
	MessageID int64         `json:"id"`
	From      *User         `json:"from"`
	To        *User         `json:"to"`
	Content   []byte        `json:"content,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	Duration  time.Duration `json:"duration"`
	Played    bool          `json:"played"`
	Path      string        `json:"path"`
	RemoteURL string        `json:"remote_url"`
}

func NewMessage(from, to *User) *Message {
	return &Message{
		From: from,
		To:   to,
	}
}
