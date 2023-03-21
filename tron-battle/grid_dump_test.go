package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGridFromStamp(t *testing.T) {
	stamp := `A2C10.18A2C3A5C2.18A10C2B4.14A7C5B4.14A2C6.6B2.14A2C12B2.14A12C2B3.13A12C2B3.13A3.5C6.1B2.13A3.5C3.4B2.13B4.4C1.6B2.16B1C5.6B2.16B14.27B3.27B3.27B3.27B3.27B3.73`
	playersStamp := `[23:9][9:10][3:15]`
	dump := `stamp: A2C10.18A2C3A5C2.18A10C2B4.14A7C5B4.14A2C6.6B2.14A2C12B2.14A12C2B3.13A12C2B3.13A3.5C6.1B2.13A3.5C3.4B2.13B4.4C1.6B2.16B1C5.6B2.16B14.27B3.27B3.27B3.27B3.27B3.73
players: [23:9][9:10][3:15]
AACCCCCCCCCC..................
AACCCAAAAACC..................
AAAAAAAAAACCBBBB..............
AAAAAAACCCCCBBBB..............
AACCCCCC......BB..............
AACCCCCCCCCCCCBB..............
AAAAAAAAAAAACCBBB.............
AAAAAAAAAAAACCBBB.............
AAA.....CCCCCC.BB.............
AAA.....CCC....BB.............
BBBB....C......BB.............
...BCCCCC......BB.............
...BBBBBBBBBBBBBB.............
..............BBB.............
..............BBB.............
..............BBB.............
..............BBB.............
..............BBB.............
..............................
..............................`
	grid := NewGridFromStamp(stamp, playersStamp)
	dump2 := fmt.Sprintf("%v", grid)
	assert.Equal(t, stamp, grid.GetStamp())
	assert.Equal(t, playersStamp, grid.GetPlayersStamp())
	assert.Equal(t, dump, dump2)
}
