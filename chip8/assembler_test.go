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
	theAssembler.SetRegister(0, 0x00)
	suite.Equal([]byte{0x60, 0x00}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSetRegister1() {
	theAssembler := NewAssembler()
	theAssembler.SetRegister(1, 0x10)
	suite.Equal([]byte{0x61, 0x10}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSetMultipleRegisters() {
	theAssembler := NewAssembler()
	theAssembler.SetRegister(1, 0x10)
	theAssembler.SetRegister(2, 0x22)
	suite.Equal([]byte{0x61, 0x10, 0x62, 0x22}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestIndexRegister() {
	theAssembler := NewAssembler()
	theAssembler.SetIndexRegister(0x123)
	suite.Equal([]byte{0xA1, 0x23}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestAddToRegister() {
	theAssembler := NewAssembler()
	theAssembler.AddToRegister(0, 0x33)
	suite.Equal([]byte{0x70, 0x33}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestDraw() {
	theAssembler := NewAssembler()
	theAssembler.Display(6, 4, 5)
	suite.Equal([]byte{0xD6, 0x45}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestGetKey() {
	theAssembler := NewAssembler()
	theAssembler.GetKey(5)
	suite.Equal([]byte{0xF5, 0x0A}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestClearScreen() {
	theAssembler := NewAssembler()
	theAssembler.ClearScreen()
	suite.Equal([]byte{0x00, 0xE0}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestJump() {
	theAssembler := NewAssembler()
	theAssembler.Jump(0x300)
	suite.Equal([]byte{0x13, 0x00}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestProgram() {
	theAssembler := NewAssembler()
	theAssembler.SetRegister(0, 0x19)
	theAssembler.SetRegister(1, 0x00)

	theAssembler.SetRegister(2, 0x20)
	theAssembler.SetRegister(3, 0x00)

	theAssembler.SetRegister(4, 0x1E)
	theAssembler.SetRegister(5, 0x00)

	theAssembler.SetRegister(6, 0x23)
	theAssembler.SetRegister(7, 0x00)

	theAssembler.SetRegister(8, 0x28)
	theAssembler.SetRegister(9, 0x00)

	theAssembler.SetIndexRegister(0x50)
	theAssembler.Display(0, 1, 5)

	theAssembler.SetIndexRegister(0x55)
	theAssembler.Display(2, 3, 5)

	theAssembler.SetIndexRegister(0x5A)
	theAssembler.Display(4, 5, 5)

	theAssembler.SetIndexRegister(0x5F)
	theAssembler.Display(6, 7, 5)

	theAssembler.SetIndexRegister(0x64)
	theAssembler.Display(8, 9, 5)

	theAssembler.GetKey(3)

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
	suite.Equal(andysProgram, theAssembler.Assemble())
}

func TestAssemblerTestSuite(t *testing.T) {
	suite.Run(t, new(AssemblerTestSuite))
}

/*

 */
