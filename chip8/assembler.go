package chip8

type Assembler struct {
	code []byte
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.code = make([]byte, 0)
	return a
}

func (a *Assembler) setRegister(index int, value int) {
	opcodes := []byte{0x60, 0x00}
	a.code = append(a.code, opcodes...)
	//	a.code[0] = 0x60
	//	a.code[1] = 0x00
}

func (a *Assembler) assemble() []byte {
	return a.code
}
