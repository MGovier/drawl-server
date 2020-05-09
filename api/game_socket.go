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
	// Get ID from URL, that should be something like /ws/<uuid>.
	ID := strings.TrimPrefix(r.URL.Path, "/ws/")
	// Check it's a real UUID (and handle crazy formats, sure).
	_, err := uuid.Parse(ID)
	if err != nil {
		log.WithError(err).Error("UUID not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameInstance, err := db.FindGameByID(ID)
	if err != nil {
		log.WithError(err).Error("game not found in WebSocket connection request")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Create a client and attach to the game hub.
	game.ServeWs(gameInstance.Hub, w, r)
}
