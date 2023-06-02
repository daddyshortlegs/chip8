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

func (a *Assembler) assemble() []byte {
	return a.code
}
