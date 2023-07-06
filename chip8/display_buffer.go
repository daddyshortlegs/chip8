package chip8

type DisplayBuffer struct {
	Pixels   [][]byte
	overflow bool
}

func NewDisplayBuffer() *DisplayBuffer {
	db := new(DisplayBuffer)
	db.Pixels = make([][]byte, 32)
	for i := range db.Pixels {
		db.Pixels[i] = make([]byte, 64)
	}
	db.overflow = false
	return db
}

func (d *DisplayBuffer) ClearScreen() {
	for i := range d.Pixels {
		for j := range d.Pixels[i] {
			d.Pixels[i][j] = 0
		}
	}
}

func (d *DisplayBuffer) DrawSprite(startAddress uint16, heightInPixels byte, x byte, y byte, memory [4096]byte) bool {
	yPos := y
	address := startAddress
	for n := 0; n < int(heightInPixels); n++ {
		value := memory[address]
		address++
		d.drawByte(value, x, yPos)
		yPos++
	}

	return d.overflow
}

func (d *DisplayBuffer) drawByte(value byte, xpos byte, ypos byte) {
	for index := 7; index >= 0; index-- {
		//fmt.Printf("Drawing at pos %d, %d\n", xpos, ypos)
		bit := GetValueAtPosition(index, value)
		if bit == 1 {
			if d.Pixels[ypos][xpos] == 1 {
				d.Pixels[ypos][xpos] = 0
				// Should set VF to 1
				d.overflow = true
			} else {
				d.Pixels[ypos][xpos] = 1
			}
		}
		xpos += 1
	}
}

func (d *DisplayBuffer) GetPixelAt(xpos byte, ypos byte) byte {
	return d.Pixels[ypos][xpos]
}
