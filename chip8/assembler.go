package chip8

type Assembler struct {
	code []byte
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.code = make([]byte, 0)
	return a
}

func (a *Assembler) ClearScreen() {
	a.buildArray(0x00, 0xE0)
}

func (a *Assembler) Jump(address uint16) {
	instruction := 0x10 | extractFirstByte(address)
	a.buildArray(instruction, extractSecondByte(address))
}

func (a *Assembler) Sub(address uint16) {
	instruction := 0x20 | extractFirstByte(address)
	a.buildArray(instruction, extractSecondByte(address))
}

func (a *Assembler) Return() {
	a.buildArray(0x00, 0xEE)
}

func (a *Assembler) SkipIfEqual(xRegister byte, value byte) {
	instruction := 0x30 | xRegister
	a.buildArray(instruction, value)
}

func (a *Assembler) SkipIfNotEqual(xRegister byte, value byte) {
	instruction := 0x40 | xRegister
	a.buildArray(instruction, value)
}

func (a *Assembler) SkipIfRegistersEqual(xRegister byte, yRegister byte) {
	instruction := 0x50 | xRegister
	second := (yRegister << 4) | 0
	a.buildArray(instruction, second)
}

func (a *Assembler) SkipIfRegistersNotEqual(xRegister byte, yRegister byte) {
	instruction := 0x90 | xRegister
	second := (yRegister << 4) | 0
	a.buildArray(instruction, second)
}

func (a *Assembler) SetRegister(index byte, value byte) {
	a.buildArray(0x60+index, value)
}

func (a *Assembler) AddToRegister(index byte, value byte) {
	a.buildArray(0x70+index, value)
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

func (a *Assembler) buildArray(key byte, value byte) {
	opcodes := []byte{key, value}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) Assemble() []byte {
	return a.code
}
