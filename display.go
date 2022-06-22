package main

import "github.com/veandco/go-sdl2/sdl"

type display struct {
	window *sdl.Window
}

func (d display) startUp() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	d.window = window
	defer d.window.Destroy()

	d.clearScreen()
	d.drawPattern()

	d.waitForExit()
}

func (d display) drawPattern() {
	surface := d.getSurface()

	d.drawPoint(surface, 0, 0)
	d.drawPoint(surface, 1, 1)
	d.drawPoint(surface, 2, 2)

	d.window.UpdateSurface()
}

func (d display) clearScreen() {
	surface := d.getSurface()
	surface.FillRect(nil, 0)
	d.window.UpdateSurface()
}

func (d display) waitForExit() {
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

func (d display) drawPoint(surface *sdl.Surface, x int32, y int32) {
	rect := sdl.Rect{x * 10, y * 10, 9, 9}
	surface.FillRect(&rect, 0x00fffff0)
}

func (d display) getSurface() *sdl.Surface {
	surface, err := d.window.GetSurface()
	if err != nil {
		panic(err)
	}
	return surface
}
