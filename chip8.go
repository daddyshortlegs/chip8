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
	m             memory
	registers     [16]byte
	indexRegister uint16
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

		theInstruction := decodeInstruction(firstByte)

		if firstByte == 0x00 {
			break
		}

		switch theInstruction {
		case SetRegister:
			v.setRegister(firstByte, secondByte)
		case AddValueToRegister:
			v.addToRegister(firstByte, secondByte)
		case SetIndexRegister:
			v.setIndexRegister(firstByte, secondByte)
		}
	}
}

func (v *chip8vm) setRegister(firstByte byte, secondByte byte) {
	mask := byte(0b00001111)
	index := firstByte & mask
	v.registers[index] = secondByte
}

func (v *chip8vm) addToRegister(firstByte byte, secondByte byte) {
	mask := byte(0b00001111)
	index := firstByte & mask
	v.registers[index] += secondByte
}

func (v *chip8vm) setIndexRegister(firstByte byte, secondByte byte) {
	shiftedBytes := uint16(firstByte) << 8
	result := shiftedBytes | uint16(secondByte)
	mask := uint16(0b0000111111111111)
	v.indexRegister = result & mask
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
