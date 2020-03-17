package messenger

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Messenger is slack incoming webhook client.
type Messenger struct {
	WebhookURL string
}

// New returns a new Messenger
func New(webhookURL string) *Messenger {
	return &Messenger{webhookURL}
}

// Send sends given message with slack incoming webhook.
func (msgr *Messenger) Send(m *Message) (*http.Response, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		msgr.WebhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res, nil
}

// Message is slack message.
type Message struct {
	Text   string        `json:"text"`
	Blocks []interface{} `json:"blocks"`
}

// Section is section block of slack message.
type Section struct {
	Type string  `json:"type"`
	Text Content `json:"text"`
}

// Context is context block of slack message.
type Context struct {
	Type     string    `json:"type"`
	Elements []Content `json:"elements"`
}

// Divider is divider block of slack message.
type Divider struct {
	Type string `json:"type"`
}

// Content is common structure of block content.
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}