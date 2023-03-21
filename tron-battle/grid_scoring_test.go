package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayerScore(t *testing.T) {
	grid := NewGrid(3)
	grid.Players[0] = &GridPlayer{0, 0}
	grid.Players[1] = &GridPlayer{GridWidth - 1, 0}
	grid.Players[2] = &GridPlayer{0, GridHeight - 1}
	assert.Equal(t, 600.0, grid.GetPlayerScore(0))
	assert.Equal(t, []float64{600.0, 600.0, 600.0}, grid.GetPlayerScores())
}

func BenchmarkGetPlayerScore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		grid := NewGrid(3)
		grid.Players[0] = &GridPlayer{0, 0}
		grid.Players[1] = &GridPlayer{GridWidth - 1, 0}
		grid.Players[2] = &GridPlayer{0, GridHeight - 1}
		grid.GetPlayerScore(0)
		grid.GetPlayerScore(1)
		grid.GetPlayerScore(2)
	}
}

func BenchmarkGetPlayerScores(b *testing.B) {
	for n := 0; n < b.N; n++ {
		grid := NewGrid(3)
		grid.Players[0] = &GridPlayer{0, 0}
		grid.Players[1] = &GridPlayer{GridWidth - 1, 0}
		grid.Players[2] = &GridPlayer{0, GridHeight - 1}	
		grid.GetPlayerScores()
	}
}
