package chip8

type instruction struct {
	first  byte
	second byte
}

type memory struct {
	bytes [4096]byte
	pc    int
}

const (
	ClearScreen = iota
	Jump
	SetRegister
	AddValueToRegister
	SetIndexRegister
	DisplayDraw
)

func (m *memory) load(bytes []byte) {
	copy(m.bytes[:], bytes)
}

func (m memory) fetch() instruction {
	return instruction{m.bytes[0], m.bytes[1]}
}

func (m *memory) decode() int {
	firstByte := m.bytes[0]
	return decodeInstruction(firstByte)
}
