package api

import (
	"time"
)

type MessageResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty" swaggerignore:"true"`
}

type JoinRoomMessage struct {
	RoomHash string `json:"roomHash"`
}

// MessageMessage represents a message sent within a room.
type MessageMessage struct {
	RoomHash string        `json:"roomHash"`
	Sender   *UserResponse `json:"sender"`
	Message  string        `json:"message"`
	SendAt   time.Time     `json:"created_at" swaggerignore:"true"`
}

// TypingMessage represents a typing status update.
type TypingMessage struct {
	RoomHash string        `json:"roomHash"`
	Sender   *UserResponse `json:"sender"`
	TyperID  int           `json:"senderID"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
