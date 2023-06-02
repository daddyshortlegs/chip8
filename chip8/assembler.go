package chip8

type Assembler struct {
	code []byte
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.code = make([]byte, 0)
	return a
}

func (a *Assembler) setRegister(index byte, value byte) {
	register := 0x60 + index
	opcodes := []byte{register, value}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) setIndexRegister(value uint16) {
	first := extractFirstByte(value)
	second := extractSecondByte(value)
	instruction := 0xA0 | first
	opcodes := []byte{instruction, second}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) display(xRegister byte, yRegister byte, n byte) {
	instruction := 0xD0 | xRegister
	second := (yRegister << 4) | n
	opcodes := []byte{instruction, second}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) getKey(xRegister byte) {
	instruction := 0xF0 | xRegister
	opcodes := []byte{instruction, 0x0A}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) assemble() []byte {
	return a.code
}
