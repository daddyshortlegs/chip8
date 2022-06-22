package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClearScreen(t *testing.T) {
	assert.Equal(t, ClearScreen, decodeInstruction(0x00), "Clear screen")
}

func TestDecodeJump(t *testing.T) {
	assert.Equal(t, Jump, decodeInstruction(0x1A), "Jump")
}

func TestDecodeSetRegister(t *testing.T) {
	assert.Equal(t, SetRegister, decodeInstruction(0x60), "Set Register")
}

func TestDecodeAddValueToRegister(t *testing.T) {
	assert.Equal(t, AddValueToRegister, decodeInstruction(0x70), "Add value to register")
}

func TestDecodeAddValueToRegister2(t *testing.T) {
	assert.Equal(t, AddValueToRegister, decodeInstruction(0x77), "Add value to register")
}

func TestDecodeSetIndexRegisterI(t *testing.T) {
	assert.Equal(t, SetIndexRegister, decodeInstruction(0xA5), "Set index register")
}

func TestDecodeDisplayDraw(t *testing.T) {
	assert.Equal(t, DisplayDraw, decodeInstruction(0xD1), "Display draw")
}
