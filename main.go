package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"

	"math/rand"
	"image/color"

	"github.com/yuin/gopher-lua"
	"time"
	"fmt"
)
const (
	LIST_LENGTH = 500 //the length of the list to be sorted
	BLOCK_WIDTH = 2 //the width of each block
	BLOCK_HEIGHT_MULT = 1 //the amount that the height of the block is multiplied by its position
	WIDTH = LIST_LENGTH * BLOCK_WIDTH //the width of the window
	HEIGHT = LIST_LENGTH * BLOCK_HEIGHT_MULT //the height of the window
	SLEEP = 0 //how many milliseconds to sleep between showings
)
var (
	imd *imdraw.IMDraw
	list = rand.Perm(LIST_LENGTH)
	stop = make(chan byte, 1)
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	
	//background
	imd = imdraw.New(nil)
	imd.Color = pixel.RGB(0, 1, 0)
	imd.Push(pixel.V(0, 0))
	imd.Color = pixel.RGB(0, 1, 1)
	imd.Push(pixel.V(WIDTH, HEIGHT))
	imd.Rectangle(0)

	for !win.Closed() {
		win.Clear(color.RGBA{0x55, 0xaa, 0xaa, 0xff})
		imd.Draw(win)
		for i, val := range list {
			rect := imdraw.New(nil)
			rect.Color = pixel.RGB(1, 1, 1)
			rect.Push(pixel.V(
				float64(i * BLOCK_WIDTH),
				float64(val * BLOCK_HEIGHT_MULT)))
			rect.Push(pixel.V(float64((i+1) * BLOCK_WIDTH), 0))
			rect.Rectangle(0)
			rect.Draw(win)
		}
		win.Update()
	}
	stop<-1
}

func show(L *lua.LState) int {
	time.Sleep(SLEEP * time.Millisecond)
	newTable := L.ToTable(1)
	if newTable.Len() != LIST_LENGTH {
		//TODO: better error handling
		panic("List of improper length given")
	}
	newList := make([]int, LIST_LENGTH)
	exists := make(map[int]bool, LIST_LENGTH)
	newTable.ForEach(func(a, b lua.LValue) {
		fmt.Println(a, b)
		if b.(lua.LNumber) <= 0 || b.(lua.LNumber) > LIST_LENGTH {
			panic("Invalid value found in list")
		}
		if exists[int(b.(lua.LNumber))] {
			panic("Duplicate value found in list")
		} else {
			exists[int(b.(lua.LNumber))] = true
		}
		newList[int(a.(lua.LNumber))-1] = int(b.(lua.LNumber))
	})
	list = newList
	return 0
}

func main() {
	for i, _ := range list {
		list[i]++
	}
	fmt.Println(list)

	go func() {
		L := lua.NewState()
		defer L.Close()
		L.SetGlobal("show", L.NewFunction(show))
		if err := L.DoFile("sort.lua"); err != nil {
			panic(err)
		}
	}()
	pixelgl.Run(run)
	<-stop
}

