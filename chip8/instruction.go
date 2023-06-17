package chip8

import "fmt"

type instruction struct {
	instr   uint16
	opCode  byte
	vx      byte
	vy      byte
	opCode2 byte
}

type opcodes func(uint16, *VM)
type furtherOpcodes func(byte, *VM)

const ClearScreen = 0x00E0
const Return = 0x00EE
const Jump = 0x1
const Subroutine = 0x2
const SkipIfEqual = 0x3
const SkipIfNotEqual = 0x4
const SkipIfRegistersEqual = 0x5
const SkipIfRegistersNotEqual = 0x9
const SetRegister = 0x6
const AddToRegister = 0x7
const BitwiseOperations = 0x8
const SetIndexRegister = 0xA
const JumpWithOffset = 0xB
const OpRandom = 0xC
const Display = 0xD
const FurtherOperations = 0xF

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

func extractIndexAndValue(instr uint16) (byte, byte) {
	firstByte := extractFirstByte(instr)
	index := getRightNibble(firstByte)
	secondByte := extractSecondByte(instr)
	return index, secondByte
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

func (i *instruction) furtherOperations(instr uint16, v *VM) {
	const Bcd = 0x33
	const FontChar = 0x29
	const GetKey = 0x0A
	const AddToIndex = 0x1E
	const Store = 0x55
	const Load = 0x65
	const GetDelayTimer = 0x07
	const SetDelayTimer = 0x15
	const SetSoundTimer = 0x18

	m := map[byte]furtherOpcodes{
		Bcd:           i.bcd,
		FontChar:      i.fontChar,
		GetKey:        i.getKey,
		AddToIndex:    i.addToIndex,
		Store:         i.store,
		Load:          i.load,
		GetDelayTimer: i.getDelayTimer,
		SetDelayTimer: i.setDelayTimer,
		SetSoundTimer: i.setSoundTimer,
	}

	m[extractSecondByte(instr)](i.vx, v)
}

func (i *instruction) bcd(vx byte, v *VM) {
	value := v.registers[vx]
	hundreds, tens, ones := splitNumberIntoUnits(value)

	address := v.indexRegister
	v.Memory[address] = hundreds
	v.Memory[address+1] = tens
	v.Memory[address+2] = ones
}

func (i *instruction) fontChar(vx byte, v *VM) {
	character := v.registers[vx]
	v.indexRegister = 0x50 + uint16(character*5)
}

func (i *instruction) getKey(vx byte, v *VM) {
	// If we get a key then suspend processing of further opcodes
	v.processInstructions = false
	key := v.display.GetKey()
	v.registers[vx] = byte(key)
	fmt.Printf("key returned is %d\n", key)
}

func (i *instruction) addToIndex(vx byte, v *VM) {
	v.indexRegister += uint16(v.registers[vx])
}

func (i *instruction) store(vx byte, v *VM) {
	max := int(vx)
	startMemory := v.indexRegister
	for i := 0; i <= max; i++ {
		v.Memory[startMemory] = v.registers[i]
		startMemory++
	}
}

func (i *instruction) load(vx byte, v *VM) {
	startMemory := v.indexRegister
	for i := 0; i <= int(vx); i++ {
		v.registers[i] = v.Memory[startMemory]
		startMemory++
	}
}

func (i *instruction) getDelayTimer(vx byte, v *VM) {
	// TODO: Test
	// FX07 sets VX to value of the delay timer
	v.registers[vx] = v.delayTimer.timer
}

func (i *instruction) setDelayTimer(vx byte, v *VM) {
	// TODO: Test
	// FX15 set the delay timer to value in VX
	v.delayTimer.timer = v.registers[vx]
}

func (i *instruction) setSoundTimer(vx byte, v *VM) {
	// TODO: Test
	// FX18 sets sound timer to value in VX
	println("vx = ", vx)
}
