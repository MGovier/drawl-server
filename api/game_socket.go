package api

import (
	"drawl-server/game"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Enable connecting to the game's WebSocket hub.
func HandleWS(w http.ResponseWriter, r *http.Request) {
	// Get Game ID and Player ID from URL, that should be something like /ws/<gameID>/<playerID>.
	IDs := r.URL.Query()
	gameID := IDs.Get("game_id")
	playerID := IDs.Get("player_id")
	if playerID == "" || gameID == "" {
		http.Error(w, "missing player or game ID", http.StatusBadRequest)
		return
	}
	_, err := uuid.Parse(gameID)
	if err != nil {
		log.WithError(err).Error("game UUID not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = uuid.Parse(playerID)
	if err != nil {
		log.WithError(err).Error("player UUID not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameInstance, err := game.FindGameByID(gameID)
	if err != nil {
		log.WithError(err).Error("game not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Check Player is in this game...
	var player *game.Player = nil
	for _, playr := range gameInstance.Players {
		if playr.ID == playerID {
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
	IP := r.Header.Get("X-Forwarded-For")
	if IP == "" {
		IP = r.RemoteAddr
	}
	log.WithFields(log.Fields{
		"IP":       IP,
		"joinCode": gameInstance.JoinCode,
		"playerID": player.ID,
	}).Debug("player connected")
}
