package event

import (
	"encoding/json"
	"sync"
)

// Message represents a pub-sub message with a topic and payload.
type Message struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}

// Client represents a WebSocket client with subscriptions and an outgoing message channel.
type Client struct {
	Send   chan []byte
	Topics map[string]bool
	mu     sync.Mutex
}

// NewClient creates a new client with a buffered send channel.
func NewClient() *Client {
	return &Client{
		Send:   make(chan []byte, 64),
		Topics: make(map[string]bool),
	}
}

// SubscribeTo adds a topic to the client's subscriptions.
func (c *Client) SubscribeTo(topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Topics[topic] = true
}

// UnsubscribeFrom removes a topic from the client's subscriptions.
func (c *Client) UnsubscribeFrom(topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Topics, topic)
}

// IsSubscribed checks if the client is subscribed to a topic.
func (c *Client) IsSubscribed(topic string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Topics[topic]
}

// Broker manages client connections and message distribution.
type Broker struct {
	Clients    map[*Client]bool
	NewClients chan *Client
	Defunct    chan *Client
	Messages   chan Message
}

// NewBroker creates a new message broker.
func NewBroker() *Broker {
	return &Broker{
		Clients:    make(map[*Client]bool),
		NewClients: make(chan *Client),
		Defunct:    make(chan *Client),
		Messages:   make(chan Message, 256),
	}
}

// Start runs the broker's event loop in a goroutine.
func (b *Broker) Start() {
	go func() {
		for {
			select {
			case client := <-b.NewClients:
				b.Clients[client] = true

			case client := <-b.Defunct:
				if b.Clients[client] {
					delete(b.Clients, client)
					close(client.Send)
				}

			case message := <-b.Messages:
				data, err := json.Marshal(message)
				if err != nil {
					continue
				}

				for client := range b.Clients {
					if client.IsSubscribed(message.Topic) {
						select {
						case client.Send <- data:
							// Message sent successfully
						default:
							// Client's send buffer is full, remove and close
							delete(b.Clients, client)
							close(client.Send)
						}
					}
				}
			}
		}
	}()
}

// Publish sends a message to all subscribers of a topic.
func (b *Broker) Publish(topic string, payload any) {
	b.Messages <- Message{
		Topic:   topic,
		Payload: payload,
	}
}
