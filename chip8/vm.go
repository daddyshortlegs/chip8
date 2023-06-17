package chip8

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
