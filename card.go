package main

import "fmt"

type CardBonus int

const (
	None CardBonus = iota
)

type CardType int

const (
	CityCardType CardType = iota
	PandemicCardType
	EventCardType
)

type Card struct {
	Name  string
	Bonus CardBonus
	Type  CardType
	City  *City
}

func (s Card) String() string {
	return fmt.Sprintf("(%s)", s.Name)
}
