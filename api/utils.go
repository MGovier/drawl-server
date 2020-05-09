package api

import (
	"drawl-server/config"
	"net/http"
)

func addXOriginHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Access-Control-Allow-Origin", config.AllowedCORSOrigin)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	next(w, r)
}

func returnXOriginHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", config.AllowedCORSOrigin)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Write([]byte("200 OK"))
}
