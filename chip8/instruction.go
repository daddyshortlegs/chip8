package chip8

const (
	ClearScreen = iota
	Jump
	SetRegister
	AddValueToRegister
	SetIndexRegister
	DisplayDraw
)

func decodeInstruction(firstByte byte, secondByte byte) int {
	if firstByte == 0x00 && secondByte == 0xE0 {
		return ClearScreen
	}

	mask := byte(0b11110000)
	firstNibble := firstByte & mask

	instructions := map[byte]int{
		0x10: Jump,
		0x60: SetRegister,
		0x70: AddValueToRegister,
		0xA0: SetIndexRegister,
		0xD0: DisplayDraw,
	}

	return instructions[firstNibble]
}
