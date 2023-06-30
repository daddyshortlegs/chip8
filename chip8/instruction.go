package chip8

import "fmt"

type Instruction struct {
	instr           uint16
	opCode          byte
	vx              byte
	vy              byte
	opCode2         byte
	opCodeFunctions map[byte]Opcode
}

type Opcode struct {
	name string
	i    instruction
}

type instruction func(uint16, *VM)
type arithmeticOpcodes func(*VM)
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

const setVxToVy = 0x00
const binaryOr = 0x01
const binaryAnd = 0x02
const logicalXor = 0x03
const addToVx = 0x04
const subtractFromVx = 0x05
const shiftRight = 0x06
const subtractFromVy = 0x07
const shiftLeft = 0x0E

const Bcd = 0x33
const FontChar = 0x29
const GetKey = 0x0A
const AddToIndex = 0x1E
const Store = 0x55
const Load = 0x65
const GetDelayTimer = 0x07
const SetDelayTimer = 0x15
const SetSoundTimer = 0x18

func NewInstruction(instr uint16) *Instruction {
	i := new(Instruction)
	i.extractNibbles(instr)
	i.opCodeFunctions = i.setupOpCodes()
	return i
}

func (i *Instruction) setupOpCodes() map[byte]Opcode {
	opCodeFunctions := map[byte]Opcode{
		ClearScreen:             Opcode{name: "ClearScreen", i: i.clearScreen},
		Return:                  Opcode{name: "Return", i: i.opReturn},
		Jump:                    Opcode{name: "Jump", i: i.jump},
		Subroutine:              Opcode{name: "Subroutine", i: i.subroutine},
		SkipIfEqual:             Opcode{name: "SkipIfEqual", i: i.skipIfEqual},
		SkipIfNotEqual:          Opcode{name: "SkipIfNotEqual", i: i.skipIfNotEqual},
		SkipIfRegistersEqual:    Opcode{name: "SkipIfRegistersEqual", i: i.skipIfRegistersEqual},
		SkipIfRegistersNotEqual: Opcode{name: "SkipIfRegistersNotEqual", i: i.skipIfRegistersNotEqual},
		SetRegister:             Opcode{name: "SetRegister", i: i.setRegister},
		AddToRegister:           Opcode{name: "AddToRegister", i: i.addToRegister},
		SetIndexRegister:        Opcode{name: "SetIndexRegister", i: i.setIndexRegister},
		JumpWithOffset:          Opcode{name: "JumpWithOffset", i: i.jumpWithOffset},
		OpRandom:                Opcode{name: "OpRandom", i: i.opRandom},
		Display:                 Opcode{name: "Display", i: i.opDisplay},
		BitwiseOperations:       Opcode{name: "BitwiseOperations", i: i.executeArithmeticInstructions},
		FurtherOperations:       Opcode{name: "FurtherOperations", i: i.furtherOperations},
	}
	return opCodeFunctions
}

func (i *Instruction) extractNibbles(instr uint16) {
	i.opCode = extractNibble(instr)
	i.vx = getRightNibble(extractFirstByte(instr))
	secondByte := extractSecondByte(instr)
	i.vy = getLeftNibble(secondByte)
	i.opCode2 = getRightNibble(secondByte)
}

func (i *Instruction) execute(instr uint16, v *VM) {
	fmt.Printf("> %s\n", i.getOpcodeName())
	fmt.Printf(">>> %x\n", i.opCode)

	function := i.getInstructionFromOpcode()
	if function == nil {
		fmt.Printf("Unknown instruction %x", i.opCode)
	} else {
		function(instr, v)
	}
}

func (i *Instruction) getInstructionFromOpcode() instruction {
	function := i.opCodeFunctions[i.opCode]
	return function.i
}

func (i *Instruction) getOpcodeName() string {
	opCodeNames := map[byte]string{
		ClearScreen:             "ClearScreen",
		Return:                  "Return",
		Jump:                    "Jump",
		Subroutine:              "Subroutine",
		SkipIfEqual:             "SkipIfEqual",
		SkipIfNotEqual:          "SkipIfNotEqual",
		SkipIfRegistersEqual:    "SkipIfRegistersEqual",
		SkipIfRegistersNotEqual: "SkipIfRegistersNotEqual",
		SetRegister:             "SetRegister",
		AddToRegister:           "AddToRegister",
		SetIndexRegister:        "SetIndexRegister",
		JumpWithOffset:          "JumpWithOffset",
		OpRandom:                "Random",
		Display:                 "Display",
		BitwiseOperations:       "Bitwise",
		FurtherOperations:       "Further",
	}

	functionName := opCodeNames[i.opCode]
	return functionName
}

func (i *Instruction) opDisplay(_ uint16, v *VM) {
	heightInPixels := i.opCode2

	v.xCoord = v.registers[i.vx] & 63
	v.yCoord = v.registers[i.vy] & 31
	v.registers[15] = 0

	//fmt.Printf("Draw index %X, xreg: %d, yreg: %d, x: %d, y: %d, numBytes: %d\n", v.indexRegister, i.vx, i.vy, v.xCoord, v.yCoord, heightInPixels)
	v.display.DrawSprite(v.indexRegister, heightInPixels, v.xCoord, v.yCoord, v.Memory)
}

func (i *Instruction) opRandom(instr uint16, v *VM) {
	randomNumber := v.random.Generate()
	secondByte := extractSecondByte(instr)
	v.registers[i.vx] = randomNumber & secondByte
}

func (i *Instruction) jumpWithOffset(instr uint16, v *VM) {
	v.pc = uint16(v.registers[0]) + extract12BitNumber(instr)
	v.pcIncrementer = 0
	fmt.Printf("Jump with offset to %X\n", v.pc)
}

func (i *Instruction) setIndexRegister(instr uint16, v *VM) {
	v.indexRegister = extract12BitNumber(instr)
}

func (i *Instruction) skipIfRegistersNotEqual(_ uint16, v *VM) {
	if v.registers[i.vx] != v.registers[i.vy] {
		v.pc += 2
	}
}

func (i *Instruction) setRegister(instr uint16, v *VM) {
	_, secondByte := extractIndexAndValue(instr)
	//fmt.Printf("SetRegister %d to %d\n", index, secondByte)
	v.registers[i.vx] = secondByte
}

func (i *Instruction) addToRegister(instr uint16, v *VM) {
	_, secondByte := extractIndexAndValue(instr)
	//fmt.Printf("Add To Register [%d] value %d\n", index, secondByte)
	v.registers[i.vx] += secondByte
}

func extractIndexAndValue(instr uint16) (byte, byte) {
	firstByte := extractFirstByte(instr)
	index := getRightNibble(firstByte)
	secondByte := extractSecondByte(instr)
	return index, secondByte
}

func (i *Instruction) skipIfRegistersEqual(_ uint16, v *VM) {
	if v.registers[i.vx] == v.registers[i.vy] {
		v.pc += 2
	}
}

func (i *Instruction) skipIfNotEqual(instr uint16, v *VM) {
	if v.registers[i.vx] != extractSecondByte(instr) {
		v.pc += 2
	}
}

func (i *Instruction) skipIfEqual(instr uint16, v *VM) {
	if v.registers[i.vx] == extractSecondByte(instr) {
		v.pc += 2
	}
}

func (i *Instruction) subroutine(instr uint16, v *VM) {
	address := extract12BitNumber(instr)
	v.theStack.Push(v.pc)
	v.pc = address
	v.pcIncrementer = 0
}

func (i *Instruction) jump(instr uint16, v *VM) {
	v.pc = extract12BitNumber(instr)
	v.pcIncrementer = 0
}

func (i *Instruction) opReturn(_ uint16, v *VM) {
	address, _ := v.theStack.Pop()
	v.pc = address
	fmt.Printf("Stack popped %X\n", v.pc)

	v.pcIncrementer = 0
}

func (i *Instruction) clearScreen(_ uint16, v *VM) {
	println("ClearScreen")
	v.display.ClearScreen()
}

func (i *Instruction) executeArithmeticInstructions(_ uint16, v *VM) {
	opCodeFunctions := map[byte]arithmeticOpcodes{
		setVxToVy:      i.setVxToVy,
		binaryOr:       i.or,
		binaryAnd:      i.and,
		logicalXor:     i.xOr,
		addToVx:        i.addToVx,
		subtractFromVx: i.subtractFromVx,
		shiftRight:     i.shiftRight,
		subtractFromVy: i.subtractFromVy,
		shiftLeft:      i.shiftLeft,
	}
	opCodeFunctions[i.opCode2](v)
}

func (i *Instruction) setVxToVy(v *VM) {
	v.registers[i.vx] = v.registers[i.vy]
}

func (i *Instruction) or(v *VM) {
	v.registers[i.vx] = v.registers[i.vx] | v.registers[i.vy]
}

func (i *Instruction) and(v *VM) {
	v.registers[i.vx] = v.registers[i.vx] & v.registers[i.vy]
}

func (i *Instruction) xOr(v *VM) {
	v.registers[i.vx] = v.registers[i.vx] ^ v.registers[i.vy]
}

func (i *Instruction) addToVx(v *VM) {
	vxRegister := v.registers[i.vx]
	vyRegister := v.registers[i.vy]

	v.registers[i.vx] = vxRegister + vyRegister
	var sum = uint16(vxRegister) + uint16(vyRegister)
	if sum > 255 {
		v.registers[15] = 1
	} else {
		v.registers[15] = 0
	}
}

func (i *Instruction) subtractFromVx(v *VM) {
	vxRegister := v.registers[i.vx]
	vyRegister := v.registers[i.vy]
	v.registers[i.vx] = vxRegister - vyRegister
	var underflowFlag byte = 1
	if vxRegister < vyRegister {
		underflowFlag = 0
	}
	v.registers[15] = underflowFlag
}

func (i *Instruction) shiftRight(v *VM) {
	overflow := v.registers[i.vy] & 0b00000001
	v.registers[15] = overflow
	v.registers[i.vx] = v.registers[i.vy] >> 1
}

func (i *Instruction) subtractFromVy(v *VM) {
	vxRegister := v.registers[i.vx]
	vyRegister := v.registers[i.vy]
	v.registers[i.vx] = vyRegister - vxRegister

	var underflowFlag byte = 1
	if vyRegister < vxRegister {
		underflowFlag = 0
	}
	v.registers[15] = underflowFlag
}

func (i *Instruction) shiftLeft(v *VM) {
	overflow := v.registers[i.vy] & 0b10000000
	v.registers[15] = overflow >> 7
	v.registers[i.vx] = v.registers[i.vy] << 1
}

func (i *Instruction) furtherOperations(instr uint16, v *VM) {
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

	names := map[byte]string{
		Bcd:           "bcd",
		FontChar:      "fontChar",
		GetKey:        "getKey",
		AddToIndex:    "addToIndex",
		Store:         "store",
		Load:          "load",
		GetDelayTimer: "getDelayTimer",
		SetDelayTimer: "setDelayTimer",
		SetSoundTimer: "setSoundTimer",
	}
	secondByte := extractSecondByte(instr)

	functionName := names[secondByte]
	fmt.Printf(">>> %s\n", functionName)
	fmt.Printf(">>> %x\n", secondByte)

	f := m[secondByte]
	f(i.vx, v)
}

func (i *Instruction) bcd(vx byte, v *VM) {
	value := v.registers[vx]
	hundreds, tens, ones := splitNumberIntoUnits(value)

	address := v.indexRegister
	v.Memory[address] = hundreds
	v.Memory[address+1] = tens
	v.Memory[address+2] = ones
}

func (i *Instruction) fontChar(vx byte, v *VM) {
	character := v.registers[vx]
	v.indexRegister = 0x50 + uint16(character*5)
}

func (i *Instruction) getKey(vx byte, v *VM) {
	// If we get a key then suspend processing of further instruction
	v.processInstructions = false
	key := v.display.GetKey()
	v.registers[vx] = byte(key)
}

func (i *Instruction) addToIndex(vx byte, v *VM) {
	v.indexRegister += uint16(v.registers[vx])
}

func (i *Instruction) store(vx byte, v *VM) {
	max := int(vx)
	startMemory := v.indexRegister
	for i := 0; i <= max; i++ {
		v.Memory[startMemory] = v.registers[i]
		startMemory++
	}
}

func (i *Instruction) load(vx byte, v *VM) {
	startMemory := v.indexRegister
	for i := 0; i <= int(vx); i++ {
		v.registers[i] = v.Memory[startMemory]
		startMemory++
	}
}

func (i *Instruction) getDelayTimer(vx byte, v *VM) {
	// TODO: Test
	// FX07 sets VX to value of the delay timer
	v.registers[vx] = v.delayTimer.timer
}

func (i *Instruction) setDelayTimer(vx byte, v *VM) {
	// TODO: Test
	// FX15 set the delay timer to value in VX
	v.setDelayTimer(v.registers[vx])
	//v.delayTimer.timer = v.registers[vx]
}

func (i *Instruction) setSoundTimer(vx byte, _ *VM) {
	// TODO: Test
	// FX18 sets sound timer to value in VX
	println("vx = ", vx)
}
