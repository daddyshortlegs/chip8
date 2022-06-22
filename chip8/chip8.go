package chip8

type instruction struct {
	first  byte
	second byte
}

type Chip8vm struct {
	memory        [4096]byte
	registers     [16]byte
	indexRegister uint16
	pc            uint16
	d             Display
}

type Display interface {
	ClearScreen()
}

func (v *Chip8vm) SetDisplay(d Display) {
	v.d = d
}

func (v *Chip8vm) Load(bytes []byte) {
	copy(v.memory[:], bytes)
}

func (v *Chip8vm) run() {
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
