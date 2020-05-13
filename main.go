package main

import (
	"drawl-server/api"
	"flag"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	log.SetLevel(log.DebugLevel)

	http.HandleFunc("/game", api.HandleNewGame)
	http.HandleFunc("/join", api.HandleJoinGame)
	http.HandleFunc("/review", api.HandleGetGameResults)
	http.HandleFunc("/", api.HandleWS)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
