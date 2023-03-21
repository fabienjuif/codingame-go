package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var (
	gridsMu        = sync.Mutex{}
	gridsMap       = make(map[string]*Grid, 100)
	playerStampReg = regexp.MustCompile(`\[(\d+):(\d+)\]`)
)

func main() {
	var grid *Grid
	scan := NewScanner()
	for {
		// N: total number of players (2 to 4).
		// P: your player number (0 to 3).
		var N, P int
		scan(&N, &P)
		if grid == nil {
			grid = NewGrid(N)
		}
		// fmt.Fprintf(os.Stderr, "N,P: %d %d\n", N, P)

		for i := 0; i < N; i++ {
			// X0: starting X coordinate of lightcycle (or -1)
			// Y0: starting Y coordinate of lightcycle (or -1)
			// X1: starting X coordinate of lightcycle (can be the same as X0 if you play before this player)
			// Y1: starting Y coordinate of lightcycle (can be the same as Y0 if you play before this player)
			var X0, Y0, X1, Y1 int
			scan(&X0, &Y0, &X1, &Y1)
			// fmt.Fprintf(os.Stderr, "putOne(a[0], %d)\nputOne(a[1], %d)\nputOne(a[2], %d)\nputOne(a[3], %d)\n", X0, Y0, X1, Y1)
			grid.MarkPlayer(i, X0, Y0)
			grid.MarkPlayer(i, X1, Y1)
		}
		// fmt.Fprintf(os.Stderr, "stamp: %v\n", grid.GetStamp())
		fmt.Fprintf(os.Stderr, "%v\n", grid)
		direction, _ := GestBestDirection(P, grid)
		fmt.Println(direction)
	}
}

func GestBestDirection(P int, grid *Grid) (Direction, float64) {
	wg := sync.WaitGroup{}
	var lS, rS, dS, uS float64
	wg.Add(4)
	go func() {
		defer wg.Done()
		g, cell := grid.GoLeft(P)
		if cell != nil {
			lS = g.GetPlayerScore(P)
		}
	}()
	go func() {
		defer wg.Done()
		g, cell := grid.GoDown(P)
		if cell != nil {
			dS = g.GetPlayerScore(P)
		}
	}()
	go func() {
		defer wg.Done()
		g, cell := grid.GoRight(P)
		if cell != nil {
			rS = g.GetPlayerScore(P)
		}
	}()
	go func() {
		defer wg.Done()
		g, cell := grid.GoUp(P)
		if cell != nil {
			uS = g.GetPlayerScore(P)
		}
	}()
	wg.Wait()
	m := max(lS, rS, dS, uS)
	fmt.Fprintf(os.Stderr, "lS:%g|rS:%g|dS:%g|uS:%g>%g >", lS, rS, dS, uS, m)
	if m == lS {
		return Direction_Left, m
	}
	if m == rS {
		return Direction_Right, m
	}
	if m == dS {
		return Direction_Down, m
	}
	return Direction_Up, m
}

// return true if already exists
func AddGrid(grid *Grid) bool {
	gridsMu.Lock()
	defer gridsMu.Unlock()
	_, exists := gridsMap[grid.GetStamp()]
	if exists {
		return true
	}
	gridsMap[grid.GetStamp()] = grid
	return false
}

const (
	GridWidth  = 30
	GridHeight = 20
)

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
//
//	GRID
//
// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
type Grid struct {
	Players []*GridPlayer
	Cells   []*Cell
}

func NewGrid(players int) *Grid {
	cells := make([]*Cell, GridWidth*GridHeight)

	for i := range cells {
		cells[i] = NewCell(i, i%GridWidth, i/GridWidth)
	}

	return &Grid{
		Cells:   cells,
		Players: make([]*GridPlayer, players),
	}
}

func NewGridFromStamp(stamp string, playersStamp string) *Grid {
	cells := make([]*Cell, GridWidth*GridHeight)
	cellIdx := 0
	lP := PlayerName_Unknown
	nS := ""
	// cells
	computeCells := func() {
		n := StringToInt(nS)
		for j := 0; j < n; j += 1 {
			cells[cellIdx] = NewCell(cellIdx, cellIdx%GridWidth, cellIdx/GridWidth)
			if lP != PlayerName_Unknown {
				cells[cellIdx].MarkFull(&lP)
			}
			cellIdx += 1
		}
	}
	for i, r := range stamp {
		rS := string(r)
		if SliceIncludes([]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}, r) {
			nS += rS
		} else {
			p := NewPlayerName(rS)
			if i != 0 {
				computeCells()
			}
			nS = ""
			lP = p
		}
	}
	computeCells()
	// players
	res := playerStampReg.FindAllStringSubmatch(playersStamp, -1)
	players := make([]*GridPlayer, len(res))
	for i, m := range res {
		players[i] = &GridPlayer{StringToInt(m[1]), StringToInt(m[2])}
	}
	return &Grid{
		Cells:   cells,
		Players: players,
	}
}

func (g *Grid) MarkPlayer(n, x, y int) {
	g.Players[n] = &GridPlayer{x, y}
	g.MarkFull(NewPlayerNameFromN(n), x, y)
}

// GetPlayerScore - FloodFill score given the current player position
func (g *Grid) GetPlayerScore(n int) float64 {
	if n >= len(g.Players) || n < 0 {
		log.Fatalf("player not found: %d", n)
	}

	playerPos := g.Players[n]
	playerCell := g.GetCell(playerPos.X, playerPos.Y)
	if playerCell == nil {
		log.Fatalf("player is a not found cell: %v", playerPos)
	}

	score := 0.0
	indexesToVisit := []int{playerCell.Index}
	alreadyAdded := map[int]bool{playerCell.Index: true}
	getNextToVisit := func() (int, bool) {
		if len(indexesToVisit) <= 0 {
			return 0, false
		}
		n, indexesToVisit = indexesToVisit[0], indexesToVisit[1:]
		return n, true
	}
	addToNextToVisit := func(cell *Cell) {
		if cell != nil {
			_, exists := alreadyAdded[cell.Index]
			if !exists {
				alreadyAdded[cell.Index] = true
				indexesToVisit = append(indexesToVisit, cell.Index)
			}
		}
	}

	addToNextToVisit(g.GetCell(playerCell.X+1, playerCell.Y))
	addToNextToVisit(g.GetCell(playerCell.X-1, playerCell.Y))
	addToNextToVisit(g.GetCell(playerCell.X, playerCell.Y-1))
	addToNextToVisit(g.GetCell(playerCell.X, playerCell.Y+1))

	for {
		cellIndex, exists := getNextToVisit()
		if !exists {
			break
		}
		cell := g.Cells[cellIndex]
		if cell == nil {
			log.Fatalf("Cell not found: %d", cellIndex)
		}
		if cell.IsVisitable() {
			score += 1
			addToNextToVisit(g.GetCell(cell.X+1, cell.Y))
			addToNextToVisit(g.GetCell(cell.X-1, cell.Y))
			addToNextToVisit(g.GetCell(cell.X, cell.Y-1))
			addToNextToVisit(g.GetCell(cell.X, cell.Y+1))
		}
	}

	return score
}

// GetPlayerScores - get all player score using GetPlayerScore(n)
//
// result is an array of each score matching player order
func (g *Grid) GetPlayerScores() []float64 {
	wg := sync.WaitGroup{}
	scores := make([]float64, len(g.Players))
	for n := range g.Players {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			scores[n] = g.GetPlayerScore(n)
		}(n)
	}
	wg.Wait()
	return scores
}

func (g *Grid) GetStamp() string {
	s := ""
	n := 0
	var last PlayerName
	for i, c := range g.Cells {
		n += 1
		if i == 0 {
			last = c.Player
			n = 0
		} else {
			if last != c.Player {
				if last == PlayerName_Unknown {
					s += fmt.Sprintf("%s%d", CellType_Empty, n)
				} else {
					s += fmt.Sprintf("%s%d", last, n)
				}
				n = 0
				last = c.Player
			}
		}
	}
	if last == PlayerName_Unknown {
		return fmt.Sprintf("%s%s%d", s, CellType_Empty, n+1)
	}
	return fmt.Sprintf("%s%s%d", s, last.slug, n+1)
}

func (g *Grid) MarkFull(player PlayerName, x, y int) (*Grid, *Cell) {
	cell := g.GetCell(x, y)
	if cell != nil && cell.IsVisitable() {
		cell.MarkFull(&player)
	}
	return g, cell
}

func (g *Grid) IsCellVisitable(x, y int) bool {
	cell := g.GetCell(x, y)
	return cell != nil && cell.Type == CellType_Empty
}

func (g *Grid) Clone() *Grid {
	// // TODO: use a pool?
	cells := make([]*Cell, len(g.Cells))
	copy(cells, g.Cells)
	// // TODO: use a pool?
	players := make([]*GridPlayer, len(g.Players))
	copy(players, g.Players)
	return &Grid{
		Players: players,
		Cells:   cells,
	}
}

// SetCell - Replace the given cell in the grid
//
// MUTATE the grid
func (g *Grid) SetCell(cell *Cell) (*Grid, *Cell) {
	if cell.Index < 0 || cell.Index >= len(g.Cells) {
		log.Fatalf("cell do not exist: %v\n", cell.Index)
	}
	g.Cells[cell.Index] = cell
	return g, cell
}

// GoRight - the given player index go to the right (if he can) and give a new grid
func (g *Grid) GoRight(n int) (*Grid, *Cell) {
	player := g.GetPlayer(n)
	return g.Go(n, player.X+1, player.Y)
}

// GoLeft - the given player index go to the right (if he can) and give a new grid
func (g *Grid) GoLeft(n int) (*Grid, *Cell) {
	player := g.GetPlayer(n)
	return g.Go(n, player.X-1, player.Y)
}

// GoUp - the given player index go to the right (if he can) and give a new grid
func (g *Grid) GoUp(n int) (*Grid, *Cell) {
	player := g.GetPlayer(n)
	return g.Go(n, player.X, player.Y-1)
}

// GoDown - the given player index go to the right (if he can) and give a new grid
func (g *Grid) GoDown(n int) (*Grid, *Cell) {
	player := g.GetPlayer(n)
	return g.Go(n, player.X, player.Y+1)
}

// Go - the given player index go to the specified location (if he can) and give a new grid & cell
func (g *Grid) Go(n, X, Y int) (*Grid, *Cell) {
	if n >= len(g.Players) || n < 0 {
		log.Fatalf("player do not exist: %v\n", n)
	}
	playerName := NewPlayerNameFromN(n)
	cell := g.GetCell(X, Y)
	if cell == nil || !cell.IsVisitable() {
		return g, nil
	}

	// new grid can be spawn, we clone the cell and the grid to mutate them
	clonedCell := cell.Clone()
	clonedCell.MarkFull(&playerName)
	clonedGrid := g.Clone()
	clonedPlayer := clonedGrid.Players[n].Clone()
	clonedGrid.Players[n] = clonedPlayer
	clonedPlayer.X = clonedCell.X
	clonedPlayer.Y = clonedCell.Y

	return clonedGrid.SetCell(clonedCell)
}

func (g *Grid) GetPlayer(n int) *GridPlayer {
	if n >= len(g.Players) || n < 0 {
		log.Fatalf("player do not exist: %v\n", n)
	}
	return g.Players[n]
}

func (g *Grid) GetCell(x, y int) *Cell {
	if x >= GridWidth || x < 0 {
		return nil
	}
	if y >= GridHeight || y < 0 {
		return nil
	}
	return g.Cells[y*GridWidth+x]
}

func (g *Grid) GetPlayersStamp() string {
	s := ""
	for _, p := range g.Players {
		s += fmt.Sprintf("[%d:%d]", p.X, p.Y)
	}
	return s
}

func (g *Grid) String() string {
	s := fmt.Sprintf("stamp: %s\n", g.GetStamp())
	s += fmt.Sprintf("players: %s\n", g.GetPlayersStamp())
	for y := 0; y < GridHeight; y += 1 {
		for x := 0; x < GridWidth; x += 1 {
			s += g.GetCell(x, y).String()
			// if x < GridWidth-1 {
			// 	s += ","
			// }
		}
		if y < GridHeight-1 {
			s += "\n"
		}
	}
	return s
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
//
//	GRID PLAYER
//
// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
type GridPlayer struct {
	X, Y int
}

func (g *GridPlayer) Clone() *GridPlayer {
	return &GridPlayer{
		X: g.X,
		Y: g.Y,
	}
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
//
//	CELL
//
// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
type Cell struct {
	Index, X, Y int
	Type        CellType
	Player      PlayerName
}

func NewCell(index, x, y int) *Cell {
	return &Cell{
		Index:  index,
		X:      x,
		Y:      y,
		Type:   CellType_Empty,
		Player: PlayerName_Unknown,
	}
}

func (c *Cell) MarkFull(player *PlayerName) *Cell {
	c.Type = CellType_Full
	if player != nil {
		c.Player = *player
	}
	return c
}

func (c *Cell) IsVisitable() bool {
	return c.Type == CellType_Empty
}

func (c *Cell) Clone() *Cell {
	return &Cell{
		Index:  c.Index,
		X:      c.X,
		Y:      c.Y,
		Type:   c.Type,
		Player: c.Player,
	}
}

func (c *Cell) String() string {
	if c.Player != PlayerName_Unknown {
		return fmt.Sprintf("%v", c.Player)
	}
	return fmt.Sprintf("%v", c.Type)
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
//
//	ENUMS
//
// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
type PlayerName struct{ slug string }

func (p PlayerName) String() string {
	return p.slug
}

func NewPlayerName(slug string) PlayerName {
	switch slug {
	case PlayerName_A.slug:
		return PlayerName_A
	case PlayerName_B.slug:
		return PlayerName_B
	case PlayerName_C.slug:
		return PlayerName_C
	case PlayerName_D.slug:
		return PlayerName_D
	default:
		return PlayerName_Unknown

	}
}

func NewPlayerNameFromN(n int) PlayerName {
	switch n {
	case 0:
		return PlayerName_A
	case 1:
		return PlayerName_B
	case 2:
		return PlayerName_C
	case 3:
		return PlayerName_D
	default:
		return PlayerName_Unknown
	}
}

var (
	PlayerName_Unknown = PlayerName{"UNKNOWN"}
	PlayerName_A       = PlayerName{"A"}
	PlayerName_B       = PlayerName{"B"}
	PlayerName_C       = PlayerName{"C"}
	PlayerName_D       = PlayerName{"D"}
)

type CellType struct{ slug string }

func (t CellType) String() string {
	return t.slug
}

func NewCellType(slug string) CellType {
	switch slug {
	case CellType_Empty.slug:
		return CellType_Empty
	case CellType_Full.slug:
		return CellType_Full
	default:
		return CellType_Unknown
	}
}

var (
	CellType_Unknown = CellType{"UNKNOWN"}
	CellType_Empty   = CellType{"."}
	CellType_Full    = CellType{"F"}
)

type Direction struct{ slug string }

func (d Direction) String() string {
	return d.slug
}

func NewDirection(slug string) Direction {
	switch slug {
	case Direction_Left.slug:
		return Direction_Left
	case Direction_Right.slug:
		return Direction_Right
	case Direction_Down.slug:
		return Direction_Down
	case Direction_Up.slug:
		return Direction_Up
	default:
		return Direction_Unknown
	}
}

var (
	Direction_Unknown = Direction{"UNKNOWN"}
	Direction_Left    = Direction{"LEFT"}
	Direction_Right   = Direction{"RIGHT"}
	Direction_Down    = Direction{"DOWN"}
	Direction_Up      = Direction{"UP"}
)

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
//
//	INPUTS
//
// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
func NewScanner() func(a ...any) (int, error) {
	_, local := os.LookupEnv("L")
	if local {
		callIdx := 0
		putOne := func(a any, i int) {
			*a.(*int) = i
		}
		return func(a ...any) (int, error) {
			n := 0
			if callIdx%4 == 0 {
				n = 2
				putOne(a[0], 3)
				putOne(a[1], 0)
			} else {
				switch callIdx {
				case 1:
					n = 4
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 5)
					putOne(a[3], 0)
				case 2:
					n = 4
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 14)
				case 3:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 29)
					putOne(a[3], 2)
				case 5:
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 4)
					putOne(a[3], 0)
				case 6:
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 13)
				case 7:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 28)
					putOne(a[3], 2)
				case 9:
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 3)
					putOne(a[3], 0)
				case 10:
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 12)
				case 11:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 27)
					putOne(a[3], 2)

				case 13:
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 2)
					putOne(a[3], 0)
				case 14:
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 11)
				case 15:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 26)
					putOne(a[3], 2)
				case 17:
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 1)
					putOne(a[3], 0)
				case 18:
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 10)
				case 19:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 25)
					putOne(a[3], 2)
				case 21:
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 0)
					putOne(a[3], 0)
				case 22:
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 9)
				case 23:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 24)
					putOne(a[3], 2)
				case 25:
					putOne(a[0], 5)
					putOne(a[1], 0)
					putOne(a[2], 0)
					putOne(a[3], 1)
				case 26:
					putOne(a[0], 19)
					putOne(a[1], 14)
					putOne(a[2], 19)
					putOne(a[3], 8)
				case 27:
					putOne(a[0], 29)
					putOne(a[1], 2)
					putOne(a[2], 23)
					putOne(a[3], 2)
				default:
				}
			}
			callIdx += 1
			return n, nil
		}
	}
	return fmt.Scan
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
//
//	UTIL
//
// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
func max(a ...float64) float64 {
	init := false
	max := 0.0
	for _, v := range a {
		if !init {
			max = v
			init = true
		} else {
			if v > max {
				max = v
			}
		}
	}
	return max
}

func SliceIncludes[T comparable](slice []T, v T) bool {
	for _, d := range slice {
		if d == v {
			return true
		}
	}
	return false
}

func StringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return n
}
