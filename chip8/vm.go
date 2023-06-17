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

type opcodes func(uint16, *VM)
type furtherOpcodes func(byte)

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

	instr := v.fetchAndIncrement()
	if instr == 0x0000 {
		return true
	}

	v.pcIncrementer = 2
	i := NewInstruction(instr)

	if instr == ClearScreen {
		i.clearScreen(instr, v)
	} else if instr == Return {
		i.opReturn(instr, v)
	} else {
		i.execute(instr, v)
	}
	return false
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

	m[extractSecondByte(instr)](i.vx)
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

func (v *VM) fetch() uint16 {
	return bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
}

func (v *VM) fetchAndIncrement() uint16 {
	i := bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
	v.pc += uint16(v.pcIncrementer)
	return i
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
