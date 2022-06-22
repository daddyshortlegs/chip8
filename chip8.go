package main

import (
	"chip8"
)

func main() {
	theDisplay := display{}
	theDisplay.startUp()

	vm := chip8.Chip8vm{}
	instruction := []byte{0x00, 0xE0}
	vm.Load(instruction)
}
