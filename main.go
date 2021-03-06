package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"github.com/yuin/gopher-lua"

	"math/rand"
	"image/color"
	"time"
	"fmt"
	"os"
	"sync"
)
var (
	CONFIG *Config
	list []int
	changed []bool
	stop = make(chan byte, 1)
	running = true
	wg = sync.WaitGroup{}
	FILENAME string
	finished = false
	showCount = 0
)

func run() {
	width := float64(CONFIG.LIST_LENGTH * CONFIG.BLOCK_WIDTH)
	height := float64(CONFIG.LIST_LENGTH) * CONFIG.BLOCK_HEIGHT_MULT

	if CONFIG.SHOWER == "circle" || CONFIG.SHOWER == "shell" || CONFIG.SHOWER == "hoops" {
		height = 2 * height
		width = height
	}
	cfg := pixelgl.WindowConfig{
		Title:  "VSF",
		Bounds: pixel.R(0, 0, width, height),
	}
	if CONFIG.VSYNC {
		cfg.VSync = true
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	
	//filename
	atlas := text.NewAtlas(
		basicfont.Face7x13,
		[]rune(FILENAME),
		text.ASCII)
	txt := text.New(pixel.V(0, height - 13), atlas)
	fmt.Fprintf(txt, "VSF: %s", FILENAME)

	last := time.Now()
	filter := make([]time.Duration, 0, CONFIG.FPSFILTER)
	var fps float64 = 0
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeySpace) {
			if !finished {
				if running {
					wg.Add(1)
				} else {
					wg.Done()
				}
				running = !running
			} else {
				list = rand.Perm(CONFIG.LIST_LENGTH)
				for i, _ := range list {
					list[i]++
				}
				showCount = 0
				go runLua()
				finished = false
			}
		}
		if win.JustPressed(pixelgl.KeyQ) {
			win.SetClosed(true)
			os.Exit(0)
		}
		win.Clear(color.RGBA{CONFIG.BG[0], CONFIG.BG[1], CONFIG.BG[2], CONFIG.BG[3]})
		txt.Draw(win, pixel.IM)
		fpsText := text.New(pixel.V(0, height - 26), atlas)
		fmt.Fprintf(fpsText, "FPS: %.2f", fps)
		fpsText.Draw(win, pixel.IM)
		showCountText := text.New(pixel.V(0, height - 39), atlas)
		fmt.Fprintf(showCountText, "# of times show is called: %d", showCount)
		showCountText.Draw(win, pixel.IM)
		if finished {
			finishedText := text.New(pixel.V(0, height - 52), atlas)
			fmt.Fprint(finishedText, "Sort finished! Press <space> to start again")
			finishedText.Draw(win, pixel.IM)
		}
		switch CONFIG.SHOWER{
		case "rect":
			rectDraw(list, changed).Draw(win)
		case "point":
			pointDraw(list, changed).Draw(win)
		case "circle":
			circleDraw(list, changed).Draw(win)
		case "block":
			blockDraw(list, changed).Draw(win)
		case "shell":
			shellDraw(list, changed).Draw(win)
		case "hoops":
			hoopDraw(list, changed).Draw(win)
		default:
			panic("Invalid shower given!")
		}
		win.Update()

		frameDiff := time.Since(last)
		last = time.Now()
		if len(filter) == cap(filter) {
			for i := 1; i < len(filter); i++ {
				filter[i-1] = filter[i]
			}
			filter[len(filter)-1] = frameDiff
		} else {
			filter = append(filter, frameDiff)
		}
		var total int64
		for _, t := range filter {
			total += int64(t)
		}
		fps = 1 / time.Duration(total / int64(len(filter))).Seconds()
	}
	stop<-1
}

func show(L *lua.LState) int {
	showCount++
	wg.Wait()
	time.Sleep(time.Duration(CONFIG.SLEEP * float64(time.Millisecond)))
	newTable := L.ToTable(1)
	if newTable.Len() != CONFIG.LIST_LENGTH {
		//TODO: better error handling
		panic("List of improper length given")
	}
	newList := make([]int, CONFIG.LIST_LENGTH)
	//exists := make(map[int]bool, CONFIG.LIST_LENGTH)
	newTable.ForEach(func(a, b lua.LValue) {
		if b.(lua.LNumber) <= 0 || int(b.(lua.LNumber)) > CONFIG.LIST_LENGTH {
			panic("Invalid value found in list")
		}
		/*if exists[int(b.(lua.LNumber))] {
			panic("Duplicate value found in list")
		} else {
			exists[int(b.(lua.LNumber))] = true
		}*/
		newList[int(a.(lua.LNumber))-1] = int(b.(lua.LNumber))
		if int(b.(lua.LNumber)) != list[int(a.(lua.LNumber))-1] {
			changed[int(a.(lua.LNumber))-1] = true
		} else {
			changed[int(a.(lua.LNumber))-1] = false
		}
	})
	list = newList
	return 0
}

func main() {
	if len(os.Args) >= 2 {
		FILENAME = os.Args[1]
	} else {
		fmt.Println("Usage: vsf <lua file> [settings file]\nlua file must contain a function called sort that sorts an array and use the show function to display that array on the screen\npress space to pause")
		return
	}
	configFile := ""
	if len(os.Args) >= 3 {
		configFile = os.Args[2]
	}
	var err error
	CONFIG, err = parse(configFile)
	if err != nil {
		panic(err)
	}
	changed = make([]bool, CONFIG.LIST_LENGTH)

	rand.Seed(time.Now().UnixNano())
	list = rand.Perm(CONFIG.LIST_LENGTH)
	for i, _ := range list {
		//list[i] = len(list) - i
		list[i]++
	}
	fmt.Println(list)

	go runLua()
	pixelgl.Run(run)
	<-stop
}

func runLua() {
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("show", L.NewFunction(show))
	if err := L.DoFile(FILENAME); err != nil {
		panic(err)
	}
	tableList := lua.LTable{}
	for i, val := range list {
		tableList.Insert(i+1, lua.LNumber(val))
	}
	err := L.CallByParam(lua.P{
		Fn: L.GetGlobal("sort"),
		NRet: 1,
		Protect: true,
	}, &tableList)
	if err != nil {
		panic(err)
	}
	finished = true
}
