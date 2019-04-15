package main

import "fmt"

type (
	// The contents of a cell (one possible square on a board)
	Cell uint8

	XY struct {
		x int
		y int
	}

	// The number of cells in a row (yFixed) or column (xFixed), these are usually printed around the outside of the
	// puzzle.
	Constraints struct {
		xFixed []uint8
		yFixed []uint8
	}

	BoardMetadata struct {
		constraints Constraints
		// Size of the board
		size XY
		// Starting cell
		start XY
		// Ending cell
		end XY
	}
	Board struct {
		cells [][]Cell
		// Pointer to the metadata - this is always the same throughout the algorithm, so use a pointer to avoid
		// unnecessary copying.
		metadata *BoardMetadata
	}
)

// Cell contents, we use a bitmask so we an reuse this as an "options for the cell contents", and then do lots of OR-ing
// of the cell contents in the algorithm
const (
	NoTrack Cell = 1 << 0 // 1
	EwTrack Cell = 1 << 1 // 2 -
	NsTrack Cell = 1 << 2 // 4 |
	SeTrack Cell = 1 << 3 // 8 ⌜
	SwTrack Cell = 1 << 4 // 16 ⌝
	NeTrack Cell = 1 << 5 // 32 ⌞
	NwTrack Cell = 1 << 6 // 64 ⌟

	Unknown = NoTrack | EwTrack | NsTrack | SeTrack | SwTrack | NeTrack | NwTrack
)

var (
	hasTrack = map[Direction]Cell{
		north: NsTrack | NeTrack | NwTrack,
		south: NsTrack | SeTrack | SwTrack,
		east:  EwTrack | NeTrack | SeTrack,
		west:  EwTrack | NwTrack | SwTrack,
	}
	hasDirection = map[Cell][2]Direction{
		EwTrack: {east, west},
		NsTrack: {north, south},
		SeTrack: {south, east},
		SwTrack: {south, west},
		NeTrack: {north, east},
		NwTrack: {north, west},
	}
)

// Copy the board.  This makes a deep copy of the cells, but keeps the same pointer to the metadata.
func (board Board) copy() Board {
	cells := make([][]Cell, board.metadata.size.y, board.metadata.size.y)
	for y := range cells {
		cells[y] = make([]Cell, board.metadata.size.x, board.metadata.size.y)
		copy(cells[y], board.cells[y])
	}
	return Board{
		cells:    cells,
		metadata: board.metadata,
	}
}

func (board Board) valueAt(xy XY) Cell {
	return board.cells[xy.y][xy.x]
}

// Return the contents of the cell at `point` with an offset in `direction`.  If this is off the edge of the board, then
// return NoTrack.  For example, on the board
//  x⌜-
//  ⌝|x
//  ⌞⌟x
// valueAtOffset(0,0, east) = ⌜
func (board Board) valueAtOffset(x int, y int, direction Direction) Cell {
	switch direction {
	case north:
		if y == 0 {
			return NoTrack
		}
		return board.cells[y-1][x]
	case west:
		if x == 0 {
			return NoTrack
		}
		return board.cells[y][x-1]
	case south:
		if y == board.metadata.size.y-1 {
			return NoTrack
		}
		return board.cells[y+1][x]
	case east:
		if x == board.metadata.size.x-1 {
			return NoTrack
		}
		return board.cells[y][x+1]
	default:
		panic("Illegal direction")
	}
}

// Format the board and constraints prettily as a string
func (board Board) prettyFormat() string {
	pretty := "\n"

	pretty += " "
	for x := 0; x < board.metadata.size.x; x++ {
		pretty += fmt.Sprintf("%v", board.metadata.constraints.xFixed[x])
	}
	pretty += "\n"

	for y := 0; y < board.metadata.size.y; y++ {
		pretty += fmt.Sprintf("%v", board.metadata.constraints.yFixed[y])
		for x := 0; x < board.metadata.size.x; x++ {
			switch board.cells[y][x] {
			case NoTrack:
				pretty += "x"
			case EwTrack:
				pretty += "-"
			case NsTrack:
				pretty += "|"
			case SeTrack:
				pretty += "⌜"
			case SwTrack:
				pretty += "⌝"
			case NeTrack:
				pretty += "⌞"
			case NwTrack:
				pretty += "⌟"
			case Unknown:
				pretty += "?"
			default:
				pretty += "Z"
			}
		}
		pretty += "\n"
	}
	return pretty
}
