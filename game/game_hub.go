package game

import log "github.com/sirupsen/logrus"

// GameHub tracks currently connected players and sends events out when necessary.
type GameHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	incomingMessages chan []byte

	// Messages to send to all clients
	broadcasts chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *GameHub {
	return &GameHub{
		incomingMessages: make(chan []byte),
		broadcasts:       make(chan []byte),
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
			// TODO: handle game messages
			log.Debug(message)
		case message := <-h.broadcasts:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}