package main

import (
	"log"
	"errors"
)

// Singly linked list of the position of unknown cells
type unknownCellPosition struct {
	x    int
	y    int
	next *unknownCellPosition
}

// Board which tracks the next unknown cell to be filled
type boardWithPosition struct {
	board           Board
	nextUnknownCell *unknownCellPosition
}

// Run the train tracks solver
func Run(board Board) (solution Board, err error) {
	log.Printf("Starting board: %v", board.prettyFormat())

	firstUnknownCell := buildListOfUnknownCells(board)

	var stack = []boardWithPosition{{
		board:           board,
		nextUnknownCell: firstUnknownCell,
	}}

	for {
		head := stack[(len(stack) - 1)]
		stack = stack[:(len(stack) - 1)]
		solution, finished := generateValidMoves(head, &stack)

		if finished {
			return solution, nil
		}

		if len(stack) == 0 {
			return solution, errors.New("no solutions found")
		}
	}
}

// Build a singly linked list of all cells in board with value Unknown
func buildListOfUnknownCells(board Board) *unknownCellPosition {
	var head = &unknownCellPosition{x: -1, y: -1}
	var prev = head

	for y := 0; y < board.metadata.size.y; y++ {
		for x := 0; x < board.metadata.size.x; x++ {
			if board.cells[y][x] == Unknown {
				prev.next = &unknownCellPosition{x: x, y: y}
				prev = prev.next
			}
		}
	}
	return head.next
}

// Given a cell `cell`, which cells are valid in the adjacent cell, in direction `direction`.
// For example, if `cell` = `NoTrack`, then valid tracks to the west of this are: |, ⌜, ⌞, and x in other words:
// NsTrack  | SeTrack | NeTrack | NoTrack
func validAdjacentCells(cell Cell, direction Direction) Cell {
	switch {
	case cell == Unknown:
		return Unknown
	case (hasTrack[direction.reverse()] & cell) != 0:
		// has a connection to the cell we are interested in
		return hasTrack[direction]
	default:
		// doesn't have a connection to the cell we are interested in
		return (0x7F) & (^ hasTrack[direction])
	}
}

func generateValidMoves(head boardWithPosition, stack *[]boardWithPosition) (solution Board, finished bool) {
	// Apply xFixed and yFixed constraints
	cellOptions := applyHorizontalConstraint(head.board, head.nextUnknownCell.y) &
		applyVerticalConstraint(head.board, head.nextUnknownCell.x)
	if cellOptions == 0 {
		return solution, false
	}

	// Apply constraints from adjacent cells
	var direction Direction;
	for direction = 0; direction < 4; direction ++ {
		cellOptions = cellOptions & validAdjacentCells(head.board.valueAtOffset(head.nextUnknownCell.x, head.nextUnknownCell.y, direction), direction)
		if cellOptions == 0 {
			return solution, false
		}
	}

	// For each option encoded in cellOptions, push the corresponding board onto the stack.
	// If any board is complete then check it doesn't contain any cycles, and if it doesn't return it as a solution.
	var i uint8
	for i = 0; i <= 7; i++ {
		option := cellOptions & (1 << i)
		if option != 0 {
			newBoard := head.board.copy()
			newBoard.cells[head.nextUnknownCell.y][head.nextUnknownCell.x] = option

			if head.nextUnknownCell.next == nil {
				// No unknown cells left
				newBoardHasCycles := hasCycles(newBoard)
				if newBoardHasCycles {
					log.Printf("(Found bad solution with cycles: %v)", newBoard.prettyFormat())
					break
				} else {
					return newBoard, true
				}
			}

			*stack = append(*stack, boardWithPosition{
				board:           newBoard,
				nextUnknownCell: head.nextUnknownCell.next,
			})
		}
	}

	return solution, false
}

// Check if the given board has any cycles (any closed loops), these make the solution invalid.
func hasCycles(board Board) bool {
	// Count number of tracks on the board
	tracksOnBoard := 0
	for y := 0; y < board.metadata.size.y; y++ {
		for x := 0; x < board.metadata.size.x; x++ {
			if board.cells[y][x] != NoTrack {
				tracksOnBoard ++
			}
		}
	}

	// Walk around the track until we get to the end or we hit the start again
	current := board.metadata.start
	tracksSeen := 1
	var prevDirection Direction = 5 // Mustn't be a valid direction
	for {
		cell := board.valueAt(current)
		// Pick a direction to walk, this can't be the direction you just came from
		directionOptions := hasDirection[cell]
		newDirection := directionOptions[0]
		if newDirection == prevDirection.reverse() {
			newDirection = directionOptions[1]
		}

		// Walk in that direction
		switch newDirection {
		case north:
			current.y --
		case east:
			current.x ++
		case south:
			current.y ++
		case west:
			current.x --
		}
		prevDirection = newDirection
		tracksSeen ++

		// Check if we are at the end, and if so check if we've seen all the tracks yet
		if board.metadata.end == current {
			if tracksSeen == tracksOnBoard {
				return false
			} else {
				return true
			}
		}
	}
}

func applyVerticalConstraint(board Board, x int) Cell {
	vConstraint := board.metadata.constraints.xFixed[x]
	var lowerLimit uint8 = 0
	var upperLimit uint8 = 0
	for y := 0; y < board.metadata.size.y; y++ {
		switch square := board.cells[y][x]; square {
		case Unknown:
			upperLimit ++
		case NoTrack:
			// nothing
		default:
			lowerLimit++
			upperLimit++
		}
	}
	if vConstraint == lowerLimit {
		return NoTrack
	}
	if vConstraint == upperLimit {
		return ^ NoTrack
	}
	return Unknown
}

func applyHorizontalConstraint(board Board, y int) Cell {
	hConstraint := board.metadata.constraints.yFixed[y]
	var lowerLimit uint8 = 0
	var upperLimit uint8 = 0
	for x := 0; x < board.metadata.size.x; x++ {
		switch square := board.cells[y][x]; square {
		case Unknown:
			upperLimit ++
		case NoTrack:
			// nothing
		default:
			lowerLimit++
			upperLimit++
		}
	}
	if hConstraint == lowerLimit {
		return NoTrack
	}
	if hConstraint == upperLimit {
		return ^ NoTrack
	}
	return Unknown
}
