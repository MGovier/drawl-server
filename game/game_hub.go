package game

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// GameHub tracks currently connected players and sends events out when necessary.
type GameHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	incomingMessages chan IncomingMessage

	// Messages to send to all clients.
	broadcasts chan []byte

	// Messages to send to specific players.
	messages chan *GameMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type GameMessage struct {
	Target  *Player
	Message *[]byte
}

type IncomingMessage struct {
	Player  *Player
	Message []byte
}

type IncomingMessageContents struct {
	Type     string      `json:"type"`
	Contents interface{} `json:"data"`
}

func newHub() *GameHub {
	return &GameHub{
		incomingMessages: make(chan IncomingMessage),
		broadcasts:       make(chan []byte),
		messages:         make(chan *GameMessage),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
	}
}

func (h *GameHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.incomingMessages:
			h.handleMessage(message)
		case message := <-h.broadcasts:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.messages:
			for client := range h.clients {
				if client.player.ID == message.Target.ID {
					select {
					case client.send <- *message.Message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func (g *GameHub) handleMessage(message IncomingMessage) {
	// Try to deserialize message from JSON as above type.
	var msg IncomingMessageContents
	err := json.Unmarshal(message.Message, &msg)
	if err != nil {
		log.WithError(err).Error("Could not unmarshal client WebSocket message")
		return
	}
	if msg.Type == "name" {
		// Data should be just a string with the new name.
		newName, ok := msg.Contents.(string)
		if !ok {
			log.WithError(err).Error("Could not cast name message contents to string")
			return
		}
		err = message.Player.SetName(newName)
		if err != nil {
			log.WithError(err).Error("Error setting a player's name")
			return
		}
	}
}
