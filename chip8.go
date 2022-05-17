package chip8

type instruction struct {
	first  byte
	second byte
}

func decodeInstruction(bytes [2]byte) instruction {
	return instruction{bytes[0], bytes[1]}
}