package floodfill

type Point struct {
	X, Y int
}

func (p *Point) ToIndex(width int) int {
	return FromPointToIndex(p.X, p.Y, width)
}

func FromPointToIndex(x, y, width int) int {
	return y*width + x
}
