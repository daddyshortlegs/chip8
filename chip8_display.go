package main

import (
	"chip8"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type Chip8Display struct {
	window *sdl.Window
}

func (d *Chip8Display) startUp() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	d.window = window

	//d.drawLetter(0, 0)
}

func (d Chip8Display) shutdown() {
	d.window.Destroy()
	sdl.Quit()
}

func (d Chip8Display) ClearScreen() {
	surface := d.getSurface()
	surface.FillRect(nil, 0)
	d.window.UpdateSurface()
}

func (d Chip8Display) DrawSprite(chip8 *chip8.Chip8VM, startAddress uint16, numberOfBytes byte, x byte, y byte) {
	yPos := y
	address := startAddress
	for n := 0; n < int(numberOfBytes); n++ {
		value := chip8.Memory[address]
		address++
		d.drawByte(value, x, yPos)
		yPos++
	}
}

func (d Chip8Display) drawByte(value byte, xpos byte, ypos byte) {
	surface := d.getSurface()

	fmt.Printf("\n")

	for index := 7; index >= 0; index-- {
		fmt.Printf("Drawing at xpos %d\n", xpos)
		bit := chip8.GetValueAtPosition(index, value)
		if bit == 1 {
			d.drawPoint(surface, xpos, ypos)
		}
		xpos += 1
	}
	d.window.UpdateSurface()
}

func (d Chip8Display) WaitForExit() {
	running := true
	for running {
		running = d.PollEvents()
	}
}

func (d Chip8Display) PollEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return false
		}
	}
	return true
}

func (d Chip8Display) drawPoint(surface *sdl.Surface, x byte, y byte) {
	rect := sdl.Rect{int32(x) * 10, int32(y) * 10, 10, 10}
	surface.FillRect(&rect, 0x00fffff0)
}

func (d Chip8Display) getSurface() *sdl.Surface {
	surface, err := d.window.GetSurface()
	if err != nil {
		panic(err)
	}
	return surface
}