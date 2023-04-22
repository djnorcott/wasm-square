package main

import (
	"math/rand"
	"strconv"
	"syscall/js"
)

type Square struct {
	X, Y, Size, Speed, Direction float64
	Color                        string
}

func logToConsole(msg string) {
	js.Global().Get("console").Call("log", msg)
}

func randomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return "rgb(" + strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ")"
}

func main() {
	window := js.Global()
	document := window.Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	width := float64(window.Get("innerWidth").Int())
	height := float64(window.Get("innerHeight").Int())
	canvas.Set("width", width)
	canvas.Set("height", height)

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
	window.Set("resizeCanvas", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		newWidth := float64(args[0].Int())
		newHeight := float64(args[1].Int())
		resizeCanvas(newWidth, newHeight)
		return nil
	}))

	addSquare := func(x, y float64) {
		size := rand.Float64()*50 + 50
		color := randomColor()
		speed := 14.0 - (size-50)/50*7.0 // Smaller squares will have a faster speed

		// make a variable called direction that is either 1.0 or -1.0
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

	initialX := (width - 100) / 2
	initialY := (height - 100) / 2
	addSquare(initialX, initialY)

	window.Set("addSquare", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		x := float64(args[0].Int())
		y := float64(args[1].Int())
		addSquare(x, y)
		return nil
	}))

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

	drawSquare := func(s Square) {
		ctx.Set("fillStyle", s.Color)
		ctx.Call("fillRect", s.X, s.Y, s.Size, s.Size)
	}

	drawSquares := func() {
		ctx.Set("fillStyle", "#f0f0f0")
		ctx.Call("fillRect", 0, 0, width, height)
		for _, s := range squares {
			drawSquare(s)
		}
	}

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updatePosition()
		drawSquares()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}
