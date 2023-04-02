package floodfill

import (
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
