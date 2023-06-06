package main

import (
	"chip8"
	"time"
)

func main() {
	chip8Display := Chip8Display{}
	defer chip8Display.shutdown()
	chip8Display.startUp()

	nano := time.Now().UnixNano()
	random := chip8.PseudoRandom{Seed: nano}

	vm := chip8.NewVM(&chip8Display, random)

	//dat, _ := ioutil.ReadFile("IBM-Logo.ch8")
	//check(err)

	//vm.Load(dat)
	//vm.Load([]byte{0xF3, 0x0A})

	vm.Load(andysProgram())
	vm.Run()

	//Chip8Display.ClearScreen()
}

func andysProgram() []byte {
	a := chip8.NewAssembler()

	a.ClearScreen()
	drawChar(a, 0, 0, 0x50)
	drawChar(a, 5, 0, 0x55)
	drawChar(a, 10, 0, 0x5A)
	drawChar(a, 15, 0, 0x5F)
	drawChar(a, 20, 0, 0x64)

	a.GetKey(3)
	a.ClearScreen()

	a.GetKey(4) // This will not wait

	return a.Assemble()
}

func drawChar(a *chip8.Assembler, xPosition int, yPosition int, startAddress int) {
	drawAtFromAddress(a, xPosition, yPosition, startAddress, 5)
}

func drawAtFromAddress(a *chip8.Assembler, xPosition int, yPosition int, startAddress int, pixelsHigh int) {
	a.SetRegister(0, byte(xPosition))
	a.SetRegister(1, byte(yPosition))

	a.SetIndexRegister(uint16(startAddress))
	a.Display(0, 1, byte(pixelsHigh))
}
