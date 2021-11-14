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
	Name      string
	Bonus     CardBonus
	Type      CardType
	VirusType VirusType
	City      *City
}

func (s Card) String() string {
	return fmt.Sprintf("(%s)", s.Name)
}
