package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGridGoRight(t *testing.T) {
	grid := NewGridFromStamp("A2C3.25A2C3.25A3.252B1.29B1.29B1.29B1.29B1.29B1.29B1.104")
	// TODO: make this loadable via a stamp too
	grid.Players = []*GridPlayer{{2, 2}, {4, 0}, {15, 10}}
	grid2, cell := grid.GoRight(0)

	assert.NotNil(t, grid2)
	assert.NotNil(t, cell)

	assert.Equal(t, cell.X, 3)
	assert.Equal(t, cell.Y, 2)
	assert.Equal(t, cell.Player, PlayerName_A)
	assert.Equal(t, cell.Type, CellType_Full)
	assert.Equal(t, grid2.Players[0].X, 3)
	assert.Equal(t, grid2.Players[0].Y, 2)

	assert.Equal(t, grid.Players[0].X, 2)
	assert.Equal(t, grid.Players[0].Y, 2)

	assert.NotSame(t, grid, grid2)
	assert.NotSame(t, grid.Players, grid2.Players)
	assert.NotSame(t, grid.Players[0], grid2.Players[0])
	assert.NotSame(t, grid.Cells, grid2.Cells)
	assert.NotSame(t, grid.Cells[cell.Index], cell)
	assert.NotSame(t, grid.Cells[cell.Index], grid2.Cells[cell.Index])
	assert.NotEqual(t, grid.GetStamp(), grid2.GetStamp())
	assert.Equal(
		t,
		"A2C3.25A2C3.25A4.251B1.29B1.29B1.29B1.29B1.29B1.29B1.104",
		grid2.GetStamp(),
	)
}
