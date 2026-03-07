package handler

import (
	"net/http"

	"Tournament/internal/adapters/driven/event"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// upgrader configures the WebSocket upgrader.
// TODO: For production, implement proper origin checking.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// clientCommand represents a command sent by a WebSocket client.
type clientCommand struct {
	Action string `json:"action"`
	Topic  string `json:"topic"`
}

// EventHandler handles WebSocket connections for the event broker.
type EventHandler struct {
	Broker *event.Broker
}

// NewEventHandler creates a new event handler with the given broker.
func NewEventHandler(broker *event.Broker) *EventHandler {
	return &EventHandler{
		Broker: broker,
	}
}

// RegisterRoutes registers the event handler routes.
func (h *EventHandler) RegisterRoutes(router chi.Router) {
	router.Get("/events/ws", h.HandleWebSocket)
}

// HandleWebSocket upgrades HTTP connections to WebSocket and manages client lifecycle.
func (h *EventHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := event.NewClient()
	h.Broker.NewClients <- client

	// Write pump: send messages from client.Send to WebSocket
	go func() {
		defer conn.Close()
		for message := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				break
			}
		}
	}()

	// Read pump: read commands from WebSocket and handle subscriptions
	defer func() {
		h.Broker.Defunct <- client
		conn.Close()
	}()

	for {
		var cmd clientCommand
		if err := conn.ReadJSON(&cmd); err != nil {
			break
		}

		switch cmd.Action {
		case "subscribe":
			client.SubscribeTo(cmd.Topic)
		case "unsubscribe":
			client.UnsubscribeFrom(cmd.Topic)
		}
	}
}
