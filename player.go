package main

import (
	"fmt"
)

type Players []*Player

type Player struct {
	Name     string
	Location string
	Hand     Deck
}

func (s Player) String() string {
	return fmt.Sprintf("|N:%s L:%s H:%s)|", s.Name, s.Location, s.Hand)
}

func (s *Player) Move(cityName string) {
	s.Location = cityName
}

func (s Players) FindByName(name string) *Player {
	for _, player := range s {
		if player.Name == name {
			return player
		}
	}

	panic(fmt.Errorf("Players.FindByName: cant find player %s", name))
}
