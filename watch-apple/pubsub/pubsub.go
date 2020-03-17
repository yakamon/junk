package pubsub

// Message is message sent from PubSub
type Message struct {
	Data []byte `json:"data"`
}
