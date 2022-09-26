package chip8

type EventType int

const (
	QuitEvent = iota
	KeyboardEvent
	NoEvent
)

type DisplayInterface interface {
	ClearScreen()
	DrawSprite(chip8 *VM, address uint16, numberOfBytes byte, x byte, y byte)
	PollEvents() EventType
	GetKey() int
}
