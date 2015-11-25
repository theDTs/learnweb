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

//Game describes a Game of Life "board", with a grid of cells that can be updated turn by turn
type Game struct {
	state, prevState [][]Cell
	rows, cols       uint
}

//NewGame returns a game with an empty board
func NewGame(rows, cols uint) *Game {
	var g Game
	g.rows, g.cols = rows, cols

	g.state = make([][]Cell, rows+2)     //+2 adds a border 1 Cell wide at the edge of the board for the purposes of counting
	g.prevState = make([][]Cell, rows+2) //We'll get the state of each generation by looking at the previous generation

	//Make a blank board
	for y := 0; y < len(g.state); y++ {
		g.state[y] = make([]Cell, cols+2)
		g.prevState[y] = make([]Cell, cols+2)
	}
	return &g
}

//RandSeed will change a bunch of cells in the middle of board to random on and off states. Requires the game to be initialized and of minimum size 4x4
func (g *Game) RandSeed() {
	//Find the center (approximate for odd height or width) Cell of the board
	xMid := g.cols / 2
	yMid := g.rows / 2
	
	//TEMP placeholder for actual random number generator
	rand := []int {0,1,0,1,1,1,1,0,0,0,1,0,1,0,1,1,}
	
	//Iterate over a 4x4 square around the center Cell
	i := 0
	for y := yMid - 1; y < yMid + 3; y++ {
		for x := xMid - 1; x < xMid +3; x++ {
			if rand[i] == 1 {
				g.state[y][x].Alive = !g.state[y][x].Alive
			}
			i++
		}
	}
	return
}

//Cell holds the state of one cell
type Cell struct {
	Alive  bool
	Player uint
	Type   uint
}
