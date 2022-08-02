package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSetRegisterInstruction(t *testing.T) {
	result := setRegisterOpcode(5, 0x14)
	assert.Equal(t, []byte{0x65, 0x14}, result)
}

func TestCreateSetIndexRegisterInstruction(t *testing.T) {
	assert.Equal(t, []byte{0xA0, 0x14}, setIndexRegisterOpcode(0x14))
	assert.Equal(t, []byte{0xA2, 0x14}, setIndexRegisterOpcode(0x214))
}

func TestCreateDrawInstruction(t *testing.T) {
	result := drawOpcode(0x5, 0xA, 5)
	assert.Equal(t, []byte{0xD5, 0xA5}, result)
}

func TestSplitInteger(t *testing.T) {
	hundreds, tens, ones := splitNumber(157)
	assert.Equal(t, byte(1), hundreds)
	assert.Equal(t, byte(5), tens)
	assert.Equal(t, byte(7), ones)
}
