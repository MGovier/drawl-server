package game

import log "github.com/sirupsen/logrus"

// GameHub tracks currently connected players and sends events out when necessary.
type GameHub struct {
	// Registered clients.
	clients map[string]*Client

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
		clients:          make(map[string]*Client),
	}
}

func (h *GameHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.player.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.player.ID]; ok {
				delete(h.clients, client.player.ID)
				close(client.send)
			}
		case message := <-h.broadcasts:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.player.ID)
				}
			}
		case message := <-h.messages:
			client, found := h.clients[message.Target.ID]
			if !found {
				log.Error("could not find player to send message")
				return
			}
			select {
			case client.send <- *message.Message:
			default:
				log.Printf("Could not send a message to a player")
				// TODO: Put the message back on the queue? Wait for reconnection?
				close(client.send)
				delete(h.clients, client.player.ID)
			}

		}
	}
}
