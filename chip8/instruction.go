package chip8

type instruction struct {
	instr   uint16
	opCode  byte
	vx      byte
	vy      byte
	opCode2 byte
}

func NewInstruction(instr uint16) *instruction {
	i := new(instruction)
	i.extractNibbles(instr)
	return i
}

func (i *instruction) extractNibbles(instr uint16) {
	i.opCode = extractNibble(instr)
	i.vx = getRightNibble(extractFirstByte(instr))
	secondByte := extractSecondByte(instr)
	i.vy = getLeftNibble(secondByte)
	i.opCode2 = getRightNibble(secondByte)
}
