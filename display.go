package main

import (
	"chip8"
	"github.com/veandco/go-sdl2/sdl"
)

type display struct {
	window *sdl.Window
}

func (d display) DrawPattern(uint16, int) {
	d.drawLetter(0, 0)
}

func (d display) ClearScreen() {
	d.clearScreen()
}

func (d *display) startUp() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	d.window = window

	d.drawLetter(0, 0)
}

func (d display) shutdown() {
	d.window.Destroy()
	sdl.Quit()
}

func (d display) drawLetter(xpos byte, ypos byte) {
	// 0010 0000   0x20
	// 0110 0000   0x60
	// 0010 0000   0x20
	// 0010 0000   0x20
	// 0111 0000   0x70

	//0xF0, 0x90, 0xF0, 0x90, 0x90

	d.drawByte(0x20, 0, 0)
	d.drawByte(0x60, 0, 1)
	d.drawByte(0x20, 0, 2)
	d.drawByte(0x20, 0, 3)
	d.drawByte(0x70, 0, 4)

	d.drawByte(0xF0, 6, 0)
	d.drawByte(0x90, 6, 1)
	d.drawByte(0xF0, 6, 2)
	d.drawByte(0x90, 6, 3)
	d.drawByte(0x90, 6, 4)
}

func (d display) drawByte(value byte, xpos byte, ypos byte) {
	surface := d.getSurface()

	for index := 7; index >= 0; index-- {
		bit7 := chip8.GetValueAtPosition(index, value)
		if bit7 == 1 {
			d.drawPoint(surface, xpos+(7-byte(index)), ypos)
		}
	}
	d.window.UpdateSurface()
}

func (d display) clearScreen() {
	surface := d.getSurface()
	surface.FillRect(nil, 0)
	d.window.UpdateSurface()
}

func (d display) WaitForExit() {
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

func (d display) drawPoint(surface *sdl.Surface, x byte, y byte) {
	rect := sdl.Rect{int32(x * 10), int32(y * 10), 10, 10}
	surface.FillRect(&rect, 0x00fffff0)
}

func (d display) getSurface() *sdl.Surface {
	println("d.window = ", d.window)
	surface, err := d.window.GetSurface()
	if err != nil {
		panic(err)
	}
	return surface
}
