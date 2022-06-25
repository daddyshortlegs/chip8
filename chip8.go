package main

import "chip8"

func main() {
	theDisplay := display{}
	defer theDisplay.shutdown()
	theDisplay.startUp()

	var display chip8.Display
	display = theDisplay

	vm := chip8.Chip8vm{}
	vm.SetDisplay(display)
	display.DrawPattern()
	instruction := []byte{0xD0, 0xD0}
	vm.Load(instruction)
	vm.Run()

	//display.ClearScreen()
	theDisplay.WaitForExit()
}
