package db

import (
	"drawl-server/game"
	"errors"
)

var activeGames []*game.Game

func init() {
	activeGames = make([]*game.Game, 0)
}

func RegisterGame(game *game.Game) {
	activeGames = append(activeGames, game)
}

func FindGameByID(gameID string) (*game.Game, error) {
	for _, game := range activeGames {
		if game.ID == gameID {
			return game, nil
		}
	}
	return nil, errors.New("game not found")
}

func FindGameByJoinCode(joinCode string) (*game.Game, error) {
	for _, game := range activeGames {
		if game.JoinCode == joinCode {
			return game, nil
		}
	}
	return nil, errors.New("game not found")
}
