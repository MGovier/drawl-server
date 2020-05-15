package game

import (
	"errors"
)

var activeGames []*Game

func init() {
	activeGames = make([]*Game, 0)
}

func RegisterGame(game *Game) {
	activeGames = append(activeGames, game)
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
	for _, game := range activeGames {
		if game.JoinCode == joinCode {
			return game, nil
		}
	}
	return nil, errors.New("game not found")
}
