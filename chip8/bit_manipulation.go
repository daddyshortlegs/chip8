package chip8

func extractNibble(address uint16) byte {
	mask := uint16(0b1111000000000000)
	u := (address & mask) >> 8
	return byte(u)
}

func (v *Chip8vm) getLeftNibble(instr byte) byte {
	mask := byte(0b11110000)
	firstNibble := instr & mask
	return firstNibble >> 4
}

func (v *Chip8vm) getRightNibble(instr byte) byte {
	mask := byte(0b00001111)
	return instr & mask
}

func extractFirstByte(address uint16) byte {
	mask := uint16(0b1111111100000000)
	u := address & mask >> 8
	return byte(u)
}

func extractSecondByte(address uint16) byte {
	mask := uint16(0b0000000011111111)
	return byte(address & mask)
}

func bytesToWord(firstByte byte, secondByte byte) uint16 {
	shiftedBytes := uint16(firstByte) << 8
	return shiftedBytes | uint16(secondByte)
}

func extract12BitNumber(address uint16) uint16 {
	mask := uint16(0b0000111111111111)
	return address & mask
}

func GetValueAtPosition(position int, value byte) byte {
	result := value >> position
	bitmask := byte(0b00000001)
	return result & bitmask
}
