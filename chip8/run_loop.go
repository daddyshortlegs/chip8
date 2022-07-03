package chip8

type Chip8vm struct {
	memory                  [4096]byte
	registers               [16]byte
	indexRegister           uint16
	pc                      uint16
	d                       Display
	previousInstructionJump bool
	xCoord                  byte
	yCoord                  byte
}

type Display interface {
	ClearScreen()
	DrawPattern()
}

func (v *Chip8vm) Init() {
	font := createFont()
	v.pc = 0x200
	copy(v.memory[0x50:], font)
}

func (v *Chip8vm) SetDisplay(d Display) {
	v.d = d
}

func (v *Chip8vm) Load(bytes []byte) {
	copy(v.memory[0x200:], bytes)
}

func (v *Chip8vm) Run() {
	v.previousInstructionJump = false
	for {
		instr := v.fetchNextInstruction()

		if instr == 0x0000 {
			break
		}

		firstNibble := extractNibble(instr)

		if instr == 0x00E0 {
			v.d.ClearScreen()
		} else if firstNibble == 0x10 {
			println("jumping")
			v.jump(instr)
			v.previousInstructionJump = true
			continue
		} else if firstNibble == 0x60 {
			v.setRegister(instr)
		} else if firstNibble == 0x70 {
			v.addToRegister(instr)
		} else if firstNibble == 0xA0 {
			v.setIndexRegister(instr)
		} else if firstNibble == 0xD0 {
			firstByte := extractFirstByte(instr)
			xRegister := v.getRightNibble(firstByte)
			secondByte := extractSecondByte(instr)
			yRegister := v.getLeftNibble(secondByte)

			v.xCoord = v.registers[xRegister] & 63
			v.yCoord = v.registers[yRegister] & 31

		}
		v.previousInstructionJump = false
	}
}

func (v *Chip8vm) fetchNextInstruction() uint16 {
	if v.previousInstructionJump == false {
		return v.fetchAndIncrement()
	}
	return v.fetch()
}

func (v *Chip8vm) fetch() uint16 {
	i := bytesToWord(v.memory[v.pc], v.memory[v.pc+1])
	return i
}

func (v *Chip8vm) fetchAndIncrement() uint16 {
	i := bytesToWord(v.memory[v.pc], v.memory[v.pc+1])
	v.pc += 2
	return i
}

func (v *Chip8vm) setRegister(instr uint16) {
	nibble := v.getRegisterIndex(instr)
	v.registers[nibble] = extractSecondByte(instr)
}

func (v *Chip8vm) addToRegister(instr uint16) {
	nibble := v.getRegisterIndex(instr)
	v.registers[nibble] += extractSecondByte(instr)
}

func (v *Chip8vm) getRegisterIndex(instr uint16) byte {
	firstByte := extractFirstByte(instr)
	return v.getRightNibble(firstByte)
}

func (v *Chip8vm) setIndexRegister(instr uint16) {
	v.indexRegister = extract12BitNumber(instr)
}

func (v *Chip8vm) jump(address uint16) {
	v.pc = extract12BitNumber(address)
}

func (v *Chip8vm) getXCoordinate() byte {
	return v.xCoord
}

func (v *Chip8vm) getYCoordinate() byte {
	return v.yCoord
}
