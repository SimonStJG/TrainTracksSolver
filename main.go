package main

import "log"

func main() {
	board := Board{
		metadata: &BoardMetadata{
			size: XY{ 10, 10},
			constraints: Constraints{
				yFixed: []uint8{2, 6, 6, 6, 3, 7, 7, 4, 1, 3},
				xFixed: []uint8{4, 4, 4, 3, 7, 4, 6, 4, 3, 6},
			},
			// These could be inferred from the grid of cells below, as they are always the only two tracks which
			// escape the grid
			start: XY {0, 3},
			end: XY {4, 9},
		},
		cells: [][]Cell {
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, SwTrack},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{NwTrack, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, EwTrack, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, SeTrack, Unknown, Unknown, Unknown, Unknown, Unknown},
		},
	}

	solution, err := Run(board)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Solution: %v", solution.prettyFormat())
}
