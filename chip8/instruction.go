package chip8

import "fmt"

type instruction struct {
	instr   uint16
	opCode  byte
	vx      byte
	vy      byte
	opCode2 byte
}

func NewInstruction(instr uint16) *instruction {
	i := new(instruction)
	i.extractNibbles(instr)
	return i
}

func (i *instruction) extractNibbles(instr uint16) {
	i.opCode = extractNibble(instr)
	i.vx = getRightNibble(extractFirstByte(instr))
	secondByte := extractSecondByte(instr)
	i.vy = getLeftNibble(secondByte)
	i.opCode2 = getRightNibble(secondByte)
}

func (i *instruction) execute(instr uint16, v *VM) {
	mOps := map[byte]opcodes{
		ClearScreen:             i.clearScreen,
		Return:                  i.opReturn,
		Jump:                    i.jump,
		Subroutine:              i.subroutine,
		SkipIfEqual:             i.skipIfEqual,
		SkipIfNotEqual:          i.skipIfNotEqual,
		SkipIfRegistersEqual:    i.skipIfRegistersEqual,
		SkipIfRegistersNotEqual: i.skipIfRegistersNotEqual,
		SetRegister:             i.setRegister,
		AddToRegister:           i.addToRegister,
		SetIndexRegister:        i.setIndexRegister,
		JumpWithOffset:          i.jumpWithOffset,
		OpRandom:                i.opRandom,
		Display:                 i.opDisplay,
		BitwiseOperations:       i.executeArithmeticInstructions,
		FurtherOperations:       i.furtherOperations,
	}
	mOps[i.opCode](instr, v)
}

func (i *instruction) opDisplay(instr uint16, v *VM) {
	heightInPixels := i.opCode2

	v.xCoord = v.registers[i.vx] & 63
	v.yCoord = v.registers[i.vy] & 31
	v.registers[15] = 0

	//fmt.Printf("Draw index %X, xreg: %display, yreg: %display, x: %display, y: %display, numBytes: %display\n", v.indexRegister, vx, vy, v.xCoord, v.yCoord, heightInPixels)
	v.display.DrawSprite(v.indexRegister, heightInPixels, v.xCoord, v.yCoord, v.Memory)
}

func (i *instruction) opRandom(instr uint16, v *VM) {
	randomNumber := v.random.Generate()
	secondByte := extractSecondByte(instr)
	v.registers[i.vx] = randomNumber & secondByte
}

func (i *instruction) jumpWithOffset(instr uint16, v *VM) {
	v.pc = uint16(v.registers[0]) + extract12BitNumber(instr)
	v.pcIncrementer = 0
	fmt.Printf("Jump with offset to %X\n", v.pc)
}

func (i *instruction) setIndexRegister(instr uint16, v *VM) {
	v.indexRegister = extract12BitNumber(instr)
}

func (i *instruction) skipIfRegistersNotEqual(instr uint16, v *VM) {
	if v.registers[i.vx] != v.registers[i.vy] {
		v.pc += 2
	}
}

func (i *instruction) setRegister(instr uint16, v *VM) {
	index, secondByte := extractIndexAndValue(instr)
	fmt.Printf("SetRegister %d to %d\n", index, secondByte)
	v.registers[index] = secondByte
}

func (i *instruction) addToRegister(instr uint16, v *VM) {
	index, secondByte := extractIndexAndValue(instr)
	fmt.Printf("Add To Register [%d] value %d\n", index, secondByte)
	v.registers[index] += secondByte
}

func (i *instruction) skipIfRegistersEqual(instr uint16, v *VM) {
	if v.registers[i.vx] == v.registers[i.vy] {
		v.pc += 2
	}
}

func (i *instruction) skipIfNotEqual(instr uint16, v *VM) {
	if v.registers[i.vx] != extractSecondByte(instr) {
		v.pc += 2
	}
}

func (i *instruction) skipIfEqual(instr uint16, v *VM) {
	if v.registers[i.vx] == extractSecondByte(instr) {
		v.pc += 2
	}
}

func (i *instruction) subroutine(instr uint16, v *VM) {
	address := extract12BitNumber(instr)
	v.pc = address
	fmt.Printf("Jump to %X\n", v.pc)
	v.theStack.Push(address)
	v.pcIncrementer = 0
}

func (i *instruction) jump(instr uint16, v *VM) {
	v.pc = extract12BitNumber(instr)
	fmt.Printf("Jump to %X\n", v.pc)
	v.pcIncrementer = 0
}

func (i *instruction) opReturn(instr uint16, v *VM) {
	address, _ := v.theStack.Pop()
	v.pc = address
	fmt.Printf("Stack popped %X\n", v.pc)

	v.pcIncrementer = 0
}

func (i *instruction) clearScreen(instr uint16, v *VM) {
	println("ClearScreen")
	v.display.ClearScreen()
}

func (i *instruction) executeArithmeticInstructions(instr uint16, v *VM) {
	const setVxToVy = 0x0
	const binaryOr = 0x1
	const binaryAnd = 0x2
	const logicalXor = 0x3
	const addToVx = 0x4
	const subtractFromVx = 0x5
	const shiftRight = 0x6
	const subtractFromVy = 0x7
	const shiftLeft = 0xE

	if i.opCode2 == setVxToVy {
		v.registers[i.vx] = v.registers[i.vy]
	} else if i.opCode2 == binaryOr {
		v.registers[i.vx] = v.registers[i.vx] | v.registers[i.vy]
	} else if i.opCode2 == binaryAnd {
		v.registers[i.vx] = v.registers[i.vx] & v.registers[i.vy]
	} else if i.opCode2 == logicalXor {
		v.registers[i.vx] = v.registers[i.vx] ^ v.registers[i.vy]
	} else if i.opCode2 == addToVx {
		vxRegister := v.registers[i.vx]
		vyRegister := v.registers[i.vy]

		v.registers[i.vx] = vxRegister + vyRegister
		var sum = uint16(vxRegister) + uint16(vyRegister)
		if sum > 255 {
			v.registers[15] = 1
		} else {
			v.registers[15] = 0
		}
	} else if i.opCode2 == subtractFromVx {
		vxRegister := v.registers[i.vx]
		vyRegister := v.registers[i.vy]
		v.registers[i.vx] = vxRegister - vyRegister
		var underflowFlag byte = 1
		if vxRegister < vyRegister {
			underflowFlag = 0
		}
		v.registers[15] = underflowFlag
	} else if i.opCode2 == shiftRight {
		overflow := v.registers[i.vy] & 0b00000001
		v.registers[15] = overflow
		v.registers[i.vx] = v.registers[i.vy] >> 1
	} else if i.opCode2 == subtractFromVy {
		vxRegister := v.registers[i.vx]
		vyRegister := v.registers[i.vy]
		v.registers[i.vx] = vyRegister - vxRegister

		var underflowFlag byte = 1
		if vyRegister < vxRegister {
			underflowFlag = 0
		}
		v.registers[15] = underflowFlag

	} else if i.opCode2 == shiftLeft {
		overflow := v.registers[i.vy] & 0b10000000
		v.registers[15] = overflow >> 7
		v.registers[i.vx] = v.registers[i.vy] << 1
	}
}
