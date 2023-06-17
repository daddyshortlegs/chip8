package chip8

type instruction struct {
	instr uint16
}

func NewInstruction(instr uint16) *instruction {
	i := new(instruction)
	i.extractNibbles(instr)
	return i
}

func (i *instruction) extractNibbles(instr uint16) (byte, byte, byte, byte) {
	opCode := extractNibble(instr)
	vx := getRightNibble(extractFirstByte(instr))
	secondByte := extractSecondByte(instr)
	vy := getLeftNibble(secondByte)
	opcode2 := getRightNibble(secondByte)
	return opCode, vx, vy, opcode2
}
