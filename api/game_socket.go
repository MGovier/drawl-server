package api

import (
	"drawl-server/db"
	"drawl-server/game"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Enable connecting to the game's WebSocket hub.
func HandleWS(w http.ResponseWriter, r *http.Request) {
	// Get Game ID and Player ID from URL, that should be something like /ws/<gameID>/<playerID>.
	IDs := strings.Split(r.URL.Path, "/")
	// Check we got real UUIDs (and handle crazy formats, sure).
	if len(IDs) != 4 {
		log.Error("Could not find required fields in WebSocket connection request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err := uuid.Parse(IDs[1])
	if err != nil {
		log.WithError(err).Error("Game UUID not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = uuid.Parse(IDs[2])
	if err != nil {
		log.WithError(err).Error("Player UUID not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameInstance, err := db.FindGameByID(IDs[1])
	if err != nil {
		log.WithError(err).Error("game not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Check Player is in this game...
	var player *game.Player = nil
	for _, playr := range gameInstance.Players {
		if playr.ID == IDs[2] {
			player = playr
			break
		}
	}
	if player == nil {
		log.WithError(err).Error("player not found in this game during WebSocket connection request")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Create a client and attach to the game hub.
	game.ServeWs(gameInstance.Hub, player, w, r)
}
