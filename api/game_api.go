package api

import (
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
	JoinCode string       `json:"joinCode"`
	GameID   string       `json:"gameID"`
	Player   *game.Player `json:"player"`
}

func handleNewGameGET(w http.ResponseWriter, r *http.Request) {
	// Start websocket server for this newGame session.
	newGame := game.NewGame()
	// Save to "DB"
	err := game.RegisterGame(newGame)
	if err != nil {
		log.WithError(err).Error("could not register game")
		http.Error(w, "could not create game", http.StatusInternalServerError)
		return
	}
	player := newGame.NewPlayer()
	resp := newGameResponse{
		JoinCode: newGame.JoinCode,
		GameID:   newGame.ID,
		Player:   player,
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
	log.WithField("joinCode", newGame.JoinCode).Debug("started new newGame")
}
