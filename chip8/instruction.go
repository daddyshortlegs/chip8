package chip8

import "fmt"

type Instruction struct {
	instr      uint16
	opCode     byte
	vx         byte
	vy         byte
	opCode2    byte
	secondByte byte
	address    uint16
	vm         *VM
}

type Opcode struct {
	name     string
	function instruction
}

type instruction func()

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

func NewInstruction(instr uint16, vm *VM) *Instruction {
	i := new(Instruction)
	i.vm = vm
	i.extractNibbles(instr)
	return i
}

func (i *Instruction) extractNibbles(instr uint16) {
	i.opCode = extractNibble(instr)
	i.vx = getRightNibble(extractFirstByte(instr))
	i.secondByte = extractSecondByte(instr)
	i.vy = getLeftNibble(i.secondByte)
	i.opCode2 = getRightNibble(i.secondByte)
	i.address = extract12BitNumber(instr)
}

func (i *Instruction) execute() {
	name := i.getOpcodeName()
	if name != "" {
		fmt.Printf("> %s\n", name)
	}

	function := i.getInstructionFromOpcode()
	if function == nil {
		fmt.Printf("Name: %s, unknown instruction %x", name, i.opCode)
	} else {
		function()
	}
}

func (i *Instruction) getInstructionFromOpcode() instruction {
	opCodeFunctions := i.primaryOpcodes()
	return opCodeFunctions[i.opCode].function
}

func (i *Instruction) getOpcodeName() string {
	opCodeFunctions := i.primaryOpcodes()
	return opCodeFunctions[i.opCode].name
}

func (i *Instruction) opDisplay() {
	heightInPixels := i.opCode2

	i.vm.xCoord = i.vm.registers[i.vx] & 63
	i.vm.yCoord = i.vm.registers[i.vy] & 31
	i.vm.registers[15] = 0

	fmt.Printf("Draw index %X, xreg: %d, yreg: %d, x: %d, y: %d, numBytes: %d\n", i.vm.indexRegister, i.vx, i.vy, i.vm.xCoord, i.vm.yCoord, heightInPixels)
	overflow := i.vm.display.DrawSprite(i.vm.indexRegister, heightInPixels, i.vm.xCoord, i.vm.yCoord, i.vm.Memory)
	if overflow == true {
		i.vm.registers[0x0F] = 1
	}
}

func (i *Instruction) opRandom() {
	randomNumber := i.vm.random.Generate()
	i.vm.registers[i.vx] = randomNumber & i.secondByte
}

func (i *Instruction) jumpWithOffset() {
	i.vm.pc = uint16(i.vm.registers[0]) + i.address
	i.vm.pcIncrementer = 0
	fmt.Printf("Jump with offset to %X\n", i.vm.pc)
}

func (i *Instruction) setIndexRegister() {
	i.vm.indexRegister = i.address
}

func (i *Instruction) skipIfRegistersNotEqual() {
	if i.vm.registers[i.vx] != i.vm.registers[i.vy] {
		i.vm.pc += 2
	}
}

func (i *Instruction) setRegister() {
	//fmt.Printf("SetRegister %d to %d\n", index, secondByte)
	i.vm.registers[i.vx] = i.secondByte
}

func (i *Instruction) addToRegister() {
	fmt.Printf("Add To Register [%d] value %d\n", i.vx, i.secondByte)
	i.vm.registers[i.vx] += i.secondByte
	fmt.Printf("&&&&&& VX = %x\n", i.vm.registers[i.vx])
}

func (i *Instruction) skipIfRegistersEqual() {
	if i.vm.registers[i.vx] == i.vm.registers[i.vy] {
		i.vm.pc += 2
	}
}

func (i *Instruction) skipIfNotEqual() {
	if i.vm.registers[i.vx] != i.secondByte {
		i.vm.pc += 2
	}
}

func (i *Instruction) skipIfEqual() {
	if i.vm.registers[i.vx] == i.secondByte {
		i.vm.pc += 2
	}
}

func (i *Instruction) subroutine() {
	i.vm.theStack.Push(i.vm.pc)
	i.vm.pc = i.address
	//i.vm.pcIncrementer = 0
}

// TODO: Need some more tests around jump, as commenting out the pcIncrementer line causes a ROM to work
// AND the Trip8 demo also works properly WITHOUT it.
func (i *Instruction) jump() {
	i.vm.pc = i.address
	//i.vm.pcIncrementer = 0
}

func (i *Instruction) opReturn() {
	address, _ := i.vm.theStack.Pop()
	i.vm.pc = address
	fmt.Printf("Stack popped %X\n", i.vm.pc)
}

func (i *Instruction) clearScreen() {
	println("ClearScreen")
	i.vm.display.ClearScreen()
}

func (i *Instruction) executeArithmeticInstructions() {
	opCodeFunctions := i.arithmeticOpcodes()

	fmt.Printf(">>> %s vx=%d vy=%d\n", opCodeFunctions[i.opCode2].name, i.vx, i.vy)

	opCodeFunctions[i.opCode2].function()
}

func (i *Instruction) setVxToVy() {
	i.vm.registers[i.vx] = i.vm.registers[i.vy]
}

func (i *Instruction) or() {
	i.vm.registers[i.vx] = i.vm.registers[i.vx] | i.vm.registers[i.vy]
}

func (i *Instruction) and() {
	i.vm.registers[i.vx] = i.vm.registers[i.vx] & i.vm.registers[i.vy]
}

func (i *Instruction) xOr() {
	i.vm.registers[i.vx] = i.vm.registers[i.vx] ^ i.vm.registers[i.vy]
}

func (i *Instruction) addToVx() {
	vxRegister := i.vm.registers[i.vx]
	vyRegister := i.vm.registers[i.vy]

	i.vm.registers[i.vx] = vxRegister + vyRegister

	fmt.Printf("******** vx = %d\n", i.vm.registers[i.vx])

	var sum = uint16(vxRegister) + uint16(vyRegister)
	if sum > 255 {
		i.vm.registers[15] = 1
	} else {
		i.vm.registers[15] = 0
	}
}

func (i *Instruction) subtractFromVx() {
	vxRegister := i.vm.registers[i.vx]
	vyRegister := i.vm.registers[i.vy]
	i.vm.registers[i.vx] = vxRegister - vyRegister
	fmt.Printf("******** vx = %d\n", i.vm.registers[i.vx])
	var underflowFlag byte = 1
	if vxRegister < vyRegister {
		underflowFlag = 0
	}
	i.vm.registers[15] = underflowFlag
}

func (i *Instruction) shiftRight() {
	overflow := i.vm.registers[i.vy] & 0b00000001
	i.vm.registers[15] = overflow
	i.vm.registers[i.vx] = i.vm.registers[i.vy] >> 1
}

func (i *Instruction) shiftRightNew() {
	i.vm.registers[i.vx] = i.vm.registers[i.vy]

	overflow := i.vm.registers[i.vx] & 0b00000001
	i.vm.registers[15] = overflow
	i.vm.registers[i.vx] = i.vm.registers[i.vx] >> 1
}

func (i *Instruction) subtractFromVy() {
	vxRegister := i.vm.registers[i.vx]
	vyRegister := i.vm.registers[i.vy]
	i.vm.registers[i.vx] = vyRegister - vxRegister

	var underflowFlag byte = 1
	if vyRegister < vxRegister {
		underflowFlag = 0
	}
	i.vm.registers[15] = underflowFlag
}

func (i *Instruction) shiftLeft() {
	overflow := i.vm.registers[i.vy] & 0b10000000
	i.vm.registers[15] = overflow >> 7
	i.vm.registers[i.vx] = i.vm.registers[i.vy] << 1
}

func (i *Instruction) shiftLeftNew() {
	i.vm.registers[i.vx] = i.vm.registers[i.vy]

	overflow := i.vm.registers[i.vy] & 0b10000000
	i.vm.registers[15] = overflow >> 7
	i.vm.registers[i.vx] = i.vm.registers[i.vx] << 1
}

func (i *Instruction) furtherOperations() {
	opcodes := i.furtherOpcodes()

	functionName := opcodes[i.secondByte].name
	fmt.Printf(">>> %s %x\n", functionName, i.secondByte)

	f := opcodes[i.secondByte].function
	f()
}

func (i *Instruction) bcd() {
	value := i.vm.registers[i.vx]
	hundreds, tens, ones := splitNumberIntoUnits(value)

	address := i.vm.indexRegister
	i.vm.Memory[address] = hundreds
	i.vm.Memory[address+1] = tens
	i.vm.Memory[address+2] = ones

	fmt.Printf("%d %d %d", hundreds, tens, ones)
}

func (i *Instruction) fontChar() {
	println("*** i.vx = ", i.vx)
	character := i.vm.registers[i.vx]
	println("** character = ", character)
	i.vm.indexRegister = 0x50 + uint16(character)*5
	fmt.Printf("** index = %x", i.vm.indexRegister)
}

func (i *Instruction) getKey() {
	// If we get a key then suspend processing of further instruction
	i.vm.processInstructions = false
	key := i.vm.display.GetKey()
	println("****** getKey = ", key)
	i.vm.registers[i.vx] = byte(key)
}

func (i *Instruction) addToIndex() {
	i.vm.indexRegister += uint16(i.vm.registers[i.vx])
}

func (i *Instruction) store() {
	max := int(i.vx)
	startMemory := i.vm.indexRegister
	for n := 0; n <= max; n++ {
		i.vm.Memory[startMemory] = i.vm.registers[n]
		startMemory++
	}
}

func (i *Instruction) load() {
	startMemory := i.vm.indexRegister
	for n := 0; n <= int(i.vx); n++ {
		i.vm.registers[n] = i.vm.Memory[startMemory]
		startMemory++
	}
}

func (i *Instruction) getDelayTimer() {
	// TODO: Test
	// FX07 sets VX to value of the delay timer
	i.vm.registers[i.vx] = i.vm.delayTimer.timer
}

func (i *Instruction) setDelayTimer() {
	// TODO: Test
	// FX15 set the delay timer to value in VX
	i.vm.setDelayTimer(i.vm.registers[i.vx])
}

func (i *Instruction) setSoundTimer() {
	// TODO: Test
	// FX18 sets sound timer to value in VX
	println("vx = ", i.vx)
}

func (i *Instruction) primaryOpcodes() map[byte]Opcode {
	return map[byte]Opcode{
		ClearScreen:             {name: "ClearScreen", function: i.clearScreen},
		Return:                  {name: "Return", function: i.opReturn},
		Jump:                    {name: "Jump", function: i.jump},
		Subroutine:              {name: "Subroutine", function: i.subroutine},
		SkipIfEqual:             {name: "SkipIfEqual", function: i.skipIfEqual},
		SkipIfNotEqual:          {name: "SkipIfNotEqual", function: i.skipIfNotEqual},
		SkipIfRegistersEqual:    {name: "SkipIfRegistersEqual", function: i.skipIfRegistersEqual},
		SkipIfRegistersNotEqual: {name: "SkipIfRegistersNotEqual", function: i.skipIfRegistersNotEqual},
		SetRegister:             {name: "SetRegister", function: i.setRegister},
		AddToRegister:           {name: "AddToRegister", function: i.addToRegister},
		SetIndexRegister:        {name: "SetIndexRegister", function: i.setIndexRegister},
		JumpWithOffset:          {name: "JumpWithOffset", function: i.jumpWithOffset},
		OpRandom:                {name: "OpRandom", function: i.opRandom},
		Display:                 {name: "Display", function: i.opDisplay},
		BitwiseOperations:       {name: "BitwiseOperations", function: i.executeArithmeticInstructions},
		FurtherOperations:       {function: i.furtherOperations},
	}
}

func (i *Instruction) arithmeticOpcodes() map[byte]Opcode {
	return map[byte]Opcode{
		setVxToVy:      {name: "SetVxToBy", function: i.setVxToVy},
		binaryOr:       {name: "Or", function: i.or},
		binaryAnd:      {name: "And", function: i.and},
		logicalXor:     {name: "Xor", function: i.xOr},
		addToVx:        {name: "AddToVx", function: i.addToVx},
		subtractFromVx: {name: "SubtractFromVx", function: i.subtractFromVx},
		shiftRight:     {name: "ShiftRight", function: i.shiftRight},
		subtractFromVy: {name: "SubtractFromVy", function: i.subtractFromVy},
		shiftLeft:      {name: "ShiftLeft", function: i.shiftLeft},
	}
}

func (i *Instruction) furtherOpcodes() map[byte]Opcode {
	return map[byte]Opcode{
		Bcd:           {name: "Bcd", function: i.bcd},
		FontChar:      {name: "FontChar", function: i.fontChar},
		GetKey:        {name: "GetKey", function: i.getKey},
		AddToIndex:    {name: "AddToIndex", function: i.addToIndex},
		Store:         {name: "Store", function: i.store},
		Load:          {name: "Load", function: i.load},
		GetDelayTimer: {name: "GetDelayTimer", function: i.getDelayTimer},
		SetDelayTimer: {name: "SetDelayTimer", function: i.setDelayTimer},
		SetSoundTimer: {name: "SetSoundTimer", function: i.setSoundTimer},
	}
}
