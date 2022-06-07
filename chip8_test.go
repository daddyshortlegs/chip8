package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func decode(instruction []byte) int {
	m := memory{}

	m.load(instruction)

	decoded := m.decode()
	return decoded
}

func TestDecodeByte(t *testing.T) {
	instruction := []byte{0x12, 0x20}
	m := memory{}

	m.load(instruction)

	decoded := m.fetch()
	assert.Equal(t, byte(0x12), decoded.first, "First nibble")
	assert.Equal(t, byte(0x20), decoded.second, "Second nibble")
}
