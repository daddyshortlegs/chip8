package chip8

func setRegisterOpcode(index byte, value byte) []byte {
	return createOpcode(byte(0x60), index, value)
}

func drawOpcode(xReg byte, yReg byte, bytesHigh byte) []byte {
	b2 := yReg<<4 | bytesHigh
	return []byte{0xD0 | xReg, b2}

}

func setIndexRegisterOpcode(value uint16) []byte {
	mask := uint16(0b0000111111111111)
	twelveBitValue := value & mask

	leftByte := twelveBitValue >> 8
	rightByte := byte(twelveBitValue)

	byte1 := 0xA0 | byte(leftByte)
	return []byte{byte1, rightByte}
}

func createOpcode(instr byte, index byte, value byte) []byte {
	result := instr | index
	return []byte{result, value}
}
