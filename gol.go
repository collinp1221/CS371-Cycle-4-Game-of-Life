// An implementation of Conway's Game of Life.
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

// Field represents a two-dimensional field of cells.
type Field struct {
	s        [][]bool
	width, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(width, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, width)
	}
	return &Field{s: s, width: width, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.width
	x %= f.width
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b     *Field
	width, h int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(width, h int) *Life {
	a := NewField(width, h)
	for i := 0; i < (width * h / 4); i++ {
		a.Set(rand.Intn(width), rand.Intn(h), true)
	}
	return &Life{
		a: a, b: NewField(width, h),
		width: width, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (grid *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < grid.h; y++ {
		for x := 0; x < grid.width; x++ {
			grid.b.Set(x, y, grid.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	grid.a, grid.b = grid.b, grid.a
}

// String returns the game board as a string.
func (grid *Life) String() string {
	var buf bytes.Buffer
	for y := 0; y < grid.h; y++ {
		for x := 0; x < grid.width; x++ {
			b := byte(' ')
			if grid.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func main() {
	grid := NewLife(40, 15)
	for i := 0; i < 1000; i++ {
		grid.Step()
		fmt.Print("\x0c", grid)     // Clear screen and print field.
		time.Sleep(time.Second / 5) //Number here controls the number of "frames" per second
	}
}
