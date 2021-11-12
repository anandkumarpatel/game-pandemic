package main

import (
	"fmt"
)

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
