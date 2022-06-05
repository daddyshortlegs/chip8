package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeByte(t *testing.T) {
	m := memory{}
	instruction := []byte{0x12, 0x20}

	m.load(instruction)

	decoded := m.fetch()
	assert.Equal(t, byte(0x12), decoded.first, "First nibble")
	assert.Equal(t, byte(0x20), decoded.second, "Second nibble")
}

func TestDecodeJump(t *testing.T) {
	m := memory{}
	instruction := []byte{0x1A, 0x2B}

	m.load(instruction)

	decoded := m.decode()
	assert.Equal(t, Jump, decoded, "Jump")
}

func TestDecodeSetRegister(t *testing.T) {
	m := memory{}
	instruction := []byte{0x60, 0x20}

	m.load(instruction)

	decoded := m.decode()
	assert.Equal(t, SetRegister, decoded, "Set Register")
}

func TestDecodeAddValueToRegister(t *testing.T) {
	m := memory{}
	instruction := []byte{0x70, 0x20}

	m.load(instruction)

	decoded := m.decode()
	assert.Equal(t, AddValueToRegister, decoded, "Add value to register")
}

func TestDecodeAddValueToRegister2(t *testing.T) {
	m := memory{}
	instruction := []byte{0x77, 0x20}

	m.load(instruction)

	decoded := m.decode()
	assert.Equal(t, AddValueToRegister, decoded, "Add value to register")
}

func TestDecodeSetIndexRegisterI(t *testing.T) {
	m := memory{}
	instruction := []byte{0xA5, 0x55}

	m.load(instruction)

	decoded := m.decode()
	assert.Equal(t, SetIndexRegister, decoded, "Set index register")
}
