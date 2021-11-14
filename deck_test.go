package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeck_CardNames(t *testing.T) {
	tests := []struct {
		name  string
		Cards []*Card
		want  []string
	}{
		{"no cards", []*Card{}, []string{}},
		{"one cards", []*Card{{Name: "one"}}, []string{"one"}},
		{"two cards", []*Card{{Name: "one"}, {Name: "two"}}, []string{"one", "two"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Deck{
				Cards: tt.Cards,
			}
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, s.CardNames())
			})
		})
	}
}
