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
	//drawChar(a, 0, 0, 0x1)
	//drawChar(a, 5, 0, 0x2)
	//drawChar(a, 10, 0, 0x3)
	//drawChar(a, 15, 0, 0x4)
	//drawChar(a, 20, 0, 0x5)
	//drawChar(a, 25, 0, 0xA)
	drawChar(a, 0xB, 30, 0)
	drawChar(a, 0xC, 35, 0)
	drawChar(a, 0xD, 40, 0)

	a.GetKey(3)
	a.ClearScreen()

	a.GetKey(4) // This will not wait

	return a.Assemble()
}

func drawChar(a *chip8.Assembler, char byte, xPosition int, yPosition int) {
	a.SetRegister(0, char)
	a.FontChar(0)

	a.SetRegister(1, byte(xPosition))
	a.SetRegister(2, byte(yPosition))
	a.Display(1, 2, byte(5))
}
