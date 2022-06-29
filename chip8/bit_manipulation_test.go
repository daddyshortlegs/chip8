package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractNibble(t *testing.T) {
	address := uint16(0x12FF)
	assert.Equal(t, byte(0x10), extractNibble(address), "")
}
func TestExtractNibble2(t *testing.T) {
	address := uint16(0x700A)
	assert.Equal(t, byte(0x70), extractNibble(address), "")
}

func TestExtractExtractFirstByte(t *testing.T) {
	address := uint16(0x12FF)
	assert.Equal(t, byte(0x12), extractFirstByte(address), "")
}

func TestExtractExtractSecondByte(t *testing.T) {
	address := uint16(0x12FF)
	assert.Equal(t, byte(0xFF), extractSecondByte(address), "")
}

func TestBytesToWord(t *testing.T) {
	assert.Equal(t, uint16(0x12FF), bytesToWord(0x12, 0xFF), "")
}

func TestExtract12BitNumber(t *testing.T) {
	assert.Equal(t, uint16(0x2FF), extract12BitNumber(0x12FF), "")
}

func TestValueAtPosition0(t *testing.T) {
	assert.Equal(t, uint8(0), getValueAtPosition(0, 0b0000))
	assert.Equal(t, uint8(1), getValueAtPosition(0, 0b0001))
	assert.Equal(t, uint8(0), getValueAtPosition(0, 0b0010))
}

func TestValueAtPosition1(t *testing.T) {
	assert.Equal(t, uint8(1), getValueAtPosition(1, 0b0010))
	assert.Equal(t, uint8(0), getValueAtPosition(1, 0b0001))
	assert.Equal(t, uint8(0), getValueAtPosition(1, 0b1101))
	assert.Equal(t, uint8(1), getValueAtPosition(1, 0b1111))
}

func TestValueAtPosition2(t *testing.T) {
	assert.Equal(t, uint8(0), getValueAtPosition(2, 0b0000))
	assert.Equal(t, uint8(0), getValueAtPosition(2, 0b0001))
	assert.Equal(t, uint8(1), getValueAtPosition(2, 0b0100))
	assert.Equal(t, uint8(1), getValueAtPosition(2, 0b0110))
	assert.Equal(t, uint8(0), getValueAtPosition(2, 0b0010))
	assert.Equal(t, uint8(1), getValueAtPosition(2, 0b1110))
}

func TestValueAtPosition3(t *testing.T) {
	assert.Equal(t, uint8(1), getValueAtPosition(3, 0b1000))
	assert.Equal(t, uint8(0), getValueAtPosition(3, 0b0000))
	assert.Equal(t, uint8(0), getValueAtPosition(3, 0b0001))
	assert.Equal(t, uint8(1), getValueAtPosition(3, 0b1111))
}
