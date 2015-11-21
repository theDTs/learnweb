package main

import (
	"github.com/gopherjs/gopherjs/js"
)

func main() {
    doc := js.Global.Get("document")
    //Clear the page
    doc.Get("body").Set("innerHTML", "")
}