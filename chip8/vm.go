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
	theStack                *stack
}

func NewVM(display DisplayInterface, random Random) *VM {
	vm := new(VM)
	vm.display = display
	vm.random = random
	vm.pc = 0x200
	vm.theStack = new(stack)
	font := createFont()
	copy(vm.Memory[0x50:], font)
	return vm
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
		} else if instr == 0x00EE {
			address, _ := v.theStack.Pop()
			v.pc = address
			v.previousInstructionJump = true

		} else if opCode == 0x1 {
			v.pc = extract12BitNumber(instr)
			fmt.Printf("Jump to %X\n", v.pc)
			v.previousInstructionJump = true
		} else if opCode == 0x2 {
			address := extract12BitNumber(instr)
			v.pc = address
			fmt.Printf("Jump to %X\n", v.pc)
			v.theStack.Push(address)
			v.previousInstructionJump = true
		} else if opCode == 0x3 {
			if v.registers[vx] == extractSecondByte(instr) {
				v.pc += 2
			}
		} else if opCode == 0x4 {
			if v.registers[vx] != extractSecondByte(instr) {
				v.pc += 2
			}
		} else if opCode == 0x5 {
			if v.registers[vx] == v.registers[vy] {
				v.pc += 2
			}
		} else if opCode == 0x9 {
			if v.registers[vx] != v.registers[vy] {
				v.pc += 2
			}
		} else if opCode == 0x6 {
			v.setRegister(instr)
		} else if opCode == 0x7 {
			v.addToRegister(instr)
		} else if opCode == 0x8 {
			v.executeArthimeticInstrucions(opcode2, vx, vy)
		} else if opCode == 0xA {
			v.indexRegister = extract12BitNumber(instr)
			fmt.Printf("Set Index Register %X\n", v.indexRegister)
		} else if opCode == 0xC {
			randomNumber := v.random.Generate()
			firstByte := extractFirstByte(instr)
			index := getRightNibble(firstByte)
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
				hundreds, tens, ones := splitNumberIntoUnits(value)

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

func (v *VM) executeArthimeticInstrucions(opcode2 byte, vx byte, vy byte) {
	const setVxToVy = 0x0
	const binaryOr = 0x1
	const binaryAnd = 0x2
	const logicalXor = 0x3
	const addToVx = 0x4
	const subtractFromVx = 0x5
	const shiftRight = 0x6
	const subtractFromVy = 0x7
	const shiftLeft = 0xE

	if opcode2 == setVxToVy {
		v.registers[vx] = v.registers[vy]
	} else if opcode2 == binaryOr {
		v.registers[vx] = v.registers[vx] | v.registers[vy]
	} else if opcode2 == binaryAnd {
		v.registers[vx] = v.registers[vx] & v.registers[vy]
	} else if opcode2 == logicalXor {
		v.registers[vx] = v.registers[vx] ^ v.registers[vy]
	} else if opcode2 == addToVx {
		vxRegister := v.registers[vx]
		vyRegister := v.registers[vy]

		v.registers[vx] = vxRegister + vyRegister
		var sum = uint16(vxRegister) + uint16(vyRegister)
		if sum > 255 {
			v.registers[15] = 1
		} else {
			v.registers[15] = 0
		}
	} else if opcode2 == subtractFromVx {
		vxRegister := v.registers[vx]
		vyRegister := v.registers[vy]
		v.registers[vx] = vxRegister - vyRegister
		var underflowFlag byte = 1
		if vxRegister < vyRegister {
			underflowFlag = 0
		}
		v.registers[15] = underflowFlag
	} else if opcode2 == shiftRight {
		overflow := v.registers[vy] & 0b00000001
		v.registers[15] = overflow
		v.registers[vx] = v.registers[vy] >> 1
	} else if opcode2 == subtractFromVy {
		vxRegister := v.registers[vx]
		vyRegister := v.registers[vy]
		v.registers[vx] = vyRegister - vxRegister

		var underflowFlag byte = 1
		if vyRegister < vxRegister {
			underflowFlag = 0
		}
		v.registers[15] = underflowFlag

	} else if opcode2 == shiftLeft {
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
	index, secondByte := extractIndexAndValue(instr)
	fmt.Printf("SetRegister %display to %display\n", index, secondByte)
	v.registers[index] = secondByte
}

func (v *VM) addToRegister(instr uint16) {
	index, secondByte := extractIndexAndValue(instr)
	fmt.Printf("Add To Register [%display] value %display\n", index, secondByte)
	v.registers[index] += secondByte
}

func extractIndexAndValue(instr uint16) (byte, byte) {
	firstByte := extractFirstByte(instr)
	index := getRightNibble(firstByte)
	secondByte := extractSecondByte(instr)
	return index, secondByte
}

func (v *VM) getXCoordinate() byte {
	return v.xCoord
}

func (v *VM) getYCoordinate() byte {
	return v.yCoord
}
