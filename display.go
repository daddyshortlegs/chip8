package main

import "github.com/veandco/go-sdl2/sdl"

// display is 64x32
func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface := getSurface(err, window)
	clearScreen(surface)

	drawPoint(surface, 0, 0)
	drawPoint(surface, 1, 1)
	drawPoint(surface, 2, 2)

	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func getSurface(err error, window *sdl.Window) *sdl.Surface {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	return surface
}

func clearScreen(surface *sdl.Surface) error {
	return surface.FillRect(nil, 0)
}

func drawPoint(surface *sdl.Surface, x int32, y int32) {
	rect := sdl.Rect{x * 10, y * 10, 9, 9}
	surface.FillRect(&rect, 0x00fffff0)
}
