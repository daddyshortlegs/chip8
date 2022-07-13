package chip8

import "fmt"

type Chip8vm struct {
	Memory                  [4096]byte
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
	DrawSprite(chip8 *Chip8vm, address uint16, numberOfBytes byte, x byte, y byte)
}

func (v *Chip8vm) Init() {
	v.pc = 0x200
	font := createFont()
	copy(v.Memory[0x50:], font)
}

func (v *Chip8vm) SetDisplay(d Display) {
	v.d = d
}

func (v *Chip8vm) Load(bytes []byte) {
	copy(v.Memory[0x200:], bytes)
}

func (v *Chip8vm) Run() {
	v.previousInstructionJump = false
	for {
		instr := v.fetchNextInstruction()

		if instr == 0x0000 {
			break
		}

		firstNibble := extractNibble(instr)
		firstByte := extractFirstByte(instr)
		secondByte := extractSecondByte(instr)
		xRegister := v.getRightNibble(firstByte)
		yRegister := v.getLeftNibble(secondByte)

		if instr == 0x00E0 {
			println("ClearScreen")
			v.d.ClearScreen()
		} else if firstNibble == 0x1 {
			v.jump(instr)
			v.previousInstructionJump = true
			return
			//continue
		} else if firstNibble == 0x6 {
			v.setRegister(instr)
		} else if firstNibble == 0x7 {
			v.addToRegister(instr)
		} else if firstNibble == 0x8 {
			secondByte2 := extractSecondByte(instr)
			opcode2 := v.getRightNibble(secondByte2)
			vx := v.getRegisterIndex(instr)
			vy := v.getLeftNibble(secondByte2)

			if opcode2 == 0 {
				v.registers[vx] = v.registers[vy]
			} else if opcode2 == 1 {
				v.registers[vx] = v.registers[vx] | v.registers[vy]
			} else if opcode2 == 2 {
				v.registers[vx] = v.registers[vx] & v.registers[vy]
			} else if opcode2 == 3 {
				v.registers[vx] = v.registers[vx] ^ v.registers[vy]
			} else if opcode2 == 4 {
				vxRegister := v.registers[vx]
				vyRegister := v.registers[vy]

				v.registers[vx] = vxRegister + vyRegister
				var sum uint16
				sum = uint16(vxRegister) + uint16(vyRegister)
				if sum > 255 {
					v.registers[15] = 1
				} else {
					v.registers[15] = 0
				}
			}

		} else if firstNibble == 0xA {
			v.setIndexRegister(instr)
		} else if firstNibble == 0xD {
			numberOfBytes := v.getRightNibble(secondByte)

			v.xCoord = v.registers[xRegister] & 63
			v.yCoord = v.registers[yRegister] & 31
			v.registers[15] = 0

			fmt.Printf("Draw index %X, xreg: %d, yreg: %d, x: %d, y: %d, numBytes: %d\n", v.indexRegister, xRegister, yRegister, v.xCoord, v.yCoord, numberOfBytes)
			v.d.DrawSprite(v, v.indexRegister, numberOfBytes, v.xCoord, v.yCoord)
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
	i := bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
	return i
}

func (v *Chip8vm) fetchAndIncrement() uint16 {
	i := bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
	v.pc += 2
	return i
}

func (v *Chip8vm) setRegister(instr uint16) {
	index := v.getRegisterIndex(instr)
	secondByte := extractSecondByte(instr)
	fmt.Printf("SetRegister %d to %d\n", index, secondByte)
	v.registers[index] = secondByte
}

func (v *Chip8vm) addToRegister(instr uint16) {
	index := v.getRegisterIndex(instr)
	secondByte := extractSecondByte(instr)
	fmt.Printf("Add To Register [%d] value %d\n", index, secondByte)
	v.registers[index] += secondByte
}

func (v *Chip8vm) getRegisterIndex(instr uint16) byte {
	firstByte := extractFirstByte(instr)
	return v.getRightNibble(firstByte)
}

func (v *Chip8vm) setIndexRegister(instr uint16) {
	v.indexRegister = extract12BitNumber(instr)
	fmt.Printf("Set Index Register %X\n", v.indexRegister)
}

func (v *Chip8vm) jump(address uint16) {
	v.pc = extract12BitNumber(address)
	fmt.Printf("Jump to %X\n", v.pc)
}

func (v *Chip8vm) getXCoordinate() byte {
	return v.xCoord
}

func (v *Chip8vm) getYCoordinate() byte {
	return v.yCoord
}
