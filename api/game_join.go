package api

import (
	"drawl-server/db"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleJoinGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addXOriginHeader(w, r, handleJoinGamePOST)
	case http.MethodOptions:
		returnXOriginHeader(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

type joinGameRequest struct {
	JoinCode string `json:"joinCode"`
}

type joinGameResponse struct {
	JoinCode string `json:"joinCode"`
	GameID   string `json:"gameID"`
}

func handleJoinGamePOST(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024) // 12 bytes = 4 UTF-8 chars plus a bit for the JSON?
	decoder := json.NewDecoder(r.Body)
	var joinRequest joinGameRequest
	err := decoder.Decode(&joinRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.WithError(err).Debug("invalid join game body contents")
		return
	}
	game, err := db.FindGameByJoinCode(joinRequest.JoinCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp := joinGameResponse{GameID: game.ID, JoinCode: game.JoinCode}
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err).Error("could not marshal JoinGame response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respJson)
	if err != nil {
		log.WithError(err).Error("could not write JoinGame response")
	}
}
