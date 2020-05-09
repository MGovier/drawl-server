package api

import (
	"drawl-server/db"
	"drawl-server/game"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleNewGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		addXOriginHeader(w, r, handleNewGameGET)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

type newGameResponse struct {
	JoinCode string `json:"joinCode"`
	GameID   string `json:"gameID"`
}

func handleNewGameGET(w http.ResponseWriter, r *http.Request) {
	// Start websocket server for this game session.
	game := game.NewGame()
	// Save to "DB"
	db.RegisterGame(game)
	// Generate a game UUID and external connection code
	resp := newGameResponse{
		JoinCode: game.JoinCode,
		GameID:   game.ID,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err).Error("could not marshal NewGame response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respJson)
	if err != nil {
		log.WithError(err).Error("could not write NewGame response")
	}
	log.WithField("joinCode", game.JoinCode).Debug("started new game")
}
