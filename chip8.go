package main

import (
	"chip8"
	"io/ioutil"
	"time"
)

func main() {
	chip8Display := Chip8Display{}
	defer chip8Display.shutdown()
	chip8Display.startUp()

	nano := time.Now().UnixNano()
	random := chip8.PseudoRandom{Seed: nano}
	vm := chip8.NewVM(chip8Display, random)

	dat, _ := ioutil.ReadFile("test_opcode.ch8")
	//check(err)

	vm.Load(dat)
	//vm.Load([]byte{
	//	0x60, 0x19, // Set register 0 to 0x00
	//	0x61, 0x00, // Set register 1 to 0x00
	//
	//0x62, 0x20, // Set register 2 to 0x05
	//0x63, 0x00, // Set register 3 to 0x00
	//
	//0x64, 0x1E, // Set register 4 to 10
	//0x65, 0x00, // Set register 5 to 0x00
	//
	//0x66, 0x23, // Set register 6 to 15
	//0x67, 0x00, // Set register 7 to 0x00
	//
	//0x68, 0x28, // Set register 8 to 20
	//0x69, 0x00, // Set register 9 to 0x00
	//
	//0xA0, 0x50, // Set Index Register to 0x50
	//0xD0, 0x15, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
	//
	//0xA0, 0x55, // Set Index Register to 0x55
	//0xD2, 0x35, // Draw, Xreg = 5, Y reg = 10, 5 bytes high

	//0xA0, 0x5A, // Set Index Register to 0x55
	//0xD4, 0x55, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
	//0xA0, 0x5F, // Set Index Register to 0x5F
	//0xD6, 0x75, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
	//0xA0, 0x64, // Set Index Register to 0x5F
	//0xD8, 0x95, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
	//})
	vm.Run()

	//Chip8Display.ClearScreen()
}
