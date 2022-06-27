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
