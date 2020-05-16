package game

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// GameHub tracks currently connected players and sends events out when necessary.
type GameHub struct {
	// Registered clients.
	clients map[string]*Client
	// Inbound messages from the clients.
	incomingMessages chan *IncomingMessage
	// Reconnecting players
	reconnections chan *Player
	// Messages to send to all clients.
	broadcasts chan *GameMessage
	// Messages to send to specific players.
	messages chan *GameMessage
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	// Clients that have been connected before
	history []string
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

func newHub(messageChannel chan *IncomingMessage, reconnectionChannel chan *Player) *GameHub {
	return &GameHub{
		incomingMessages: messageChannel,
		reconnections:    reconnectionChannel,
		broadcasts:       make(chan *GameMessage, 32),
		messages:         make(chan *GameMessage, 64),
		register:         make(chan *Client, 10),
		unregister:       make(chan *Client, 10),
		clients:          make(map[string]*Client),
		history:          make([]string, 0),
	}
}

func (h *GameHub) run() {
	timeout := time.After(6 * time.Hour)
	running := true
	for running {
		select {
		case client := <-h.register:
			h.clients[client.player.ID] = client
			var previousClient = false
			for _, oldClientID := range h.history {
				if oldClientID == client.player.ID {
					// They must be reconnecting, wb!
					log.WithField("player", client.player).Debug("player reconnected")
					// Lets give them their last update again in case they missed it.
					previousClient = true
					h.reconnections <- client.player
				}
			}
			if !previousClient {
				h.history = append(h.history, client.player.ID)
			}
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
					log.Debug("could not send broadcast")
					close(client.send)
					delete(h.clients, client.player.ID)
				}
			}
		case message := <-h.messages:
			if message == nil || message.Target == nil {
				// I can't send this!
				continue
			}
			client, found := h.clients[message.Target.ID]
			if !found {
				log.Error("could not find player to send message")
				h.putMessageBack(message)
				continue
			}
			select {
			case client.send <- message:
			default:
				log.Printf("Could not send a message to a player")
				h.putMessageBack(message)
				close(client.send)
				delete(h.clients, client.player.ID)
			}
		case <-timeout:
			for _, client := range h.clients {
				close(client.send)
			}
			running = false
		}
	}
}

// Put messages back on the send queue after an interval
func (h *GameHub) putMessageBack(message *GameMessage) {
	go func() {
		select {
		case <-time.After(3 * time.Second):
			h.messages <- message
			log.Debug("placed messaged back on the queue")
		}
	}()
}
