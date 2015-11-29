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
	state, prevState, board [][]Cell
	rows, cols       uint
}

//NewGame returns a game with an empty board
func NewGame(rows, cols uint) *Game {
	var g Game
	g.rows, g.cols = rows, cols

	g.state = makeBoard (rows+2, cols+2)     //+2 adds a border 1 Cell wide at each edge of the board, for the purposes of counting around the edge cells
	g.prevState = makeBoard(rows+2, cols+2) //We'll get the state of each generation by looking at the previous generation
	g.board = makeBoard(rows, cols)
	
	return &g
}

func makeBoard (rows, cols uint) [][]Cell {
	
	//Create anonymous backing Array for b 
	back := make([]Cell, rows*cols)
	
	b := make([][]Cell, rows)
	
	for y := 0; y < rows; y++ {
		//Map each row onto a subslice of the backing Array the length of cols
		b[y] := back [y*cols:(y+1)*cols]
	}
	
	return b
}


//RandSeed will change a bunch of cells in the middle of board to random on and off states. Requires the game to be initialized and of minimum size 4x4@
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

//Board returns a grid of the Game's current board for copying or displaying; the grid can also be updated by providing it to Game.Update() (avoiding excessive allocations)
func (*g Game) Board() [][]Cell {
	
	//Initialize blank board of Cells
	b := makeBoard(g.rows, g.cols)
	
	//Copy board
	for y := 0; y < g.rows; y++ {
		x :=0; x < g.cols; x++ {
			b[y][x] = g.board[y][x] 
		}
	}
	
	return b
}

func (*g Game) Update(b [][]Cell) {
	
	//Copy board
	for y := 0; y < g.rows; y++ {
		x :=0; x < g.cols; x++ {
			b[y][x] = g.board[y][x] 
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
