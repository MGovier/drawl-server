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
	ID              string                `json:"gameID"`
	Hub             *GameHub              `json:"-"`
	JoinCode        string                `json:"joinCode"`
	Players         []*Player             `json:"players"`
	PlayerMap       map[string]*Player    `json:"-"`
	PlayersFinished []*Player             `json:"-"`
	Stage           GameStage             `json:"gameStage"`
	Journeys        []*WordJourney        `json:"wordJourneys"`
	Round           int                   `json:"round"`
	Limit           int                   `json:"limit"`
	GameEvents      chan *IncomingMessage `json:"-"`
	// Send events when game events occur to check if the round is complete!
	GameProgressChecker chan struct{} `json:"-"`
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
	game.PlayerMap = make(map[string]*Player)
	game.Players = make([]*Player, 0)
	game.Journeys = make([]*WordJourney, 0)
	game.GameProgressChecker = make(chan struct{})
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
	if g.Stage == GAME_RUNNING {
		log.Debug("Attempted to start an already running game")
		return
	}
	// Players can not change at this point. Stops broadcasts.
	g.Stage = GAME_RUNNING
	// One final broadcast to make sure we didn't miss any in the last second...
	g.sendPlayers()
	// Reset finished players
	g.PlayersFinished = make([]*Player, 0)
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
	g.PlayerMap[newPlayer.ID] = newPlayer
	return newPlayer
}

func (g *Game) run() {
	for {
		select {
		case incomingMessage := <-g.GameEvents:
			go g.HandleMessage(incomingMessage)
		case <-g.GameProgressChecker:
			g.checkAndAdvanceRound()
		}
	}
}

func (g *Game) checkAndAdvanceRound() {
	// The review process has been finished (that is handled by incrementing the round then giving to sendNextRound)
	if g.Round == g.Limit {
		if len(g.PlayersFinished) == len(g.Players) {
			g.sendResults()
			g.Stage = GAME_ENDED
		}
		return
	}
	waitingFor := make([]*Player, 0)
	for _, journey := range g.Journeys {
		if len(journey.Plays)-1 <= g.Round {
			waitingFor = append(waitingFor, journey.Order[g.Round])
		}
	}
	if len(waitingFor) == 0 {
		g.Round++
		g.sendNextRoundToPlayers()
	}
	log.Debugf("Waiting for: %v", waitingFor)
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
	if msg.Type == "drawing" {
		// Find what journey this belongs to based on the round and player order.
		// TODO: Make better!
		for _, journey := range g.Journeys {
			if journey.Order[g.Round].ID == message.Player.ID {
				// Right journey!
				drawData, ok := msg.Contents.(string)
				if !ok {
					log.Error("could not cast drawing data to string")
					return
				}
				journey.Plays = append(journey.Plays, &Drawing{
					Drawing: drawData,
					Player:  message.Player,
				})
				g.GameProgressChecker <- struct{}{}
			}
		}
	}
	if msg.Type == "guess" {
		// Find what journey this belongs to based on the round and player order.
		// TODO: Make better!
		for _, journey := range g.Journeys {
			if journey.Order[g.Round].ID == message.Player.ID {
				// Right journey!
				guess, ok := msg.Contents.(string)
				if !ok {
					log.Error("could not cast guess data to string")
					return
				}
				journey.Plays = append(journey.Plays, &Word{
					Word:   guess,
					Player: message.Player,
				})
				g.GameProgressChecker <- struct{}{}
			}
		}
	}
	if msg.Type == "award" {
		// TODO: Stop people awarding the same person multiple times per game
		target, ok := msg.Contents.(string)
		if !ok {
			log.WithError(err).Error("Could not cast award message contents to string")
			return
		}
		if target == message.Player.ID {
			// Nice tryyyyyy
			return
		}
		g.PlayerMap[target].Points++
	}
	if msg.Type == "done" {
		for _, player := range g.PlayersFinished {
			if player.ID == message.Player.ID {
				// Already registered as done
				return
			}
		}
		g.PlayersFinished = append(g.PlayersFinished, message.Player)
		g.GameProgressChecker <- struct{}{}
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
	// Reset things in case it's a new round.
	g.Journeys = make([]*WordJourney, 0)
	for offset, _ := range g.Players {
		word := generateWord()
		order := g.calculatePlayOrder(offset)
		startingPlay := &Word{
			Word:   word,
			Player: nil,
		}
		newJourney := &WordJourney{
			Order: order,
			Plays: make([]GamePlay, 0),
		}
		newJourney.Plays = append(newJourney.Plays, startingPlay)
		g.Journeys = append(g.Journeys, newJourney)
	}
}

func (g *Game) registerDrawing() {

}

func (g *Game) calculatePlayOrder(offset int) []*Player {
	order := make([]*Player, 0)
	for i, _ := range g.Players {
		order = append(order, g.Players[(i+offset)%len(g.Players)])
	}
	return order
}
