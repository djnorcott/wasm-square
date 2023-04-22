package main

import (
	"math/rand"
	"strconv"
	"syscall/js"
)

// Define a struct type for a square, with X, Y, Size, Speed, Direction and Color fields
type Square struct {
	X, Y, Size, Speed, Direction float64
	Color                        string
}

// Define a function to log a message to the JavaScript console
func logToConsole(msg string) {
	js.Global().Get("console").Call("log", msg)
}

// Define a function to generate a random RGB color
func randomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return "rgb(" + strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ")"
}

func main() {
	// Get the global window object and the document object
	window := js.Global()
	document := window.Get("document")

	// Get the canvas element from the document and its 2D context
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	// Get the initial width and height of the window and set the canvas size accordingly
	width := float64(window.Get("innerWidth").Int())
	height := float64(window.Get("innerHeight").Int())
	canvas.Set("width", width)
	canvas.Set("height", height)

	// Define a function to resize the canvas and all squares
	var squares []Square
	resizeCanvas := func(newWidth, newHeight float64) {
		ratioW := newWidth / width
		ratioH := newHeight / height
		for i, s := range squares {
			squares[i].X = s.X * ratioW
			squares[i].Y = s.Y * ratioH
			squares[i].Size = s.Size * ratioW
		}
		width = newWidth
		height = newHeight
		canvas.Set("width", width)
		canvas.Set("height", height)
	}

	// Set a JavaScript function for resizing the canvas. When the canvas is resized in the browser, the JavaScript 'resizeCanvas' function is called.
	// This function, in turn, calls the Go 'resizeCanvas' function, allowing us to update the canvas size and all the squares on it. We can define a
	// JavaScript function using 'js.FuncOf' that takes in the arguments of the Go function and calls the Go function with those arguments when invoked.
	window.Set("resizeCanvas", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		newWidth := float64(args[0].Int())
		newHeight := float64(args[1].Int())
		resizeCanvas(newWidth, newHeight)
		return nil
	}))

	// Define a function to add a new square at the given position, with a random size and color
	addSquare := func(x, y float64) {
		maxSize := 150.0
		size := rand.Float64()*maxSize + 20
		color := randomColor()
		speed := (maxSize - size + 40) / 8

		direction := 1.0
		if rand.Intn(2) == 0 {
			direction = -1.0
		}

		square := Square{
			X:         x - size/2,
			Y:         y - size/2,
			Size:      size,
			Speed:     speed,
			Direction: direction,
			Color:     color,
		}
		squares = append(squares, square)
	}

	// Add the first square in the center of the canvas
	initialX := (width - 100) / 2
	initialY := (height - 100) / 2
	addSquare(initialX, initialY)

	// Set a JavaScript function for adding a square, calling the Go function
	window.Set("addSquare", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		x := float64(args[0].Int())
		y := float64(args[1].Int())
		addSquare(x, y)
		return nil
	}))

	// Define a function to update the position of all squares and bounce them off the walls
	updatePosition := func() {
		for i, s := range squares {
			newX := s.X + s.Speed*s.Direction
			if newX <= 0 {
				squares[i].X = 0
				squares[i].Direction = 1.0
			} else if newX+s.Size >= width {
				squares[i].X = width - s.Size
				squares[i].Direction = -1.0
			} else {
				squares[i].X = newX
			}
		}
	}

	// Define a function to draw a square on the canvas
	drawSquare := func(s Square) {
		ctx.Set("fillStyle", s.Color)
		ctx.Call("fillRect", s.X, s.Y, s.Size, s.Size)
	}

	// Define a function to draw all squares on the canvas
	drawSquares := func() {
		ctx.Set("fillStyle", "#f0f0f0")
		ctx.Call("fillRect", 0, 0, width, height)
		for _, s := range squares {
			drawSquare(s)
		}
	}

	// Define a JavaScript function to continuously render frames, updating positions and redrawing squares
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updatePosition()
		drawSquares()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	// Start rendering frames
	js.Global().Call("requestAnimationFrame", renderFrame)

	// Block the main goroutine from exiting
	select {}
}
