package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrEmptySenderID is returned when a message is missing a SenderID.
	ErrEmptySenderID = errors.New("validation error: SenderID cannot be empty")
	// ErrEmptyReceiverID is returned when a message is missing a ReceiverID.
	ErrEmptyReceiverID = errors.New("validation error: ReceiverID cannot be empty")
	// ErrEmptyPayload is returned when a message is missing a Payload.
	ErrEmptyPayload = errors.New("validation error: Payload cannot be empty")
	// ErrEmptyTimestamp is returned when a message is missing a Timestamp.
	ErrEmptyTimestamp = errors.New("validation error: Timestamp cannot be empty")
)

type Message struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	Payload    string `json:"payload"`
	Timestamp  int64  `json:"timestamp"`
}

func NewMessage(senderID, receiverID, payload string) (*Message, error) {
	if senderID == "" {
		return nil, ErrEmptySenderID
	}
	if receiverID == "" {
		return nil, ErrEmptyReceiverID
	}
	if payload == "" {
		return nil, ErrEmptyPayload
	}

	return &Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Payload:    payload,
		Timestamp:  time.Now().UTC().Unix(),
	}, nil
}

func ParseFromJSON(s string) (*Message, error) {
	var m Message
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("error while unmarshalling message: %w", err)
	}

	if m.SenderID == "" {
		return nil, ErrEmptySenderID
	}
	if m.ReceiverID == "" {
		return nil, ErrEmptyReceiverID
	}
	if m.Payload == "" {
		return nil, ErrEmptyPayload
	}
	if m.Timestamp == 0 {
		return nil, ErrEmptyTimestamp
	}

	return &m, nil
}
