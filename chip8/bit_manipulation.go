package chip8

func extractNibble(firstByte byte) byte {
	mask := byte(0b00001111)
	return firstByte & mask
}

func extract12BitNumber(firstByte byte, secondByte byte) uint16 {
	shiftedBytes := uint16(firstByte) << 8
	result := shiftedBytes | uint16(secondByte)
	mask := uint16(0b0000111111111111)
	return result & mask
}
