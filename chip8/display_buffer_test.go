package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DisplayBufferTestSuite struct {
	suite.Suite
}

func (suite *DisplayBufferTestSuite) TestEmptyDisplayBuffer() {
	displayBuffer := NewDisplayBuffer()
	pixel := displayBuffer.GetPixelAt(0, 0)
	suite.Equal(uint8(0), pixel)
}

func (suite *DisplayBufferTestSuite) TestDrawByte() {
	displayBuffer := NewDisplayBuffer()
	displayBuffer.drawByte(uint8(0b10000000), 0, 0)

	suite.Equal(uint8(1), displayBuffer.GetPixelAt(0, 0))
	suite.Equal(uint8(0), displayBuffer.GetPixelAt(1, 0))
}

func (suite *DisplayBufferTestSuite) TestWholeScreenPainted() {
	displayBuffer := NewDisplayBuffer()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			displayBuffer.drawByte(uint8(0b10000000), byte(x), byte(y))
		}
	}

	suite.Equal(false, verifyAllBlank(displayBuffer))
}

func (suite *DisplayBufferTestSuite) TestSpriteIsClipped_WhenEdgeOfScreen() {
	displayBuffer := NewDisplayBuffer()

	memory := [4096]byte{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00}

	displayBuffer.DrawSprite(0, 1, 63, 31, memory)

	suite.Equal(false, verifyAllBlank(displayBuffer))
}

func (suite *DisplayBufferTestSuite) TestScreenIsClipped() {
	displayBuffer := NewDisplayBuffer()

	for y := 0; y < 64; y++ {
		for x := 0; x < 128; x++ {
			displayBuffer.drawByte(uint8(0b10000000), byte(x), byte(y))
		}
	}

	suite.Equal(false, verifyAllBlank(displayBuffer))
}

func (suite *DisplayBufferTestSuite) TestClearScreen() {
	displayBuffer := NewDisplayBuffer()
	displayBuffer.drawByte(uint8(0b11111111), 0, 0)

	displayBuffer.ClearScreen()

	suite.Equal(uint8(0), displayBuffer.GetPixelAt(0, 0))
	suite.Equal(uint8(0), displayBuffer.GetPixelAt(1, 0))
}

func (suite *DisplayBufferTestSuite) TestWholeScreenCleared() {
	displayBuffer := NewDisplayBuffer()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			displayBuffer.drawByte(uint8(0b10000000), byte(x), byte(y))
		}
	}

	displayBuffer.ClearScreen()

	suite.Equal(true, verifyAllBlank(displayBuffer))
}

func verifyAllBlank(displayBuffer *DisplayBuffer) bool {
	allBlanks := true
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			pixel := displayBuffer.GetPixelAt(byte(x), byte(y))
			if pixel == 1 {
				allBlanks = false
			}
		}
	}
	return allBlanks
}

func (suite *DisplayBufferTestSuite) TestXor() {
	displayBuffer := NewDisplayBuffer()
	displayBuffer.drawByte(uint8(0b11111111), 0, 0)
	displayBuffer.drawByte(uint8(0b11111111), 0, 0)

	suite.Equal(uint8(0), displayBuffer.GetPixelAt(0, 0))
}

func (suite *DisplayBufferTestSuite) TestXorPattern() {
	displayBuffer := NewDisplayBuffer()
	displayBuffer.drawByte(uint8(0b10101010), 0, 0)
	displayBuffer.drawByte(uint8(0b11111111), 0, 0)

	suite.Equal(uint8(0), displayBuffer.GetPixelAt(0, 0))
	suite.Equal(uint8(1), displayBuffer.GetPixelAt(1, 0))
	suite.Equal(uint8(0), displayBuffer.GetPixelAt(2, 0))
	suite.Equal(uint8(1), displayBuffer.GetPixelAt(3, 0))
	suite.Equal(uint8(0), displayBuffer.GetPixelAt(4, 0))
	suite.Equal(uint8(1), displayBuffer.GetPixelAt(5, 0))
	suite.Equal(uint8(0), displayBuffer.GetPixelAt(6, 0))
	suite.Equal(uint8(1), displayBuffer.GetPixelAt(7, 0))
}

func TestDisplayBufferSuite(t *testing.T) {
	suite.Run(t, new(DisplayBufferTestSuite))
}
