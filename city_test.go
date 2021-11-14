package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCity_HasVirus(t *testing.T) {
	tests := []struct {
		name string
		s    City
		want bool
	}{
		{"no virus", City{}, false},
		{"no virus", City{VirusCounts: VirusCounts{Black: 0}}, false},
		{"one virus", City{VirusCounts: VirusCounts{Black: 1}}, true},
		{"two virus", City{VirusCounts: VirusCounts{Black: 1, Blue: 2}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.s.HasVirus())
		})
	}
}
