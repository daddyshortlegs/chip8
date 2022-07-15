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

func (suite *Chip8TestSuite) TestClearScreen() {
	m := mockDisplay{false, drawPatternValues{}}
	var display Display
	display = &m

	suite.vm.SetDisplay(display)
	suite.vm.Load([]byte{0x00, 0xE0})
	suite.vm.Run()

	suite.True(m.screenCleared)
}

func (suite *Chip8TestSuite) TestGetCoordinatesFromRegisters_whenDraw() {
	m := mockDisplay{false, drawPatternValues{}}
	var display Display
	display = &m
	suite.vm.SetDisplay(display)

	suite.vm.registers[5] = 20
	suite.vm.registers[10] = 30
	suite.vm.Load([]byte{0xD5, 0xA0})

	suite.vm.Run()

	suite.Equal(byte(20), suite.vm.getXCoordinate())
	suite.Equal(byte(30), suite.vm.getYCoordinate())
}

func (suite *Chip8TestSuite) TestCoordinatesShouldWrap() {
	m := mockDisplay{false, drawPatternValues{}}
	var display Display
	display = &m
	suite.vm.SetDisplay(display)

	suite.vm.registers[5] = 64
	suite.vm.registers[10] = 32
	suite.vm.indexRegister = 0x200
	suite.vm.Load([]byte{0xD5, 0xA5})

	suite.vm.Run()

	suite.Equal(byte(0), suite.vm.getXCoordinate())
	suite.Equal(byte(0), suite.vm.getYCoordinate())
}

func (suite *Chip8TestSuite) TestInitialMemoryContainsFont() {
	bytes := suite.vm.Memory[0x50:0x09F]
	suite.Equal(byte(0xF0), bytes[0], "First byte")
	suite.Equal(byte(0x90), bytes[1], "Second byte")
	suite.Equal(byte(0x80), bytes[len(bytes)-1], "last byte")
}

func (suite *Chip8TestSuite) TestLoadPlacesCodeAtCorrectPlace() {
	suite.vm.Load([]byte{0xD5, 0xA0})

	bytes := suite.vm.Memory[0x200:]
	suite.Equal(byte(0xD5), bytes[0], "First byte")
	suite.Equal(byte(0xA0), bytes[1], "Second byte")
}

func (suite *Chip8TestSuite) TestDraw() {
	m := mockDisplay{false, drawPatternValues{}}
	var display Display
	display = &m
	suite.vm.SetDisplay(display)

	instructions1 := setRegisterOpcode(0x5, 0x14)
	instructions2 := setRegisterOpcode(0xA, 0x1E)
	indexInstruction := setIndexRegisterOpcode(0x050)
	drawInstruction := drawOpcode(0x5, 0xA, 5)

	result := append(instructions1, instructions2...)
	result = append(result, indexInstruction...)
	result = append(result, drawInstruction...)
	suite.vm.Load(result)

	suite.vm.Run()

	suite.Equal(byte(20), m.values.x)
	suite.Equal(byte(30), m.values.y)
	suite.Equal(uint16(0x50), m.values.address)
	suite.Equal(byte(5), m.values.numberOfBytes)
}

func (suite *Chip8TestSuite) TestVXIsSetToVY() {
	suite.executeInstruction([]byte{
		0x65, 0x14, // Set register 5 to 0x14 (20)
		0x80, 0x50, // Set register 0 to what's in register 5
	})

	suite.Equal(byte(0x14), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryORofVXVY() {
	suite.executeInstruction([]byte{
		0x60, 0x0F, // Set register 0 to 0x0F
		0x61, 0xF0, // Set register 1 to 0xF0
		0x80, 0x11, // Set register 0 to what's in register 0 & 1 ORd together
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryANDofVXVY() {
	suite.executeInstruction([]byte{
		0x60, 0b00001111, // Set register 0 to ...
		0x61, 0b00110011, // Set register 1 to ...
		0x80, 0x12, // Set register 0 to what's in register 0 & 1 ANDd together
	})

	suite.Equal(byte(0b00000011), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryXORofVXVY() {
	suite.executeInstruction([]byte{
		0x60, 0b00001111, // Set register 0 to ...
		0x61, 0b00110011, // Set register 1 to ...
		0x80, 0x13, // Set register 0 to what's in register 0 & 1 ANDd together
	})

	suite.Equal(byte(0b00111100), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestAddWithNoCarry() {
	suite.executeInstruction([]byte{
		0x60, 0x0A, // Set register 0 to ...
		0x61, 0x0A, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestAddWithCarry() {
	suite.executeInstruction([]byte{
		0x60, 0xFF, // Set register 0 to ...
		0x61, 0x01, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})

	suite.Equal(byte(0x00), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestCarryFlagIsSetTo0AfterPreviousCarry() {
	suite.vm = Chip8vm{}
	suite.vm.Init()

	suite.vm.registers[15] = 1

	suite.vm.Load([]byte{
		0x60, 0x0A, // Set register 0 to ...
		0x61, 0x0A, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})
	suite.vm.Run()

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVY() {
	suite.executeInstruction([]byte{
		0x60, 0x0A, // Set register 0 to 10
		0x61, 0x01, // Set register 1 to 1
		0x80, 0x15, // Set VX to 10 - 1
	})

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVYUnderflow() {
	suite.executeInstruction([]byte{
		0x60, 0x0A, // Set register 0 to 10
		0x61, 0x0B, // Set register 1 to 1
		0x80, 0x15, // Set VX to 10 - 11
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVX() {
	suite.executeInstruction([]byte{
		0x60, 0x01, // Set register 0 to 1
		0x61, 0x0A, // Set register 1 to 10
		0x80, 0x17, // Set VX to 10 - 1
	})

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVXUnderflow() {
	suite.executeInstruction([]byte{
		0x60, 0x0B, // Set register 0 to 11
		0x61, 0x0A, // Set register 1 to 10
		0x80, 0x17, // Set VX to 10 - 1
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

//8XY6
func (suite *Chip8TestSuite) TestVXShiftRight() {
	suite.executeInstruction([]byte{
		0x61, 0b11111110, // Set register 1 to 10
		0x80, 0x16, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b11111110), suite.vm.registers[1])
	suite.Equal(byte(0b01111111), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

//8XY6
func (suite *Chip8TestSuite) TestVXShiftRightWithOverflow() {
	suite.executeInstruction([]byte{
		0x61, 0b00110001, // Set register 1 to 10
		0x80, 0x16, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b00110001), suite.vm.registers[1])
	suite.Equal(byte(0b00011000), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXShiftLeft() {
	suite.executeInstruction([]byte{
		0x61, 0b01111110, // Set register 1 to 10
		0x80, 0x1E, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b01111110), suite.vm.registers[1])
	suite.Equal(byte(0b11111100), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXShiftLeftWithOverflow() {
	suite.executeInstruction([]byte{
		0x61, 0b11111100, // Set register 1 to 10
		0x80, 0x1E, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b11111100), suite.vm.registers[1])
	suite.Equal(byte(0b11111000), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
