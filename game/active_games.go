package game

import (
	"errors"
	"math/rand"
)

var activeGames []*Game
var joinCodes map[string]*Game

func init() {
	activeGames = make([]*Game, 0)
	joinCodes = make(map[string]*Game)
}

// Register the game and create a join code for it
func RegisterGame(game *Game) error {
	activeGames = append(activeGames, game)
	joinCode, err := generateJoinCode()
	if err != nil {
		return err
	}
	game.JoinCode = joinCode
	return nil
}

func UnregisterGame(gameID string) {
	if len(activeGames) < 2 {
		activeGames = make([]*Game, 0)
		return
	}
	for i, game := range activeGames {
		if game.ID == gameID {
			activeGames[i] = activeGames[len(activeGames)-1]
			activeGames = activeGames[:len(activeGames)-1]
		}
	}
}

func FindGameByID(gameID string) (*Game, error) {
	for _, game := range activeGames {
		if game.ID == gameID {
			return game, nil
		}
	}
	return nil, errors.New("game not found")
}

func FindGameByJoinCode(joinCode string) (*Game, error) {
	game, found := joinCodes[joinCode]
	if !found {
		return nil, errors.New("game not found, or no longer joinable")
	}
	return game, nil
}

func RemoveGameJoinCode(game *Game) {
	delete(joinCodes, game.JoinCode)
}

// Generate a random string of A-Z chars with len 4
func generateJoinCode() (string, error) {
	for i := 0; i < 100; i++ {
		code := randomCode()
		if _, found := joinCodes[code]; found {
			if !found {
				return code, nil
			}
		}
	}
	return "", errors.New("could not find free join code")
}

func randomCode() string {
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
