package floodfill

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessEmpty(t *testing.T) {
	expected := `00 00 00
00 00 00`
	ff := NewFloodFill(3, 2)
	// ff.MarkPointUnreachable(&Point{0, 1})
	// ff.MarkPointUnreachable(&Point{0, 2})
	// ff.MarkPointUnreachable(&Point{0, 3})
	ff.Process()
	assert.Equal(t, expected, ff.String())
}

func TestProcess1(t *testing.T) {
	expected := `00 00 00
++ ++ ++
01 01 01`
	ff := NewFloodFill(3, 3)
	ff.MarkPointUnreachable(&Point{0, 1})
	ff.MarkPointUnreachable(&Point{1, 1})
	ff.MarkPointUnreachable(&Point{2, 1})
	ff.Process()
	assert.Equal(t, expected, ff.String())
}

func TestProcess2(t *testing.T) {
	expected := `00 00 00
++ ++ 00
00 00 00`
	ff := NewFloodFill(3, 3)
	ff.MarkPointUnreachable(&Point{0, 1})
	ff.MarkPointUnreachable(&Point{1, 1})
	ff.Process()
	assert.Equal(t, expected, ff.String())
}

func TestProcess3(t *testing.T) {
	expected := `00 00 00
++ ++ 00
01 01 ++`
	ff := NewFloodFill(3, 3)
	ff.MarkPointUnreachable(&Point{0, 1})
	ff.MarkPointUnreachable(&Point{1, 1})
	ff.MarkPointUnreachable(&Point{2, 2})
	ff.Process()
	assert.Equal(t, expected, ff.String())
}

func TestProcess4(t *testing.T) {
	expected := `00 00 00
++ ++ 00
01 01 ++
++ ++ 02
02 02 02
02 02 02`
	expected2 := `00 00 00
++ ++ 00
01 01 ++
01 ++ 01
01 01 01
01 01 01`
	expected3 := `00 00 00
++ ++ 00
00 00 00
00 ++ 00
00 00 00
00 00 00`
	ff := NewFloodFill(3, 6)
	ff.MarkPointUnreachable(&Point{0, 1})
	ff.MarkPointUnreachable(&Point{1, 1})
	ff.MarkPointUnreachable(&Point{2, 2})
	ff.MarkPointUnreachable(&Point{0, 3})
	ff.MarkPointUnreachable(&Point{1, 3})
	ff.Process()
	assert.Equal(t, expected, ff.String())
	ff.MarkPointEmpty(&Point{0, 3}).Process()
	assert.Equal(t, expected2, ff.String())
	ff.MarkPointEmpty(&Point{2, 2}).Process()
	assert.Equal(t, expected3, ff.String())
}

func Benchmark1(b *testing.B) {
	ff := NewFloodFill(3, 6)
	ff.MarkPointUnreachable(&Point{0, 1})
	ff.MarkPointUnreachable(&Point{1, 1})
	ff.MarkPointUnreachable(&Point{2, 2})
	ff.MarkPointUnreachable(&Point{0, 3})
	ff.MarkPointUnreachable(&Point{1, 3})
	for n := 0; n < b.N; n++ {
		ff.
			MarkPointUnreachable(&Point{0, 3}).
			MarkPointUnreachable(&Point{2, 2}).
			Process()
		ff.
			MarkPointEmpty(&Point{0, 3}).
			MarkPointUnreachable(&Point{2, 2}).
			Process()
		ff.
			MarkPointEmpty(&Point{0, 3}).
			MarkPointEmpty(&Point{2, 2}).
			Process()
	}
}

func getBigCorridor() FloodFill {
	ff := NewFloodFill(20, 20)
	for x := 0; x < 19; x += 1 {
		ff.MarkPointUnreachable(&Point{x, 1})
	}
	for y := 1; y < 20; y += 1 {
		ff.MarkPointUnreachable(&Point{18, y})
	}
	return ff
}

func TestBigCorridor(t *testing.T) {
	ff := getBigCorridor()
	ff.Process()
	assert.Equal(t, strings.TrimSpace(`
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
`), ff.String())
	ff.MarkPointEmpty(&Point{18, 19}).Process()
	assert.Equal(t, strings.TrimSpace(`
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
`), ff.String())
	ff.MarkPointUnreachable(&Point{18, 19}).Process()
	assert.Equal(t, strings.TrimSpace(`
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 00
`), ff.String())
}

func Benchmark2(b *testing.B) {
	ff := getBigCorridor()
	for n := 0; n < b.N; n++ {
		ff.MarkPointEmpty(&Point{18, 19}).Process()
		ff.MarkPointUnreachable(&Point{18, 19}).Process()
	}
}

func setSquare(ff FloodFill, x, y, w, h int) FloodFill {
	for xx := x; xx < w+x; xx += 1 {
		ff.MarkPointUnreachable(&Point{xx, y})
	}
	for xx := x; xx < w+x; xx += 1 {
		ff.MarkPointUnreachable(&Point{xx, y + h - 1})
	}
	for yy := y; yy < h+y; yy += 1 {
		ff.MarkPointUnreachable(&Point{x, yy})
	}
	for yy := y; yy < h+y; yy += 1 {
		ff.MarkPointUnreachable(&Point{x + w - 1, yy})
	}
	return ff
}

func getLotOfZones() FloodFill {
	ff := NewFloodFill(20, 20)
	// zone 1
	ff = setSquare(ff, 3, 1, 7, 5)
	// zone 3
	ff = setSquare(ff, 12, 3, 7, 5)
	// zone 2
	ff = setSquare(ff, 1, 3, 3, 14)
	// zone 5
	ff = setSquare(ff, 3, 11, 12, 3)
	// zone 7
	ff = setSquare(ff, 13, 13, 7, 7)
	return ff
}

func TestLotOfZones(t *testing.T) {
	ff := getLotOfZones()
	ff.Process()
	assert.Equal(t, strings.TrimSpace(`
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 ++ ++ ++ ++ ++ ++ ++ 00 00 00 00 00 00 00 00 00 00
00 00 00 ++ 01 01 01 01 01 ++ 00 00 00 00 00 00 00 00 00 00
00 ++ ++ ++ 01 01 01 01 01 ++ 00 00 ++ ++ ++ ++ ++ ++ ++ 00
00 ++ 02 ++ 01 01 01 01 01 ++ 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 02 ++ ++ ++ ++ ++ ++ ++ 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 ++ ++ ++ ++ ++ ++ ++ 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 ++ 02 ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ 00 00 00 00 00
00 ++ 02 ++ 05 05 05 05 05 05 05 05 05 05 ++ 00 00 00 00 00
00 ++ 02 ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++
00 ++ 02 ++ 00 00 00 00 00 00 00 00 00 ++ 07 07 07 07 07 ++
00 ++ 02 ++ 00 00 00 00 00 00 00 00 00 ++ 07 07 07 07 07 ++
00 ++ ++ ++ 00 00 00 00 00 00 00 00 00 ++ 07 07 07 07 07 ++
00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 07 07 07 07 07 ++
00 00 00 00 00 00 00 00 00 00 00 00 00 ++ 07 07 07 07 07 ++
00 00 00 00 00 00 00 00 00 00 00 00 00 ++ ++ ++ ++ ++ ++ ++
`), ff.String())
	ff.MarkPointUnreachable(&Point{0, 8}).Process()
	assert.Equal(t, strings.TrimSpace(`
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 ++ ++ ++ ++ ++ ++ ++ 00 00 00 00 00 00 00 00 00 00
00 00 00 ++ 01 01 01 01 01 ++ 00 00 00 00 00 00 00 00 00 00
00 ++ ++ ++ 01 01 01 01 01 ++ 00 00 ++ ++ ++ ++ ++ ++ ++ 00
00 ++ 02 ++ 01 01 01 01 01 ++ 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 02 ++ ++ ++ ++ ++ ++ ++ 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 02 ++ 00 00 00 00 00 00 00 00 ++ ++ ++ ++ ++ ++ ++ 00
++ ++ 02 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
05 ++ 02 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
05 ++ 02 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
05 ++ 02 ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ 00 00 00 00 00
05 ++ 02 ++ 06 06 06 06 06 06 06 06 06 06 ++ 00 00 00 00 00
05 ++ 02 ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++
05 ++ 02 ++ 05 05 05 05 05 05 05 05 05 ++ 08 08 08 08 08 ++
05 ++ 02 ++ 05 05 05 05 05 05 05 05 05 ++ 08 08 08 08 08 ++
05 ++ ++ ++ 05 05 05 05 05 05 05 05 05 ++ 08 08 08 08 08 ++
05 05 05 05 05 05 05 05 05 05 05 05 05 ++ 08 08 08 08 08 ++
05 05 05 05 05 05 05 05 05 05 05 05 05 ++ 08 08 08 08 08 ++
05 05 05 05 05 05 05 05 05 05 05 05 05 ++ ++ ++ ++ ++ ++ ++
`), ff.String())
	ff.MarkPointEmpty(&Point{3, 4}).MarkPointEmpty(&Point{2, 16}).Process()
	assert.Equal(t, strings.TrimSpace(`
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 ++ ++ ++ ++ ++ ++ ++ 00 00 00 00 00 00 00 00 00 00
00 00 00 ++ 01 01 01 01 01 ++ 00 00 00 00 00 00 00 00 00 00
00 ++ ++ ++ 01 01 01 01 01 ++ 00 00 ++ ++ ++ ++ ++ ++ ++ 00
00 ++ 01 01 01 01 01 01 01 ++ 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 01 ++ ++ ++ ++ ++ ++ ++ 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 01 ++ 00 00 00 00 00 00 00 00 ++ 03 03 03 03 03 ++ 00
00 ++ 01 ++ 00 00 00 00 00 00 00 00 ++ ++ ++ ++ ++ ++ ++ 00
++ ++ 01 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
01 ++ 01 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
01 ++ 01 ++ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
01 ++ 01 ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ 00 00 00 00 00
01 ++ 01 ++ 06 06 06 06 06 06 06 06 06 06 ++ 00 00 00 00 00
01 ++ 01 ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++ ++
01 ++ 01 ++ 01 01 01 01 01 01 01 01 01 ++ 08 08 08 08 08 ++
01 ++ 01 ++ 01 01 01 01 01 01 01 01 01 ++ 08 08 08 08 08 ++
01 ++ 01 ++ 01 01 01 01 01 01 01 01 01 ++ 08 08 08 08 08 ++
01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 08 08 08 08 08 ++
01 01 01 01 01 01 01 01 01 01 01 01 01 ++ 08 08 08 08 08 ++
01 01 01 01 01 01 01 01 01 01 01 01 01 ++ ++ ++ ++ ++ ++ ++
	`), ff.String())
}

func Benchmark3(b *testing.B) {
	ff := getLotOfZones()
	ff.Process()
	for n := 0; n < b.N; n++ {
		ff.MarkPointUnreachable(&Point{0, 8}).Process()
		ff.MarkPointEmpty(&Point{3, 4}).MarkPointEmpty(&Point{2, 16}).Process()
		ff.MarkPointEmpty(&Point{0, 8}).Process()
		ff.MarkPointUnreachable(&Point{3, 4}).MarkPointUnreachable(&Point{2, 16}).Process()
	}
}
