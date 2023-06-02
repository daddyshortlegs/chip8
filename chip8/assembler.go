package chip8

type Assembler struct {
	code []byte
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.code = make([]byte, 0)
	return a
}

func (a *Assembler) SetRegister(index byte, value byte) {
	register := 0x60 + index
	opcodes := []byte{register, value}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) SetIndexRegister(value uint16) {
	first := extractFirstByte(value)
	second := extractSecondByte(value)
	instruction := 0xA0 | first
	opcodes := []byte{instruction, second}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) Display(xRegister byte, yRegister byte, n byte) {
	instruction := 0xD0 | xRegister
	second := (yRegister << 4) | n
	opcodes := []byte{instruction, second}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) GetKey(xRegister byte) {
	instruction := 0xF0 | xRegister
	opcodes := []byte{instruction, 0x0A}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) Assemble() []byte {
	return a.code
}
