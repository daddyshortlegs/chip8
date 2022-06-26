package chip8

type Chip8vm struct {
	memory                  [4096]byte
	registers               [16]byte
	indexRegister           uint16
	pc                      uint16
	d                       Display
	previousInstructionJump bool
}

type Display interface {
	ClearScreen()
	DrawPattern()
}

func (v *Chip8vm) SetDisplay(d Display) {
	v.d = d
}

func (v *Chip8vm) Load(bytes []byte) {
	copy(v.memory[:], bytes)
}

func (v *Chip8vm) Run() {
	v.previousInstructionJump = false
	for {
		var instr instruction
		instr = v.fetchNextInstruction(instr)

		if instr.first == 0x00 && instr.second == 0x00 {
			break
		}

		firstNibble := v.getNibble(instr)

		if instr.first == 0x00 && instr.second == 0xE0 {
			v.d.ClearScreen()
		} else if firstNibble == 0x10 {
			v.jump(instr.first, instr.second)
			v.previousInstructionJump = true
			continue
		} else if firstNibble == 0x60 {
			v.setRegister(instr.first, instr.second)
		} else if firstNibble == 0x70 {
			v.addToRegister(instr.first, instr.second)
		} else if firstNibble == 0xA0 {
			v.setIndexRegister(instr.first, instr.second)
		}
		v.previousInstructionJump = false
	}
}

func (v *Chip8vm) fetchNextInstruction(instr instruction) instruction {
	if v.previousInstructionJump == false {
		instr = v.fetchAndIncrement()
	} else {
		instr = v.fetch()
	}
	return instr
}

func (v *Chip8vm) getNibble(instr instruction) byte {
	mask := byte(0b11110000)
	firstNibble := instr.first & mask
	return firstNibble
}

func (v *Chip8vm) fetch() instruction {
	i := instruction{v.memory[v.pc], v.memory[v.pc+1]}
	return i
}

func (v *Chip8vm) fetchAndIncrement() instruction {
	i := instruction{v.memory[v.pc], v.memory[v.pc+1]}
	v.pc += 2
	return i
}

func (v *Chip8vm) setRegister(firstByte byte, secondByte byte) {
	v.registers[extractNibble(firstByte)] = secondByte
}

func (v *Chip8vm) addToRegister(firstByte byte, secondByte byte) {
	v.registers[extractNibble(firstByte)] += secondByte
}

func (v *Chip8vm) setIndexRegister(firstByte byte, secondByte byte) {
	v.indexRegister = extract12BitNumber(firstByte, secondByte)
}

func (v *Chip8vm) jump(firstByte byte, secondByte byte) {
	v.pc = extract12BitNumber(firstByte, secondByte)
}
