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

func (suite *Chip8TestSuite) TestSetRegister() {
	suite.vm = chip8vm{}
	instruction := []byte{0x60, 0x01}
	suite.vm.load(instruction)
	suite.vm.run()

	suite.Equal(byte(0x01), suite.vm.registers[0], "Register V0")
}
func (suite *Chip8TestSuite) TestSetRegisterToSomefinkElse() {
	suite.vm = chip8vm{}
	instruction := []byte{0x60, 0x02}
	suite.vm.load(instruction)
	suite.vm.run()

	suite.Equal(byte(0x02), suite.vm.registers[0], "Register V0")
}

func (suite *Chip8TestSuite) TestSetRegister2() {
	suite.vm = chip8vm{}
	instruction := []byte{0x61, 0xFF}
	suite.vm.load(instruction)
	suite.vm.run()

	suite.Equal(byte(0xFF), suite.vm.registers[1], "Register V0")
}

func (suite *Chip8TestSuite) TestSetRegister2ToSomethingElse() {
	suite.vm = chip8vm{}
	instruction := []byte{0x61, 0xEE}
	suite.vm.load(instruction)
	suite.vm.run()

	suite.Equal(byte(0xEE), suite.vm.registers[1], "Register V0")
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
