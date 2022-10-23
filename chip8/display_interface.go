package chip8

type EventType int

const (
	QuitEvent = iota
	KeyboardEvent
	NoEvent
)

type DisplayInterface interface {
	ClearScreen()
	DrawSprite(address uint16, numberOfBytes byte, x byte, y byte, memory [4096]byte)
	PollEvents() EventType
	GetKey() int
}
