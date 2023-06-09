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
	a.buildArray(0x30+xRegister, value)
}

func (a *Assembler) SkipIfNotEqual(xRegister byte, value byte) {
	a.buildArray(0x40+xRegister, value)
}

func (a *Assembler) SkipIfRegistersEqual(xRegister byte, yRegister byte) {
	a.buildArray(0x50+xRegister, (yRegister<<4)|0)
}

func (a *Assembler) SkipIfRegistersNotEqual(xRegister byte, yRegister byte) {
	a.buildArray(0x90+xRegister, (yRegister<<4)|0)
}

func (a *Assembler) SetRegister(xRegister byte, value byte) {
	a.buildArray(0x60+xRegister, value)
}

func (a *Assembler) AddToRegister(xRegister byte, value byte) {
	a.buildArray(0x70+xRegister, value)
}

func (a *Assembler) Set(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|0)
}

func (a *Assembler) Or(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|1)
}

func (a *Assembler) And(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|2)
}

func (a *Assembler) Xor(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|3)
}

func (a *Assembler) Add(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|4)
}

func (a *Assembler) Subtract(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|5)
}

func (a *Assembler) SubtractLast(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|7)
}

func (a *Assembler) ShiftRight(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|6)
}

func (a *Assembler) ShiftLeft(xRegister byte, yRegister byte) {
	a.buildArray(0x80+xRegister, (yRegister<<4)|0xE)
}

func (a *Assembler) SetIndexRegister(value uint16) {
	instruction := 0xA0 | extractFirstByte(value)
	a.buildArray(instruction, extractSecondByte(value))
}

func (a *Assembler) SetJumpWithOffset(address uint16) {
	instruction := 0xB0 | extractFirstByte(address)
	a.buildArray(instruction, extractSecondByte(address))
}

func (a *Assembler) Random(xRegister byte, value byte) {
	a.buildArray(0xC0+xRegister, value)
}

func (a *Assembler) Display(xRegister byte, yRegister byte, n byte) {
	a.buildArray(0xD0+xRegister, (yRegister<<4)|n)
}

func (a *Assembler) SkipIfKeyPressed(xRegister byte) {
	a.buildArray(0xE0+xRegister, 0x9E)
}

func (a *Assembler) SkipIfKeyNotPressed(xRegister byte) {
	a.buildArray(0xE0+xRegister, 0xA1)
}

func (a *Assembler) GetDelayTimer(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x07)
}

func (a *Assembler) SetDelayTimer(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x15)
}

func (a *Assembler) SetSoundTimer(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x18)
}

func (a *Assembler) AddToIndex(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x1E)
}

func (a *Assembler) GetKey(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x0A)
}

func (a *Assembler) FontChar(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x29)
}

func (a *Assembler) BCD(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x33)
}

func (a *Assembler) Store(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x55)
}

func (a *Assembler) Load(xRegister byte) {
	a.buildArray(0xF0+xRegister, 0x65)
}

func (a *Assembler) buildArray(key byte, value byte) {
	opcodes := []byte{key, value}
	a.code = append(a.code, opcodes...)
}

func (a *Assembler) Assemble() []byte {
	return a.code
}
