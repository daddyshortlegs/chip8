package chip8

type Assembler struct {
	code []byte
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.code = make([]byte, 0)
	return a
}

func (a *Assembler) buildArray(key byte, value byte) {
	opcodes := []byte{key, value}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) SetRegister(index byte, value byte) {
	a.buildArray(0x60+index, value)
}

func (a *Assembler) SetIndexRegister(value uint16) {
	instruction := 0xA0 | extractFirstByte(value)
	a.buildArray(instruction, extractSecondByte(value))
}

func (a *Assembler) Display(xRegister byte, yRegister byte, n byte) {
	second := (yRegister << 4) | n
	a.buildArray(0xD0|xRegister, second)
}

func (a *Assembler) GetKey(xRegister byte) {
	a.buildArray(0xF0|xRegister, 0x0A)
}

func (a *Assembler) Assemble() []byte {
	return a.code
}
