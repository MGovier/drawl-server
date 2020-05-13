package game

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// Updates we can send to players
// These include player updates, and round info

type gameUpdate struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

func (g *Game) sendPlayers() {
	players, _ := json.Marshal(g.Players)
	gameUpdate, _ := json.Marshal(gameUpdate{
		Type: "players",
		Data: players,
	})
	g.Hub.broadcasts <- gameUpdate

}

// Give each player their appropriate words to draw
func (g *Game) sendNextRoundToPlayers() {
	if g.Round == g.Limit {
		// Game Over!
		// TODO: Create end of game summaries and point system...
		update := gameUpdate{
			Type: "review",
			Data: nil,
		}
		messageBytes, err := json.Marshal(update)
		if err != nil {
			log.WithError(err).Error("problem marshalling game update to JSON")
			return
		}
		select {
		case g.Hub.broadcasts <- messageBytes:
		default:
			log.Error("Could not dispatch message")
		}
		return
	}
	for _, journey := range g.Journeys {
		// Chose stage in journey according to round.
		var update gameUpdate
		if g.Round%2 == 0 {
			update = gameUpdate{
				Type: "word",
				Data: []byte(journey.Plays[g.Round].GetPlay()),
			}
		} else {
			update = gameUpdate{
				Type: "drawing",
				Data: []byte(journey.Plays[g.Round].GetPlay()),
			}
		}
		messageBytes, err := json.Marshal(update)
		if err != nil {
			log.WithError(err).Error("problem marshalling game update to JSON")
			return
		}
		msg := &GameMessage{
			Target:  journey.Order[g.Round],
			Message: &messageBytes,
		}
		select {
		case g.Hub.messages <- msg:
		default:
			log.Error("Could not dispatch message")
		}
	}
}

// TODO: Give a disconnected client everything they need to catch up...
func (g *Game) reconnectPlayer() {

}
