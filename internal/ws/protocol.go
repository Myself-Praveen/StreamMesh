package ws

import (
	"encoding/json"
	"time"
)

// MessageType defines the type of the WebSocket message
type MessageType string

const (
	TypePublish     MessageType = "publish"
	TypeSubscribe   MessageType = "subscribe"
	TypeUnsubscribe MessageType = "unsubscribe"
	TypeHistory     MessageType = "history"
	TypeAck         MessageType = "ack"
	TypeError       MessageType = "error"
)

// Envelope is the standard JSON frame for all WebSocket messages
type Envelope struct {
	Type    MessageType     `json:"type"`
	Topic   string          `json:"topic,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"` // Message ID for acks
	TS      int64           `json:"ts,omitempty"` // Unix timestamp
}

// ParseMessage parses a raw byte array into an Envelope
func ParseMessage(data []byte) (*Envelope, error) {
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return &env, nil
}

// NewMessage creates a new outgoing message envelope
func NewMessage(msgType MessageType, topic string, payload interface{}, id string) (*Envelope, error) {
	var rawPayload json.RawMessage
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		rawPayload = b
	}

	return &Envelope{
		Type:    msgType,
		Topic:   topic,
		Payload: rawPayload,
		ID:      id,
		TS:      time.Now().UnixNano(),
	}, nil
}

// Encode serializes the Envelope to JSON bytes
func (e *Envelope) Encode() ([]byte, error) {
	return json.Marshal(e)
}
