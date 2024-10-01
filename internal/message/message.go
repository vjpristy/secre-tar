package message

import (
	"encoding/json"
	"time"
)

type Message struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func NewMessage(from, to, content string) *Message {
	return &Message{
		From:    from,
		To:      to,
		Content: content,
		Time:    time.Now(),
	}
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func Deserialize(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return &m, err
}
