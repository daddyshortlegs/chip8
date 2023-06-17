package chip8

import (
	"fmt"
)

type VM struct {
	Memory              [4096]byte
	registers           [16]byte
	indexRegister       uint16
	pc                  uint16
	pcIncrementer       int
	display             DisplayInterface
	processInstructions bool
	xCoord              byte
	yCoord              byte
	random              Random
	theStack            *stack
	delayTimer          *DelayTimer
}

type opcodes func(uint16)
type furtherOpcodes func(byte)

func NewVM(display DisplayInterface, random Random) *VM {
	vm := new(VM)
	vm.display = display
	vm.random = random
	vm.pc = 0x200
	vm.pcIncrementer = 2
	vm.theStack = new(stack)
	vm.processInstructions = true
	font := createFont()
	copy(vm.Memory[0x50:], font)
	vm.delayTimer = NewDelayTimer()
	return vm
}

func (v *VM) Load(bytes []byte) {
	copy(v.Memory[0x200:], bytes)
}

func (v *VM) Run() {
	v.delayTimer.Start()
	for {
		if v.processInstructions == true {
			quit := v.fetchAndProcessInstruction()
			if quit == true {
				return
			}
		}

		eventType := v.display.PollEvents()
		if eventType == KeyboardEvent {
			v.processInstructions = true
		}
		if eventType == QuitEvent {
			return
		}
	}
}

func (v *VM) fetchAndProcessInstruction() (quit bool) {
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
	const Random = 0xC
	const Display = 0xD
	const FurtherOperations = 0xF

	//mOps := map[byte]opcodes{
	//	ClearScreen:             v.clearScreen,
	//	Return:                  v.opReturn,
	//	Jump:                    v.jump,
	//	Subroutine:              v.subroutine,
	//	SkipIfEqual:             v.skipIfEqual,
	//	SkipIfNotEqual:          v.skipIfNotEqual,
	//	SkipIfRegistersEqual:    v.skipIfRegistersEqual,
	//	SkipIfRegistersNotEqual: v.skipIfRegistersNotEqual,
	//	SetRegister:             v.setRegister,
	//	AddToRegister:           v.addToRegister,
	//	SetIndexRegister:        v.setIndexRegister,
	//	JumpWithOffset:          v.jumpWithOffset,
	//	Random:                  v.opRandom,
	//	Display:                 v.opDisplay,
	//}

	instr := v.fetchAndIncrement()
	if instr == 0x0000 {
		return true
	}

	i := instruction{instr}
	opCode, _, _, _ := i.extractNibbles()

	v.pcIncrementer = 2

	if instr == ClearScreen {
		v.clearScreen(instr)
	} else if instr == Return {
		v.opReturn(instr)
	} else if opCode == Jump {
		v.jump(instr)
	} else if opCode == Subroutine {
		v.subroutine(instr)
	} else if opCode == SkipIfEqual {
		v.skipIfEqual(instr)
	} else if opCode == SkipIfNotEqual {
		v.skipIfNotEqual(instr)
	} else if opCode == SkipIfRegistersEqual {
		v.skipIfRegistersEqual(instr)
	} else if opCode == SkipIfRegistersNotEqual {
		v.skipIfRegistersNotEqual(instr)
	} else if opCode == SetRegister {
		v.setRegister(instr)
	} else if opCode == AddToRegister {
		v.addToRegister(instr)
	} else if opCode == BitwiseOperations {
		v.executeArthimeticInstrucions(instr)
	} else if opCode == SetIndexRegister {
		v.setIndexRegister(instr)
	} else if opCode == JumpWithOffset {
		v.jumpWithOffset(instr)
	} else if opCode == Random {
		v.opRandom(instr)
	} else if opCode == Display {
		v.opDisplay(instr)
	} else if opCode == FurtherOperations {
		v.furtherOperations(instr)
	}
	return false
}

func (v *VM) furtherOperations(instr uint16) {
	i := instruction{instr}
	_, vx, _, _ := i.extractNibbles()

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
		Bcd:           v.bcd,
		FontChar:      v.fontChar,
		GetKey:        v.getKey,
		AddToIndex:    v.addToIndex,
		Store:         v.store,
		Load:          v.load,
		GetDelayTimer: v.getDelayTimer,
		SetDelayTimer: v.setDelayTimer,
		SetSoundTimer: v.setSoundTimer,
	}

	m[extractSecondByte(instr)](vx)
}

func (v *VM) opDisplay(instr uint16) {
	i := instruction{instr}
	_, vx, vy, opcode2 := i.extractNibbles()

	heightInPixels := opcode2

	v.xCoord = v.registers[vx] & 63
	v.yCoord = v.registers[vy] & 31
	v.registers[15] = 0

	//fmt.Printf("Draw index %X, xreg: %display, yreg: %display, x: %display, y: %display, numBytes: %display\n", v.indexRegister, vx, vy, v.xCoord, v.yCoord, heightInPixels)
	v.display.DrawSprite(v.indexRegister, heightInPixels, v.xCoord, v.yCoord, v.Memory)
}

func (v *VM) opRandom(instr uint16) {
	i := instruction{instr}
	_, vx, _, _ := i.extractNibbles()

	randomNumber := v.random.Generate()
	secondByte := extractSecondByte(instr)
	v.registers[vx] = randomNumber & secondByte
}

func (v *VM) jumpWithOffset(instr uint16) {
	v.pc = uint16(v.registers[0]) + extract12BitNumber(instr)
	v.pcIncrementer = 0
	fmt.Printf("Jump with offset to %X\n", v.pc)
}

func (v *VM) setIndexRegister(instr uint16) {
	v.indexRegister = extract12BitNumber(instr)
}

func (v *VM) skipIfRegistersNotEqual(instr uint16) {
	i := instruction{instr}
	_, vx, vy, _ := i.extractNibbles()

	if v.registers[vx] != v.registers[vy] {
		v.pc += 2
	}
}

func (v *VM) skipIfRegistersEqual(instr uint16) {
	i := instruction{instr}
	_, vx, vy, _ := i.extractNibbles()

	if v.registers[vx] == v.registers[vy] {
		v.pc += 2
	}
}

func (v *VM) skipIfNotEqual(instr uint16) {
	i := instruction{instr}
	_, vx, _, _ := i.extractNibbles()

	if v.registers[vx] != extractSecondByte(instr) {
		v.pc += 2
	}
}

func (v *VM) skipIfEqual(instr uint16) {
	i := instruction{instr}
	_, vx, _, _ := i.extractNibbles()

	if v.registers[vx] == extractSecondByte(instr) {
		v.pc += 2
	}
}

func (v *VM) subroutine(instr uint16) {
	address := extract12BitNumber(instr)
	v.pc = address
	fmt.Printf("Jump to %X\n", v.pc)
	v.theStack.Push(address)
	v.pcIncrementer = 0
}

func (v *VM) jump(instr uint16) {
	v.pc = extract12BitNumber(instr)
	fmt.Printf("Jump to %X\n", v.pc)
	v.pcIncrementer = 0
}

func (v *VM) opReturn(uint16) {
	address, _ := v.theStack.Pop()
	v.pc = address
	fmt.Printf("Stack popped %X\n", v.pc)

	v.pcIncrementer = 0
}

func (v *VM) clearScreen(uint16) {
	println("ClearScreen")
	v.display.ClearScreen()
}

func (v *VM) bcd(vx byte) {
	value := v.registers[vx]
	hundreds, tens, ones := splitNumberIntoUnits(value)

	address := v.indexRegister
	v.Memory[address] = hundreds
	v.Memory[address+1] = tens
	v.Memory[address+2] = ones
}

func (v *VM) fontChar(vx byte) {
	character := v.registers[vx]
	v.indexRegister = 0x50 + uint16(character*5)
}

func (v *VM) getKey(vx byte) {
	// If we get a key then suspend processing of further opcodes
	v.processInstructions = false
	key := v.display.GetKey()
	v.registers[vx] = byte(key)
	fmt.Printf("key returned is %d\n", key)
}

func (v *VM) addToIndex(vx byte) {
	v.indexRegister += uint16(v.registers[vx])
}

func (v *VM) store(vx byte) {
	max := int(vx)
	startMemory := v.indexRegister
	for i := 0; i <= max; i++ {
		v.Memory[startMemory] = v.registers[i]
		startMemory++
	}
}

func (v *VM) load(vx byte) {
	startMemory := v.indexRegister
	for i := 0; i <= int(vx); i++ {
		v.registers[i] = v.Memory[startMemory]
		startMemory++
	}
}

func (v *VM) getDelayTimer(vx byte) {
	// TODO: Test
	// FX07 sets VX to value of the delay timer
	v.registers[vx] = v.delayTimer.timer
}

func (v *VM) setDelayTimer(vx byte) {
	// TODO: Test
	// FX15 set the delay timer to value in VX
	v.delayTimer.timer = v.registers[vx]
}

func (v *VM) setSoundTimer(vx byte) {
	// TODO: Test
	// FX18 sets sound timer to value in VX
	println("vx = ", vx)
}

func (v *VM) executeArthimeticInstrucions(instr uint16) {

	i := instruction{instr}
	_, vx, vy, opcode2 := i.extractNibbles()

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

func (v *VM) fetch() uint16 {
	return bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
}

func (v *VM) fetchAndIncrement() uint16 {
	i := bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
	v.pc += uint16(v.pcIncrementer)
	return i
}

func (v *VM) setRegister(instr uint16) {
	index, secondByte := extractIndexAndValue(instr)
	fmt.Printf("SetRegister %d to %d\n", index, secondByte)
	v.registers[index] = secondByte
}

func (v *VM) addToRegister(instr uint16) {
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

func (v *VM) getXCoordinate() byte {
	return v.xCoord
}

func (v *VM) getYCoordinate() byte {
	return v.yCoord
}
