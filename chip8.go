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
	Jump = iota
	SetRegister
	AddValueToRegister
	SetIndexRegister
)

func (m *memory) load(bytes []byte) {
	copy(m.bytes[:], bytes)
}

func (m memory) fetch() instruction {
	return instruction{m.bytes[0], m.bytes[1]}
}

func (m *memory) decode() int {
	firstByte := m.bytes[0]
	mask := byte(0b11110000)

	firstNibble := firstByte & mask

	if firstNibble == 0x10 {
		return Jump
	}

	if firstNibble == 0x60 {
		return SetRegister
	}

	if firstNibble == 0x70 {
		return AddValueToRegister
	}

	return SetIndexRegister
}
