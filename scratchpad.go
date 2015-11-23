package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const cellSize = 10

func main() {
	doc := js.Global.Get("document")
	//Clear the page
	doc.Get("body").Set("innerHTML", "")
	doc.Set("title", "Canvas Tests")

	canvas := doc.Call("createElement", "canvas")
	canvas.Set("height", "500")
	canvas.Set("width", "500")
	canvas.Set("id", "myCanvas")
	doc.Get("body").Call("appendChild", canvas)

	ctx := canvas.Call("getContext", "2d")
	ctx.Set("fillStyle", "red")

	ctx.Call("fillRect", 0, 0, 50, 50)

}

type Game struct {
	Board
}

type Board struct {
	 state, prevState [][]Cell
	 Rows, Cols uint
}

func New(rows, cols uint) Board {
	
	var b Board
	b.Rows, b.Cols = rows, cols
	
	b.state = make([][]Cell, rows+2) //+2 adds a border 1 Cell wide at the edge of the board for the purposes of counting
	b.prevState = make([][]Cell, rows+2) //We'll get the state of each generation by looking at the previous generation
	
	//Make a blank board
	for y := 0; y < len(state); y++ {
		state[y] = make([]Cell, cols+2)
		prevState[y] = make([]Cell, cols+2)
	}
	
	return b
}

type Cell struct {
	Alive bool 
	Player int
}

