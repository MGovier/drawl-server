package api

import (
	"drawl-server/db"
	"encoding/json"
	"net/http"
)

func HandleGetGameReview(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		addXOriginHeader(w, r, handleGetGameReviewGET)
	case http.MethodOptions:
		returnXOriginHeader(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func handleGetGameReviewGET(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		http.Error(w, "missing game_id query parameter", http.StatusBadRequest)
		return
	}
	game, err := db.FindGameByID(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Serialize the whole game, oh boy
	jsnData, err := json.Marshal(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsnData)
}
