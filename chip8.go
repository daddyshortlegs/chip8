package chip8

type instruction struct {
	first  byte
	second byte
}

type memory struct {
	bytes [4096]byte
	pc    int
}

type chip8vm struct {
	m         memory
	registers [16]byte
}

func (v *chip8vm) load(bytes []byte) {
	v.m.load(bytes)
}

func (v *chip8vm) run() {

	var firstByte byte
	for ok := true; ok; ok = !(firstByte == 0x00) {
		instruction := v.m.fetch()

		firstByte = instruction.first
		secondByte := instruction.second
		decodeInstruction(firstByte)

		if firstByte != 0x00 {
			v.setRegister(firstByte, secondByte)
		}
	}
}

func (v *chip8vm) setRegister(firstByte byte, secondByte byte) {
	mask := byte(0b00001111)
	secondNibble := firstByte & mask
	v.registers[secondNibble] = secondByte
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

func (m *memory) fetch() instruction {
	i := instruction{m.bytes[m.pc], m.bytes[m.pc+1]}
	m.pc += 2
	return i
}
