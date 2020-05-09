package game

// The wild ride the drawings and guesses hopefully go through!
type WordJourney struct {
	StartingWord string `json:"startingWord"`
	StartingPlayer *Player `json:"startingPlayer"`
	WordUUID string `json:"wordUUID"`
	Plays []*GamePlay `json:"gamePlays"`
}

// A game play can be either a guess or a drawing, this is used to track all events in a game and... eventually...
// allow a disconnected client to reconnect and catch back up.

type GamePlay interface {
	GetPlay() string
	GetPlayer() *Player
}

type Guess struct {
	Guess string `json:"guess"`
	Player *Player `json:"player"`
}

func (g *Guess) GetPlay() string {
	return g.Guess
}

func (g *Guess) GetPlayer() *Player {
	return g.Player
}

type Drawing struct {
	Drawing string `json:"drawing"`
	Player *Player `json:"player"`
}

func (d *Drawing) GetPlay() string {
	return d.Drawing
}

func (d *Drawing) GetPlayer() *Player {
	return d.Player
}
