package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeByte(t *testing.T) {
	instruction := [2]byte {0x12, 0x20}

	decoded := decodeInstruction(instruction)
	assert.Equal(t, byte(0x12), decoded.first, "First nibble")
	assert.Equal(t, byte(0x20), decoded.second, "Second nibble")
}

type instruction struct {
	first byte
	second byte
}


func decodeInstruction(bytes [2]byte) instruction {
	return instruction{bytes[0], bytes[1]}
}