package game

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// Updates we can send to players
// These include player updates, and round info

type gameUpdate struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (g *Game) sendPlayers() {
	players, _ := json.Marshal(g.Players)
	g.Hub.broadcasts <- players
}

// Give each player their appropriate words to draw
func (g *Game) sendNextRoundToPlayers() {
	for _, journey := range g.Journeys {
		// Chose stage in journey according to round.
		if g.Round == 0 {
			update := gameUpdate{
				Type: "newWord",
				Data: journey.StartingWord,
			}
			messageBytes, err := json.Marshal(update)
			if err != nil {
				log.WithError(err).Error("problem marshalling game update to JSON")
			}
			msg := &GameMessage{
				Target:  journey.StartingPlayer,
				Message: &messageBytes,
			}
			g.Hub.messages <- msg
		}
	}
}

// TODO: Give a disconnected client everything they need to catch up...
func (g *Game) reconnectPlayer() {

}
