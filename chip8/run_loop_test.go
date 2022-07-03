package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Chip8TestSuite struct {
	suite.Suite
	vm Chip8vm
}

func (suite *Chip8TestSuite) SetupTest() {
	suite.vm = Chip8vm{}
	suite.vm.Init()
}

func (suite *Chip8TestSuite) TestFetchInstruction() {
	suite.vm.Load([]byte{0x12, 0x20})

	decoded := suite.vm.fetchAndIncrement()
	suite.Equal(uint16(0x1220), decoded, "First byte")
}

func (suite *Chip8TestSuite) TestFetchNextInstruction() {
	suite.vm.Load([]byte{0x12, 0x20, 0x33, 0x44})
	suite.vm.fetchAndIncrement()

	decoded := suite.vm.fetchAndIncrement()

	suite.Equal(uint16(0x3344), decoded, "First byte")
}

func (suite *Chip8TestSuite) TestSetRegisters() {
	verifyRegisterSet(suite, []byte{0x60, 0xFF}, 0, 0xFF)
	verifyRegisterSet(suite, []byte{0x61, 0xEE}, 1, 0xEE)
	verifyRegisterSet(suite, []byte{0x62, 0xDD}, 2, 0xDD)
	verifyRegisterSet(suite, []byte{0x63, 0xCC}, 3, 0xCC)
	verifyRegisterSet(suite, []byte{0x64, 0xBB}, 4, 0xBB)
}

func verifyRegisterSet(suite *Chip8TestSuite, instruction []byte, register int, result int) {
	suite.executeInstruction(instruction)
	suite.Equal(byte(result), suite.vm.registers[register])
}

func (suite *Chip8TestSuite) TestFetchAndSetAllRegisters() {
	suite.vm.Load([]byte{0x60, 0x11, 0x61, 0x12, 0x65, 0xCC})
	suite.vm.Run()

	suite.Equal(byte(0x11), suite.vm.registers[0])
	suite.Equal(byte(0x12), suite.vm.registers[1])
	suite.Equal(byte(0xCC), suite.vm.registers[5])
}

func (suite *Chip8TestSuite) executeInstruction(data []byte) {
	suite.vm = Chip8vm{}
	suite.vm.Init()
	suite.vm.Load(data)
	suite.vm.Run()
}

func (suite *Chip8TestSuite) TestAddToRegister() {
	suite.executeInstruction([]byte{0x70, 0x0A})
	suite.Equal(byte(0x0A), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetAndAddToRegister() {
	suite.executeInstruction([]byte{0x60, 0x01, 0x70, 0x0A})
	suite.Equal(byte(0x0B), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetIndexRegister() {
	suite.executeInstruction([]byte{0xA0, 0x0A})
	suite.Equal(uint16(0x0A), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestSetIndexRegisterWith12BitValue() {
	suite.executeInstruction([]byte{0xAF, 0xFF})
	suite.Equal(uint16(0xFFF), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestSetJumpToAddress() {
	suite.executeInstruction([]byte{0x13, 0x00})
	suite.Equal(uint16(0x300), suite.vm.pc)
}

type mockDisplay struct {
	screenCleared bool
}

func (m *mockDisplay) DrawPattern(uint16, int) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDisplay) ClearScreen() {
	m.screenCleared = true
}

func (suite *Chip8TestSuite) TestClearScreen() {
	m := mockDisplay{false}
	var display Display
	display = &m

	suite.vm.SetDisplay(display)
	suite.vm.Load([]byte{0x00, 0xE0})
	suite.vm.Run()

	suite.True(m.screenCleared)
}

func (suite *Chip8TestSuite) TestGetCoordinatesFromRegisters_whenDraw() {
	suite.vm.registers[5] = 20
	suite.vm.registers[10] = 30
	suite.vm.Load([]byte{0xD5, 0xA0})

	suite.vm.Run()

	suite.Equal(byte(20), suite.vm.getXCoordinate())
	suite.Equal(byte(30), suite.vm.getYCoordinate())
}

func (suite *Chip8TestSuite) TestCoordinatesShouldWrap() {
	suite.vm.registers[5] = 64
	suite.vm.registers[10] = 32
	suite.vm.Load([]byte{0xD5, 0xA0})

	suite.vm.Run()

	suite.Equal(byte(0), suite.vm.getXCoordinate())
	suite.Equal(byte(0), suite.vm.getYCoordinate())
}

func (suite *Chip8TestSuite) TestInitialMemoryContainsFont() {
	bytes := suite.vm.memory[0x50:0x09F]
	suite.Equal(byte(0xF0), bytes[0], "First byte")
	suite.Equal(byte(0x90), bytes[1], "Second byte")
	suite.Equal(byte(0x80), bytes[len(bytes)-1], "last byte")
}

func (suite *Chip8TestSuite) TestLoadPlacesCodeAtCorrectPlace() {
	suite.vm.Load([]byte{0xD5, 0xA0})

	bytes := suite.vm.memory[0x200:]
	suite.Equal(byte(0xD5), bytes[0], "First byte")
	suite.Equal(byte(0xA0), bytes[1], "Second byte")
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
