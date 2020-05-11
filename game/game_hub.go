package game

import log "github.com/sirupsen/logrus"

// GameHub tracks currently connected players and sends events out when necessary.
type GameHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	incomingMessages chan *IncomingMessage

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

func newHub(messageChannel chan *IncomingMessage) *GameHub {
	return &GameHub{
		incomingMessages: messageChannel,
		broadcasts:       make(chan []byte, 500),
		messages:         make(chan *GameMessage, 500),
		register:         make(chan *Client, 10),
		unregister:       make(chan *Client, 10),
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
			found := false
		PlayerLoop:
			for client := range h.clients {
				if client.player.ID == message.Target.ID {
					select {
					case client.send <- *message.Message:
						found = true
						break PlayerLoop
					default:
						log.Printf("Could not send a message to a player")
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			if !found {
				log.Error("could not find player to send message")
			}
		}
	}
}
