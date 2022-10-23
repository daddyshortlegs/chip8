package main

import (
	"chip8"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type Chip8Display struct {
	window         *sdl.Window
	keyCode        sdl.Keycode
	keyPressed     bool
	virtualDisplay [][]byte
}

func (k *Chip8Display) GetKey() int {
	k.keyPressed = false
	return int(k.keyCode)
}

func (d *Chip8Display) startUp() {
	d.virtualDisplay = d.NewVirtualDisplay()

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

func (d Chip8Display) DrawSprite(chip8 *chip8.VM, startAddress uint16, heightInPixels byte, x byte, y byte) {
	yPos := y
	address := startAddress
	for n := 0; n < int(heightInPixels); n++ {
		value := chip8.Memory[address]
		address++
		d.drawByte(value, x, yPos)
		yPos++
	}
}

func (d Chip8Display) NewVirtualDisplay() [][]byte {
	virtualDisplay := make([][]byte, 32)
	for i := range virtualDisplay {
		virtualDisplay[i] = make([]byte, 64)
	}
	return virtualDisplay
}

func (d Chip8Display) drawByte(value byte, xpos byte, ypos byte) {
	for index := 7; index >= 0; index-- {
		fmt.Printf("Drawing at pos %d, %d\n", xpos, ypos)
		bit := chip8.GetValueAtPosition(index, value)
		if bit == 1 {
			if d.virtualDisplay[ypos][xpos] == 1 {
				d.virtualDisplay[ypos][xpos] = 0
			} else {
				d.virtualDisplay[ypos][xpos] = 1
				// Should set VF to 1
			}
		}
		xpos += 1
	}

	d.writeDisplay()
}

func (d Chip8Display) writeDisplay() {
	surface := d.getSurface()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if d.virtualDisplay[y][x] == 1 {
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
