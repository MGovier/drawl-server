package game

import "errors"

type Player struct {
	ID   string `json:"playerID"`
	Name string `json:"playerName"`
}

func (p *Player) SetName(name string) error {
	if len(name) > 15 {
		return errors.New("name too long")
	}
	if len(name) < 1 {
		return errors.New("name empty")
	}
	p.Name = name
	return nil
}
