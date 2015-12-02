package main

import (
  "github.com/gopherjs/gopherjs/js"
)

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

type display struct {
  rows, cols, height, width, cellSize int
  Board
  Game
  canvas  *js.Object
  context *js.Object
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

  g.Update(g.board)

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
      g.board[y][x] = g.state[y][x]
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
