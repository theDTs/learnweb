package main

import (
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	doc := js.Global.Get("document")
	//Clear the page
	doc.Get("body").Set("innerHTML", "")
	doc.Set("title", "Canvas Tests")

	var d display

	d.rows, d.cols = 50, 50
	d.cellSize = 10
	d.height = d.rows * d.cellSize
	d.width = d.cols * d.cellSize

	d.Game = NewGame(uint(d.rows), uint(d.cols))
	d.Game.SeedAcorn()
	d.Board = d.Game.GetBoard()

	d.canvas = doc.Call("createElement", "canvas")
	d.canvas.Set("height", d.height)
	d.canvas.Set("width", d.width)
	d.canvas.Set("id", "myCanvas")
	doc.Get("body").Call("appendChild", d.canvas)

	d.ctx = d.canvas.Call("getContext", "2d")
	d.Draw()
	js.Global.Call("setInterval", func(){d.DrawNext()}, 500)

}

type display struct {
	rows, cols, height, width, cellSize int
	Board
	*Game
	canvas *js.Object
	ctx    *js.Object
}

func (d *display) Draw() {
	//Clear canvas
	d.ctx.Set("fillStyle", "white")
	d.ctx.Call("fillRect", 0, 0, d.width, d.height)
	d.ctx.Set("fillStyle", "red")

	//Update copy of Game's state.
	d.Game.Update(d.Board)

	//Draw living cells
	rows, cols := d.Game.Rows(), d.Game.Cols()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if d.Board[y][x].Alive {
				d.ctx.Call("fillRect",
					x*d.cellSize, // fillRect uses "x, y" order
					y*d.cellSize,
					d.cellSize,
					d.cellSize)
			}
		}
	}

	return
}

func (d *display) DrawNext() {
	d.Game.Step()
	d.Draw()
	return
}

//Game describes a Game of Life "board", with a grid of cells that can be updated turn by turn
type Game struct {
	state, prevState, board Board
	rows, cols              uint
}

//Board represents a grid of cellular automata. Note that coordinates are reversed when indexing (y, x) to match the order "rows, columns"
type Board [][]Cell

//Cell holds the state of one cell
type Cell struct {
	Alive  bool
	Player uint
	Type   uint
}

//NewGame returns a game with an empty board
func NewGame(rows, cols uint) *Game {
	var g Game
	g.rows, g.cols = rows, cols

	g.state = makeBoard(rows+2, cols+2)     //+2 adds a border 1 Cell wide at each edge of the board, for the purposes of counting around the edge cells
	g.prevState = makeBoard(rows+2, cols+2) //We'll get the state of each generation by looking at the previous generation
	g.board = makeBoard(rows, cols)

	return &g
}

//Rows gets the number of rows
func (g *Game) Rows() int {
	return int(g.rows)
}

//Cols gets the number of columns
func (g *Game) Cols() int {
	return int(g.cols)
}

var offset = [8]struct{ y, x int }{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 1},
	{0, 1},
	{1, 1},
	{-1, 0},
	{1, 0},
}

//Step calculates which are alive and dead to advance the simulation one generation.
func (g *Game) Step() {

	rows, cols := int(g.rows), int(g.cols)
	var neighbors int

	//Swapping the boards allows us to overwrite the state from two turns ago, while reading the state from last turn
	g.state, g.prevState = g.prevState, g.state

	//Iterate over every row, skipping the border rows
	for y := 1; y <= rows; y++ {
		//Cell by cell, skipping borders
		for x := 1; x <= cols; x++ {

			// Count the neighbors
			neighbors = 0
			for i := 0; i < 8; i++ {
				if g.prevState[offset[i].y+y][offset[i].x+x].Alive {
					neighbors++
				}
			}

			/*Implement game of life ruleset:
			  With exactly 2 neighbors: Alive cell stays alive and dead cell stays dead
			  With exactly 3 neighbors: Alive cell stays alive and dead cell also lives
			  Any other number of neighbors: Alive cell dies and dead cell stays dead
			*/
			if neighbors == 2 {
				g.state[y][x] = g.prevState[y][x]
			} else if neighbors == 3 {
				g.state[y][x].Alive = true
			} else {
				g.state[y][x].Alive = false
			}

		}
	}

	//Update the copy
	for y := 0; y < int(g.rows); y++ {
		for x := 0; x < int(g.cols); x++ {
			//+1 to avoid copying the border cells
			g.board[y][x] = g.state[y+1][x+1]
		}
	}

	return

}

func makeBoard(rows, cols uint) Board {

	//Create anonymous backing Array for b
	back := make([]Cell, rows*cols)

	b := make(Board, rows)

	for y := 0; y < int(rows); y++ {
		//Map each row onto a subslice of the backing Array the length of cols
		b[y] = back[y*int(cols) : (y+1)*int(cols)]
	}

	return b
}

//SeedRand will change a bunch of cells in the middle of the game board to random on and off states. Requires the game to be initialized and of minimum size 4x4
func (g *Game) SeedRand() {

	//Find the center (approximate for odd height or width) Cell of the board
	xMid := g.cols / 2
	yMid := g.rows / 2

	//TEMP placeholder for actual random number generator
	rand := []int{0, 1, 0, 1, 1, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1}

	//Iterate over a 4x4 square around the center Cell
	i := 0
	for y := yMid - 1; y < yMid+3; y++ {
		for x := xMid - 1; x < xMid+3; x++ {
			if rand[i] == 1 {
				g.state[y][x].Alive = !g.state[y][x].Alive
			}
			i++
		}
	}

	//Update the copy
	for y := 0; y < int(g.rows); y++ {
		for x := 0; x < int(g.cols); x++ {
			g.board[y][x] = g.state[y+1][x+1]
		}
	}

	return
}

//SeedAcorn will clear the game board and place the "acorn", a long-lived "methuselah" pattern, in the center. Requires an initialized game with a board of at least 3 x 7
func (g *Game) SeedAcorn() {
	//Acorn pattern with rows end-to-end
	acorn := []int{0, 1, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		1, 1, 0, 0, 1, 1, 1,
	}

	//For performance in gopherjs
	rows, cols := int(g.rows), int(g.cols)

	//Clear the board
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.state[y][x].Alive = false
		}
	}

	//Center the pattern on the board
	startRow := (rows - 1) / 2
	startCol := (cols - 1) / 2

	//Copy acorn pattern onto center of board
	for y, i := 0, 0; y < 3; y++ {
		for x := 0; x < 7; x++ {
			if acorn[i] == 1 {
				g.state[startRow+y][startCol+x].Alive = true
			}
			i++
		}
	}

	//Update the copy
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = g.state[y+1][x+1]
		}
	}

	return
}

//GetBoard returns a grid of the Game's current board for copying or displaying; the grid can also be updated by providing it to Game.Update() (avoiding excessive allocations)
func (g *Game) GetBoard() Board {

	//Initialize blank board of Cells
	b := makeBoard(g.rows, g.cols)

	//Copy board
	for y := 0; y < int(g.rows); y++ {
		for x := 0; x < int(g.cols); x++ {
			b[y][x] = g.board[y][x]
		}
	}

	return b
}

//Update takes a Board and copies the game's current state to it.
func (g *Game) Update(b Board) {

	//Copy board
	for y := 0; y < int(g.rows); y++ {
		for x := 0; x < int(g.cols); x++ {
			b[y][x] = g.board[y][x]
		}
	}

	return
}
