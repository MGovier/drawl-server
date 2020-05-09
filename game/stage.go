package game

type GameStage string

const (
	GAME_STARTING GameStage = "gameStarting"
	GAME_RUNNING GameStage = "gameRunning"
	GAME_ENDED GameStage = "gameEnded"
)
