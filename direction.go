package main

type Direction uint8

const (
	north Direction = 0
	east  Direction = 1
	south Direction = 2
	west  Direction = 3
)

func (direction Direction) reverse() Direction {
	return (direction + 2) % 4
}
