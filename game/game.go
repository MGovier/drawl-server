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
	ID         string                `json:"gameID"`
	Hub        *GameHub              `json:"-"`
	JoinCode   string                `json:"joinCode"`
	Players    []*Player             `json:"players"`
	Stage      GameStage             `json:"gameStage"`
	Journeys   []*WordJourney        `json:"wordJourneys"`
	Round      int                   `json:"round"`
	Limit      int                   `json:"limit"`
	GameEvents chan *IncomingMessage `json:"-"`
}

// Start a new game up, and return the UUID and join code.
func NewGame() *Game {
	game := Game{}
	game.GameEvents = make(chan *IncomingMessage)
	// Generate UUID
	ID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Entropy problems, oh my")
	}
	game.ID = ID.String()
	// Start websocket server
	hub := newHub(game.GameEvents)
	go hub.run()
	game.Hub = hub
	// Create a game JoinCode
	game.JoinCode = generateJoinCode()
	// Set stage as starting
	game.Stage = GAME_STARTING
	// Init arrays
	game.Players = make([]*Player, 0)
	game.Journeys = make([]*WordJourney, 0)
	go game.run()
	// Start broadcast of player names until game begins
	game.broadcastPlayers()
	return &game
}

func (g *Game) broadcastPlayers() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if g.Stage != GAME_STARTING {
					ticker.Stop()
					return
				}
				g.sendPlayers()
			}
		}
	}()
}

func (g *Game) StartGame() {
	if g.Stage != GAME_STARTING {
		log.Debug("Attempted to start an already running game")
		return
	}
	// Players can not change at this point. Stops broadcasts.
	g.Stage = GAME_RUNNING
	// One final broadcast to make sure we didn't miss any in the last second...
	g.sendPlayers()
	// Number of rounds is just number of players
	g.Limit = len(g.Players)
	g.Round = 0
	// Create starting words and distribute them.
	g.startJourneys()
	g.sendNextRoundToPlayers()
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

func (g *Game) run() {
	for {
		select {
		case incomingMessage := <-g.GameEvents:
			g.HandleMessage(incomingMessage)
		}
	}
}

func (g *Game) HandleMessage(message *IncomingMessage) {
	// Try to deserialize message from JSON as above type.
	var msg IncomingMessageContents
	err := json.Unmarshal(message.Message, &msg)
	if err != nil {
		log.WithError(err).Error("Could not unmarshal client WebSocket message")
		return
	}
	if msg.Type == "name" {
		// Data should be just a string with the new name.
		newName, ok := msg.Contents.(string)
		if !ok {
			log.WithError(err).Error("Could not cast name message contents to string")
			return
		}
		err = message.Player.SetName(newName)
		if err != nil {
			log.WithError(err).Error("Error setting a player's name")
			return
		}
	}
	if msg.Type == "start" {
		// Check correct player started the game for *essential security*.
		if message.Player != g.Players[0] {
			log.Error("Incorrect player tried to start the game")
			return
		}
		g.StartGame()
	}
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

func (g *Game) startJourneys() {
	for _, player := range g.Players {
		word := generateWord()
		g.Journeys = append(g.Journeys, &WordJourney{
			StartingWord:   word,
			StartingPlayer: player,
			Plays:          make([]*GamePlay, 0),
		})
	}
}
