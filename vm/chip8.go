package vm

type instruction struct {
	first  byte
	second byte
}

type chip8vm struct {
	memory        [4096]byte
	registers     [16]byte
	indexRegister uint16
	pc            uint16
}

func (v *chip8vm) load(bytes []byte) {
	copy(v.memory[:], bytes)
}

func (v *chip8vm) run() {
	var firstByte byte
	var theInstruction int
	for ok := true; ok; ok = !(firstByte == 0x00) {
		var instr instruction
		if theInstruction != Jump {
			instr = v.fetchAndIncrement()
		} else {
			instr = v.fetch()
		}

		firstByte = instr.first
		secondByte := instr.second

		theInstruction = decodeInstruction(firstByte)

		if firstByte == 0x00 {
			break
		}

		switch theInstruction {
		case Jump:
			v.jump(firstByte, secondByte)
		case SetRegister:
			v.setRegister(firstByte, secondByte)
		case AddValueToRegister:
			v.addToRegister(firstByte, secondByte)
		case SetIndexRegister:
			v.setIndexRegister(firstByte, secondByte)
		}
	}
}

func (v *chip8vm) fetch() instruction {
	i := instruction{v.memory[v.pc], v.memory[v.pc+1]}
	return i
}

func (v *chip8vm) fetchAndIncrement() instruction {
	i := instruction{v.memory[v.pc], v.memory[v.pc+1]}
	v.pc += 2
	return i
}

func (v *chip8vm) setRegister(firstByte byte, secondByte byte) {
	v.registers[extractNibble(firstByte)] = secondByte
}

func (v *chip8vm) addToRegister(firstByte byte, secondByte byte) {
	v.registers[extractNibble(firstByte)] += secondByte
}

func (v *chip8vm) setIndexRegister(firstByte byte, secondByte byte) {
	v.indexRegister = extract12BitNumber(firstByte, secondByte)
}

func (v *chip8vm) jump(firstByte byte, secondByte byte) {
	v.pc = extract12BitNumber(firstByte, secondByte)
}
