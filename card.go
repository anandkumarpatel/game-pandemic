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
	Bonus     CardBonus
	Type      CardType
	VirusType VirusType
	Name      string
}

func (s Card) String() string {
	return fmt.Sprintf("(%s)", s.Name)
}
