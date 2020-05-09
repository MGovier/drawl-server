package game

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Game struct {
	ID       string         `json:"gameID"`
	Hub      *GameHub       `json:"-"`
	JoinCode string         `json:"joinCode"`
	Players  []*Player      `json:"players"`
	Stage    GameStage      `json:"gameStage"`
	Journeys []*WordJourney `json:"wordJourneys"`
}

// Start a new game up, and return the UUID and join code.
func NewGame() *Game {
	game := Game{}
	// Generate UUID
	ID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Entropy problems, oh my")
	}
	game.ID = ID.String()
	// Start websocket server
	hub := newHub()
	go hub.run()
	game.Hub = hub
	// Create a game JoinCode
	game.JoinCode = generateJoinCode()
	// Set stage as starting
	game.Stage = GAME_STARTING
	// Init arrays
	game.Players = make([]*Player, 0)
	game.Journeys = make([]*WordJourney, 0)
	// broadcast some stuff for the time being...
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				players, _ := json.Marshal(game.Players)
				game.Hub.broadcasts <- players
			}
		}
	}()
	return &game
}

func (g *Game) StartGame() {
	// Create starting words and distribute them.
}

func (g *Game) NewPlayer() *Player {
	name := fmt.Sprintf("Player %v", len(g.Players))
	playerID, err := uuid.NewRandom()
	if err != nil {
		log.WithError(err).Fatal("error creating UUID for NewPlayer")
	}
	newPlayer := &Player{
		ID:   playerID.String(),
		Name: name,
	}
	g.Players = append(g.Players, newPlayer)
	return newPlayer
}

// Generate a random string of A-Z chars with len 4
func generateJoinCode() string {
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
