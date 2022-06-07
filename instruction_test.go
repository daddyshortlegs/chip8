package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClearScreen(t *testing.T) {
	assert.Equal(t, ClearScreen, decode([]byte{0x00, 0xE0}), "Clear screen")
}

func TestDecodeJump(t *testing.T) {
	assert.Equal(t, Jump, decode([]byte{0x1A, 0x2B}), "Jump")
}

func TestDecodeSetRegister(t *testing.T) {
	assert.Equal(t, SetRegister, decode([]byte{0x60, 0x20}), "Set Register")
}

func TestDecodeAddValueToRegister(t *testing.T) {
	assert.Equal(t, AddValueToRegister, decode([]byte{0x70, 0x20}), "Add value to register")
}

func TestDecodeAddValueToRegister2(t *testing.T) {
	assert.Equal(t, AddValueToRegister, decode([]byte{0x77, 0x20}), "Add value to register")
}

func TestDecodeSetIndexRegisterI(t *testing.T) {
	assert.Equal(t, SetIndexRegister, decode([]byte{0xA5, 0x55}), "Set index register")
}

func TestDecodeDisplayDraw(t *testing.T) {
	assert.Equal(t, DisplayDraw, decode([]byte{0xD1, 0x55}), "Display draw")
}
