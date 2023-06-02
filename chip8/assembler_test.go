package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssemblerTestSuite struct {
	suite.Suite
}

func (suite *AssemblerTestSuite) TestSetRegister0() {
	theAssembler := NewAssembler()
	theAssembler.setRegister(0, 0x00)
	suite.Equal([]byte{0x60, 0x00}, theAssembler.assemble())
}

func (suite *AssemblerTestSuite) TestSetRegister1() {
	theAssembler := NewAssembler()
	theAssembler.setRegister(1, 0x10)
	suite.Equal([]byte{0x61, 0x10}, theAssembler.assemble())
}

func (suite *AssemblerTestSuite) TestSetMultipleRegisters() {
	theAssembler := NewAssembler()
	theAssembler.setRegister(1, 0x10)
	theAssembler.setRegister(2, 0x22)
	suite.Equal([]byte{0x61, 0x10, 0x62, 0x22}, theAssembler.assemble())
}

func (suite *AssemblerTestSuite) TestIndexRegister() {
	theAssembler := NewAssembler()
	theAssembler.setIndexRegister(0x123)
	suite.Equal([]byte{0xA1, 0x23}, theAssembler.assemble())
}

func (suite *AssemblerTestSuite) TestDraw() {
	theAssembler := NewAssembler()
	theAssembler.display(6, 4, 5)
	suite.Equal([]byte{0xD6, 0x45}, theAssembler.assemble())
}

func (suite *AssemblerTestSuite) TestGetKey() {
	theAssembler := NewAssembler()
	theAssembler.getKey(5)
	suite.Equal([]byte{0xF5, 0x0A}, theAssembler.assemble())
}

func (suite *AssemblerTestSuite) TestProgram() {
	theAssembler := NewAssembler()
	theAssembler.setRegister(0, 0x19)
	theAssembler.setRegister(1, 0x00)

	theAssembler.setRegister(2, 0x20)
	theAssembler.setRegister(3, 0x00)

	theAssembler.setRegister(4, 0x1E)
	theAssembler.setRegister(5, 0x00)

	theAssembler.setRegister(6, 0x23)
	theAssembler.setRegister(7, 0x00)

	theAssembler.setRegister(8, 0x28)
	theAssembler.setRegister(9, 0x00)

	theAssembler.setIndexRegister(0x50)
	theAssembler.display(0, 1, 5)

	theAssembler.setIndexRegister(0x55)
	theAssembler.display(2, 3, 5)

	theAssembler.setIndexRegister(0x5A)
	theAssembler.display(4, 5, 5)

	theAssembler.setIndexRegister(0x5F)
	theAssembler.display(6, 7, 5)

	theAssembler.setIndexRegister(0x64)
	theAssembler.display(8, 9, 5)

	theAssembler.getKey(3)

	andysProgram := []byte{
		0x60, 0x19,
		0x61, 0x00,

		0x62, 0x20,
		0x63, 0x00,

		0x64, 0x1E,
		0x65, 0x00,

		0x66, 0x23,
		0x67, 0x00,

		0x68, 0x28,
		0x69, 0x00,

		0xA0, 0x50,
		0xD0, 0x15,

		0xA0, 0x55,
		0xD2, 0x35,

		0xA0, 0x5A,
		0xD4, 0x55,

		0xA0, 0x5F,
		0xD6, 0x75,

		0xA0, 0x64,
		0xD8, 0x95,

		0xF3, 0x0A,
	}
	suite.Equal(andysProgram, theAssembler.assemble())
}

func TestAssemblerTestSuite(t *testing.T) {
	suite.Run(t, new(AssemblerTestSuite))
}

/*

 */
