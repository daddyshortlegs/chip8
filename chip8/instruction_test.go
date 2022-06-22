package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClearScreen(t *testing.T) {
	assert.Equal(t, ClearScreen, decodeInstruction(0x00, 0), "Clear screen")
}

func TestDecodeJump(t *testing.T) {
	assert.Equal(t, Jump, decodeInstruction(0x1A, 0), "Jump")
}

func TestDecodeSetRegister(t *testing.T) {
	assert.Equal(t, SetRegister, decodeInstruction(0x60, 0), "Set Register")
}

func TestDecodeAddValueToRegister(t *testing.T) {
	assert.Equal(t, AddValueToRegister, decodeInstruction(0x70, 0), "Add value to register")
}

func TestDecodeAddValueToRegister2(t *testing.T) {
	assert.Equal(t, AddValueToRegister, decodeInstruction(0x77, 0), "Add value to register")
}

func TestDecodeSetIndexRegisterI(t *testing.T) {
	assert.Equal(t, SetIndexRegister, decodeInstruction(0xA5, 0), "Set index register")
}

func TestDecodeDisplayDraw(t *testing.T) {
	assert.Equal(t, DisplayDraw, decodeInstruction(0xD1, 0), "Display draw")
}
