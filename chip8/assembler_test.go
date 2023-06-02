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

func TestAssemblerTestSuite(t *testing.T) {
	suite.Run(t, new(AssemblerTestSuite))
}

/*
andysProgram := []byte{
0x60, 0x19, // Set register 0 to 0x00
0x61, 0x00, // Set register 1 to 0x00

0x62, 0x20, // Set register 2 to 0x05
0x63, 0x00, // Set register 3 to 0x00

0x64, 0x1E, // Set register 4 to 10
0x65, 0x00, // Set register 5 to 0x00

0x66, 0x23, // Set register 6 to 15
0x67, 0x00, // Set register 7 to 0x00

0x68, 0x28, // Set register 8 to 20
0x69, 0x00, // Set register 9 to 0x00

0xA0, 0x50, // Set Index Register to 0x50
0xD0, 0x15, // Draw, Xreg = 5, Y reg = 10, 5 bytes high

0xA0, 0x55, // Set Index Register to 0x55
0xD2, 0x35, // Draw, Xreg = 5, Y reg = 10, 5 bytes high

0xA0, 0x5A, // Set Index Register to 0x55
0xD4, 0x55, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
0xA0, 0x5F, // Set Index Register to 0x5F
0xD6, 0x75, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
0xA0, 0x64, // Set Index Register to 0x5F
0xD8, 0x95, // Draw, Xreg = 5, Y reg = 10, 5 bytes high
0xF3, 0x0A, // Wait for key
}
*/
