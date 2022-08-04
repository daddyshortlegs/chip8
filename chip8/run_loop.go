package chip8

import (
	"fmt"
)

type VM struct {
	Memory                  [4096]byte
	registers               [16]byte
	indexRegister           uint16
	pc                      uint16
	display                 DisplayInterface
	previousInstructionJump bool
	xCoord                  byte
	yCoord                  byte
	random                  Random
}

func NewVM(display DisplayInterface, random Random) *VM {
	vm := new(VM)
	vm.display = display
	vm.random = random
	vm.pc = 0x200
	font := createFont()
	copy(vm.Memory[0x50:], font)
	return vm
}

type DisplayInterface interface {
	ClearScreen()
	DrawSprite(chip8 *VM, address uint16, numberOfBytes byte, x byte, y byte)
	PollEvents() bool
}

func (v *VM) Load(bytes []byte) {
	copy(v.Memory[0x200:], bytes)
}

func (v *VM) Run() {
	v.previousInstructionJump = false
	for {
		instr := v.fetchNextInstruction()
		if instr == 0x0000 {
			break
		}

		i := instruction{instr}
		opCode, vx, vy, opcode2 := i.extractNibbles()

		if instr == 0x00E0 {
			println("ClearScreen")
			v.display.ClearScreen()
		} else if opCode == 0x1 {
			v.jump(instr)
			v.previousInstructionJump = true
			//continue
		} else if opCode == 0x6 {
			v.setRegister(instr)
		} else if opCode == 0x7 {
			v.addToRegister(instr)
		} else if opCode == 0x8 {
			v.executeArthimeticInstrucions(opcode2, vx, vy)
		} else if opCode == 0xA {
			v.setIndexRegister(instr)
		} else if opCode == 0xC {
			randomNumber := v.random.Generate()
			index := v.getRegisterIndex(instr)
			secondByte := extractSecondByte(instr)
			v.registers[index] = randomNumber & secondByte
		} else if opCode == 0xD {
			numberOfBytes := opcode2

			v.xCoord = v.registers[vx] & 63
			v.yCoord = v.registers[vy] & 31
			v.registers[15] = 0

			fmt.Printf("Draw index %X, xreg: %display, yreg: %display, x: %display, y: %display, numBytes: %display\n", v.indexRegister, vx, vy, v.xCoord, v.yCoord, numberOfBytes)
			v.display.DrawSprite(v, v.indexRegister, numberOfBytes, v.xCoord, v.yCoord)
		} else if opCode == 0xF {

			secondByte := extractSecondByte(instr)

			if secondByte == 0x33 {
				value := v.registers[vx]
				hundreds, tens, ones := splitNumber(value)

				address := v.indexRegister
				v.Memory[address] = hundreds
				v.Memory[address+1] = tens
				v.Memory[address+2] = ones

			} else if secondByte == 0x29 {
				character := v.registers[vx]
				v.indexRegister = 0x50 + uint16(character*5)
			}

		} else {
			v.previousInstructionJump = false
		}
		quit := v.display.PollEvents()
		if quit == false {
			return
		}
	}
}

func splitNumber(number byte) (byte, byte, byte) {
	hundreds := number / 100
	tens := (number % 100) / 10
	ones := number % 10
	return hundreds, tens, ones
}

func (v *VM) executeArthimeticInstrucions(opcode2 byte, vx byte, vy byte) {
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
	} else if opcode2 == 5 {
		vxRegister := v.registers[vx]
		vyRegister := v.registers[vy]
		v.registers[vx] = vxRegister - vyRegister
		var underflowFlag byte = 1
		if vxRegister < vyRegister {
			underflowFlag = 0
		}
		v.registers[15] = underflowFlag
	} else if opcode2 == 6 {
		overflow := v.registers[vy] & 0b00000001
		v.registers[15] = overflow
		v.registers[vx] = v.registers[vy] >> 1
	} else if opcode2 == 7 {
		vxRegister := v.registers[vx]
		vyRegister := v.registers[vy]
		v.registers[vx] = vyRegister - vxRegister

		var underflowFlag byte = 1
		if vyRegister < vxRegister {
			underflowFlag = 0
		}
		v.registers[15] = underflowFlag

	} else if opcode2 == 0xE {
		overflow := v.registers[vy] & 0b10000000
		v.registers[15] = overflow >> 7
		v.registers[vx] = v.registers[vy] << 1
	}
}

func (v *VM) fetchNextInstruction() uint16 {
	if v.previousInstructionJump == false {
		return v.fetchAndIncrement()
	}
	return v.fetch()
}

func (v *VM) fetch() uint16 {
	return bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
}

func (v *VM) fetchAndIncrement() uint16 {
	i := bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
	v.pc += 2
	return i
}

func (v *VM) setRegister(instr uint16) {
	index := v.getRegisterIndex(instr)
	secondByte := extractSecondByte(instr)
	fmt.Printf("SetRegister %display to %display\n", index, secondByte)
	v.registers[index] = secondByte
}

func (v *VM) addToRegister(instr uint16) {
	index := v.getRegisterIndex(instr)
	secondByte := extractSecondByte(instr)
	fmt.Printf("Add To Register [%display] value %display\n", index, secondByte)
	v.registers[index] += secondByte
}

func (v *VM) getRegisterIndex(instr uint16) byte {
	firstByte := extractFirstByte(instr)
	return getRightNibble(firstByte)
}

func (v *VM) setIndexRegister(instr uint16) {
	v.indexRegister = extract12BitNumber(instr)
	fmt.Printf("Set Index Register %X\n", v.indexRegister)
}

func (v *VM) jump(address uint16) {
	v.pc = extract12BitNumber(address)
	//fmt.Printf("Jump to %X\n", v.pc)
}

func (v *VM) getXCoordinate() byte {
	return v.xCoord
}

func (v *VM) getYCoordinate() byte {
	return v.yCoord
}
