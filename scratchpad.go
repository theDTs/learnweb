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

	ctx := doc.Call("getElementById", "myCanvas").Call("getContext", "2d")
	ctx.Set("fillStyle", "red")

	ctx.Call("fillRect", 0, 0, 50, 50)

}
