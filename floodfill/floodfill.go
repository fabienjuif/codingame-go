package floodfill

import (
	"fmt"
	"math"
)

type FloodFill interface {
	// Call this function before reading values.
	//
	// You can call it after mutations are done, eg:
	//   ff.MarkPointEmpty(point).MarkPointUnreachable(point2).Process()
	Process()
	GetZoneFromPoint(p *Point) Zone
	GetZoneFromIndex(i int) Zone
	MarkPointEmpty(p *Point) FloodFill
	MarkPointUnreachable(p *Point) FloodFill
	MarkIndexEmpty(i int) FloodFill
	MarkIndexUnreachable(i int) FloodFill
	GetCells() []Zone
	String() string
}

type floodfill struct {
	// cells zones
	//  - 255 means empty (default value)
	//  - 254 means unreachable
	//  - from 0, there are zones
	cells []byte
	zones     byte
	width     int
	height    int
}

func (f *floodfill) String() string {
	table := ``
	for y := 0; y < f.height; y += 1 {
		row := ``
		for x := 0; x < f.width; x += 1 {
			cell := f.cells[FromPointToIndex(x, y, f.width)]
			if cell == cellType_Empty {
				row += ".."
			} else if cell == cellType_Unreachable {
				row += "++"
			} else {
				row += fmt.Sprintf("%02v", cell)
			}
			if x < f.width-1 {
				row += " "
			}
		}
		table += row
		if y < f.height-1 {
			table += "\n"
		}
	}
	return table
}

func (f *floodfill) GetCells() []byte {
	return f.cells
}

func (f *floodfill) MarkIndexEmpty(i int) FloodFill {
	f.cells[i] = byte(cellType_Empty)
	return f
}

func (f *floodfill) MarkIndexUnreachable(i int) FloodFill {
	f.cells[i] = byte(cellType_Unreachable)
	return f
}

func (f *floodfill) MarkPointEmpty(p *Point) FloodFill {
	return f.MarkIndexEmpty(p.ToIndex(f.width))
}

// MarkPointUnreachable implements FloodFill
func (f *floodfill) MarkPointUnreachable(p *Point) FloodFill {
	return f.MarkIndexUnreachable(p.ToIndex(f.width))
}

func (f *floodfill) GetZoneFromIndex(i int) Zone {
	return f.cells[i]
}

func (f *floodfill) GetZoneFromPoint(p *Point) Zone {
	return f.cells[p.ToIndex(f.width)]
}

func (f *floodfill) Process() {
	curZone := byte(0)
	for y := 0; y < f.height; y += 1 {
		for x := 0; x < f.width; x += 1 {
			curIdx := FromPointToIndex(x, y, f.width)
			cell := f.cells[curIdx]
			if cell == cellType_Unreachable {
				continue
			}
			var lValue *byte
			var tValue *byte
			if x > 0 {
				left := f.cells[FromPointToIndex(x-1, y, f.width)]
				if left != cellType_Unreachable {
					lValue = &left
				}
			}
			if y > 0 {
				top := f.cells[FromPointToIndex(x, y-1, f.width)]
				if top != cellType_Unreachable {
					tValue = &top
				}
			}

			if lValue != nil && tValue != nil {
				f.cells[curIdx] = *lValue
				if lValue != tValue {
					f.swap(*lValue, *tValue)
				}
			} else if lValue != nil {
				f.cells[curIdx] = *lValue
			} else if tValue != nil {
				f.cells[curIdx] = *tValue
			} else {
				f.cells[curIdx] = curZone
				curZone += 1
			}
		}
	}
}

func (f *floodfill) swap(val1, val2 byte) {
	min := val1
	max := val2
	if val1 > val2 {
		min = val2
		max = val1
	}
	for i, c := range f.cells {
		if c == max {
			f.cells[i] = min
		}
	}
}

func NewFloodFill(width, height int) FloodFill {
	cells := make([]byte, width*height)
	for i := range cells {
		cells[i] = cellType_Empty
	}
	return &floodfill{
		cells:  cells,
		zones:  0,
		width:  width,
		height: height,
	}
}

type Zone = byte

var (
	cellType_Empty       = byte(math.MaxUint8)
	cellType_Unreachable = byte(math.MaxUint8 - 1)
)
