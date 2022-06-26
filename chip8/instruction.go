package chip8

type instruction struct {
	first  byte
	second byte
}

type Instruction interface {
	Execute()
}

type ClearScreenInstruction struct {
}

func Execute(instr instruction) {

}
