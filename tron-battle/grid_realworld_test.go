package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// FIXME: still a bug
func TestOneGrid(t *testing.T) {
	grid := NewGridFromStamp("A29.1A29.25C2.2A1.23C4A3.23C1.2C1A4.19C7.23C1.10C20.10C2.19C2.8C1.20C1.8C20.1C1.27C3.247", "[29:4][-1:-1][22:4]")
	fmt.Printf("%v", grid)
	now := time.Now()
	_, score := grid.MinMax(0, 3)
	fmt.Printf("time: %v", time.Since(now))
	assert.Greater(t, 0, score)
}