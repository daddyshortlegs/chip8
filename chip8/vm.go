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

func (v *VM) fetchAndIncrement() uint16 {
	i := bytesToWord(v.Memory[v.pc], v.Memory[v.pc+1])
	v.pc += uint16(v.pcIncrementer)
	return i
}

func (v *VM) getXCoordinate() byte {
	return v.xCoord
}

func (v *VM) getYCoordinate() byte {
	return v.yCoord
}

func (v *VM) setDelayTimer(b byte) {
	v.delayTimer.setTimer(b)
}
