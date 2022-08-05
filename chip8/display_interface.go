package chip8

type DisplayInterface interface {
	ClearScreen()
	DrawSprite(chip8 *VM, address uint16, numberOfBytes byte, x byte, y byte)
	PollEvents() bool
}
