package main

import (
	"chip8"
)

func main() {
	theDisplay := display{}
	theDisplay.startUp()

	vm := chip8.Chip8vm{}
	instruction := []byte{0x12, 0x20}
	vm.Load(instruction)
}
