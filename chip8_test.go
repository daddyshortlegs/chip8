package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Chip8TestSuite struct {
	suite.Suite
	m  memory
	vm chip8vm
}

func (suite *Chip8TestSuite) SetupTest() {
	println("**** setup test")
}

func (suite *Chip8TestSuite) TestFetchInstruction() {
	suite.m = memory{}
	instruction := []byte{0x12, 0x20}
	suite.m.load(instruction)

	decoded := suite.m.fetch()
	suite.Equal(byte(0x12), decoded.first, "First byte")
	suite.Equal(byte(0x20), decoded.second, "Second byte")
}

func (suite *Chip8TestSuite) TestFetchNextInstruction() {
	suite.m = memory{}
	instruction := []byte{0x12, 0x20, 0x33, 0x44}
	suite.m.load(instruction)
	suite.m.fetch()

	decoded := suite.m.fetch()

	suite.Equal(byte(0x33), decoded.first, "First byte")
	suite.Equal(byte(0x44), decoded.second, "Second byte")
}

func (suite *Chip8TestSuite) TestSetRegisters() {
	verifyRegisterSet(suite, []byte{0x60, 0xFF}, 0, 0xFF)
	verifyRegisterSet(suite, []byte{0x61, 0xEE}, 1, 0xEE)
	verifyRegisterSet(suite, []byte{0x62, 0xDD}, 2, 0xDD)
	verifyRegisterSet(suite, []byte{0x63, 0xCC}, 3, 0xCC)
	verifyRegisterSet(suite, []byte{0x64, 0xBB}, 4, 0xBB)
}

func verifyRegisterSet(suite *Chip8TestSuite, instruction []byte, register int, result int) {
	suite.vm = chip8vm{}
	suite.vm.load(instruction)
	suite.vm.run()

	suite.Equal(byte(result), suite.vm.registers[register])
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
