package chip8

type instruction struct {
	instr uint16
}

func (i instruction) extractNibbles() (byte, byte, byte, byte) {
	opCode := extractNibble(i.instr)
	vx := getRightNibble(extractFirstByte(i.instr))
	secondByte := extractSecondByte(i.instr)
	vy := getLeftNibble(secondByte)
	opcode2 := getRightNibble(secondByte)
	return opCode, vx, vy, opcode2
}
