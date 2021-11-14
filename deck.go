package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Deck struct {
	Cards []*Card
}

func (s Deck) String() string {
	return fmt.Sprintf("size: %d %s", len(s.Cards), s.Cards)
}

func (s Deck) CardNames() []string {
	out := []string{}
	for _, card := range s.Cards {
		out = append(out, card.Name)
	}
	return out
}

func (s Deck) Count() int {
	return len(s.Cards)
}

func (s *Deck) AddCard(card *Card) {
	s.Cards = append(s.Cards, card)
}

func (s *Deck) AddDeck(deck *Deck) {
	s.Cards = append(s.Cards, deck.Cards...)
}

func (s *Deck) Draw() *Card {
	out := s.Cards[0]
	s.RemoveCard(out.Name)
	return out
}

func (s *Deck) RemoveCard(cardName string) *Card {
	i := s.FindCardIndexWithName(cardName)
	out := s.Cards[i]
	s.Cards = append(s.Cards[:i], s.Cards[i+1:]...)
	return out
}

func (s *Deck) Clear() {
	s.Cards = []*Card{}
}

func (s *Deck) FindCardIndexWithName(name string) int {
	for i, v := range s.Cards {
		if v.Name == name {
			return i
		}
	}

	return -1
}

func (s *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s.Cards), func(i, j int) { s.Cards[i], s.Cards[j] = s.Cards[j], s.Cards[i] })
}
