package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCanBeExtendedInDirection(t *testing.T) {
	assert.Equal(
		t,
		NsTrack | SeTrack | NeTrack | NoTrack,
		validAdjacentCells(NoTrack, west),
	)

	assert.Equal(
		t,
		EwTrack | SeTrack | SwTrack | NoTrack,
		validAdjacentCells(EwTrack, north),
	)

	assert.Equal(
		t,
		NsTrack | NeTrack | NwTrack,
		validAdjacentCells(NsTrack, north),
	)
}

func TestSolutions(t *testing.T) {
	board := Board{
		cells: [][]Cell{
			{Unknown, Unknown, Unknown, SwTrack, Unknown, Unknown},
			{SwTrack, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown},
			{Unknown, Unknown, SeTrack, Unknown, Unknown, Unknown},
		},
		metadata: &BoardMetadata{
			size: XY {6, 6},
			constraints: Constraints{
				yFixed: [] uint8{3, 3, 5, 3, 3, 4},
				xFixed: [] uint8{4, 5, 2, 4, 2, 4},
			},
			start: XY{ 0, 1},
			end: XY {2, 5},
		},
	}

	solution, err := Run(board)
	assert.Equal(t, err, nil)
	assert.Equal(
		t,
		solution.prettyFormat(),
		`
 452424
3x⌜-⌝xx
3⌝|x|xx
5||x⌞-⌝
3||xxx|
3⌞⌟xxx|
4xx⌜--⌟
`)
}
