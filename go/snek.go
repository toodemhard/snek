package main

import (
	"math/rand"
	"snek/pkg/ring_array"
)

type cellState uint8

const (
	emptyCell cellState = 0
	snakeCell cellState = 1
	foodCell  cellState = 2
)

const (
	gridWidth  int = 12
	gridHeight int = 9
	gridSize       = gridWidth * gridHeight
	initSize   int = 5
)

type vector struct {
	x int
	y int
}

func vectorAdd(a vector, b vector) vector {
	return vector{a.x + b.x, a.y + b.y}
}

func vectorEquals(a vector, b vector) bool {
	return (a.x == b.x) && (a.y == b.y)
}

type model struct {
	currentView view

	food vector

	snake       ring_array.RingArray[vector]
	direction   vector
	inputBuffer ring_array.RingArray[vector]
}

func (m *model) grow() {
	tail, _ := m.snake.Tail()
	m.snake.PushBack(tail)
}

func outOfBounds(position vector) bool {
	if position.x > gridWidth-1 || position.x < 0 {
		return true
	}
	if position.y > gridHeight-1 || position.y < 0 {
		return true
	}
	return false
}

func collision(snake ring_array.RingArray[vector]) bool {
	cells := snake.Array()
	for i := len(cells) - 2; i >= 0; i-- {
		if vectorEquals(cells[len(cells)-1], cells[i]) {
			return true
		}
	}
	return false
}

func randomEmpty(m model) vector {
	grid := stateGrid(m)
	var empty vector
	found := false
	for found == false {
		empty = vector{rand.Intn(gridWidth), rand.Intn(gridHeight)}
		if grid[gridIndex(empty)] == emptyCell {
			found = true
		}
	}
	return empty
}

func stateGrid(m model) []cellState {
	grid := make([]cellState, gridSize)
	for _, snakeNode := range m.snake.Array() {
		grid[gridIndex(snakeNode)] = snakeCell
	}
	grid[gridIndex(m.food)] = foodCell
	return grid
}

func tickUpdate(m model) (model, bool) {
	input, err := m.inputBuffer.PopBack()
	if err == nil {
		m.direction = input
	}
	head, _ := m.snake.Head()
	newHead := vectorAdd(head, m.direction)

	m.snake.PushFront(newHead)
	m.snake.PopBack()

	if outOfBounds(newHead) || collision(m.snake) {
		m.currentView = gameOver
		return m, true
	}

	if vectorEquals(newHead, m.food) {
		m.grow()
		m.food = randomEmpty(m)
	}

	return m, false
}

func gridIndex(vector vector) int {
	return vector.y*gridWidth + vector.x
}
