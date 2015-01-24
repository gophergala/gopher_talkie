package common

import (
	"time"
)

type Message struct {
	MessageID string        `json:"id"`
	From      *User         `json:"from"`
	To        *User         `json:"to"`
	Content   []byte        `json:"content,omitempty"`
	Timestamp time.Time     `json:"time"`
	Duration  time.Duration `json:"duration"`
	Played    bool          `json:"played"`
}
