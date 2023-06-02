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
	a.SetRegister(0, 0x19)
	a.SetRegister(1, 0x00)

	a.SetRegister(2, 0x20)
	a.SetRegister(3, 0x00)

	a.SetRegister(4, 0x1E)
	a.SetRegister(5, 0x00)

	a.SetRegister(6, 0x23)
	a.SetRegister(7, 0x00)

	a.SetRegister(8, 0x28)
	a.SetRegister(9, 0x00)

	a.SetIndexRegister(0x50)
	a.Display(0, 1, 5)

	a.SetIndexRegister(0x55)
	a.Display(2, 3, 5)

	a.SetIndexRegister(0x5A)
	a.Display(4, 5, 5)

	a.SetIndexRegister(0x5F)
	a.Display(6, 7, 5)

	a.SetIndexRegister(0x64)
	a.Display(8, 9, 5)

	a.GetKey(3)

	return a.Assemble()
}
