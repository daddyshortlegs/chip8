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
	suite.vm = chip8vm{}
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

func (suite *Chip8TestSuite) TestFetchAndSetAllRegisters() {
	suite.vm = chip8vm{}
	data := []byte{0x60, 0x11, 0x61, 0x12, 0x65, 0xCC}
	suite.vm.load(data)
	suite.vm.run()

	suite.Equal(byte(0x11), suite.vm.registers[0])
	suite.Equal(byte(0x12), suite.vm.registers[1])
	suite.Equal(byte(0xCC), suite.vm.registers[5])
}

func (suite *Chip8TestSuite) TestAddToRegister() {
	suite.vm = chip8vm{}
	data := []byte{0x70, 0x0A}
	suite.vm.load(data)
	suite.vm.run()

	suite.Equal(byte(0x0A), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestsSetAndAddToRegister() {
	suite.vm = chip8vm{}
	data := []byte{0x60, 0x01, 0x70, 0x0A}
	suite.vm.load(data)
	suite.vm.run()

	suite.Equal(byte(0x0B), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetIndexRegister() {
	suite.vm = chip8vm{}
	data := []byte{0xA0, 0x0A}
	suite.vm.load(data)
	suite.vm.run()

	suite.Equal(uint16(0x0A), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestSetIndexRegisterWith12BitValue() {
	suite.vm = chip8vm{}
	data := []byte{0xAF, 0xFF}
	suite.vm.load(data)
	suite.vm.run()

	suite.Equal(uint16(0xFFF), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestSetJumpToAddress() {
	suite.vm = chip8vm{}
	data := []byte{0x12, 0x00}
	suite.vm.load(data)
	suite.vm.run()

	suite.Equal(uint16(0x200), suite.vm.pc)
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
