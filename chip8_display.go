package main

import (
	"chip8"
	"github.com/veandco/go-sdl2/sdl"
)

type Chip8Display struct {
	window        *sdl.Window
	keyCode       sdl.Keycode
	keyPressed    bool
	displayBuffer *chip8.DisplayBuffer
}

func (k *Chip8Display) GetKey() int {
	k.keyPressed = false
	return int(k.keyCode)
}

func (d *Chip8Display) startUp() {
	d.displayBuffer = chip8.NewDisplayBuffer()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	d.window = window
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

func (d Chip8Display) DrawSprite(startAddress uint16, heightInPixels byte, x byte, y byte, memory [4096]byte) {
	d.displayBuffer.DrawSprite(startAddress, heightInPixels, x, y, memory)

	d.writeDisplay()
}

func (d Chip8Display) writeDisplay() {
	surface := d.getSurface()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if d.displayBuffer.GetPixelAt(byte(x), byte(y)) == 1 {
				d.drawPoint(surface, byte(x), byte(y))
			}
		}
	}

	d.window.UpdateSurface()
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

func (d Chip8Display) PollEvents() (quit chip8.EventType) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return chip8.QuitEvent
		case *sdl.KeyboardEvent:
			println("keyboard %s", t.Keysym.Sym)
			d.keyCode = t.Keysym.Sym
			d.keyPressed = true
			return chip8.KeyboardEvent
		}
	}
	return chip8.NoEvent
}
