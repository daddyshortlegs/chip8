package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Chip8TestSuite struct {
	suite.Suite
	vm *Chip8VM
	m  mockDisplay
}

const FONT_MEMORY = 0x50
const SET_REGISTER_0 = 0x60
const SET_REGISTER_1 = 0x61
const SET_REGISTER_2 = 0x62
const SET_REGISTER_3 = 0x63
const SET_REGISTER_4 = 0x64
const ADD_REGISTER_0 = 0x70
const FONT_REGISTER_0 = 0xF0

func (suite *Chip8TestSuite) SetupTest() {
	suite.m = mockDisplay{false, drawPatternValues{}, true}
	suite.vm = NewChip8VM(&suite.m)
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
	verifyRegisterSet(suite, []byte{SET_REGISTER_0, 0xFF}, 0, 0xFF)
	verifyRegisterSet(suite, []byte{SET_REGISTER_1, 0xEE}, 1, 0xEE)
	verifyRegisterSet(suite, []byte{SET_REGISTER_2, 0xDD}, 2, 0xDD)
	verifyRegisterSet(suite, []byte{SET_REGISTER_3, 0xCC}, 3, 0xCC)
	verifyRegisterSet(suite, []byte{SET_REGISTER_4, 0xBB}, 4, 0xBB)
}

func verifyRegisterSet(suite *Chip8TestSuite, instruction []byte, register int, result int) {
	suite.executeInstruction(instruction)
	suite.Equal(byte(result), suite.vm.registers[register])
}

func (suite *Chip8TestSuite) TestFetchAndSetAllRegisters() {
	suite.vm.Load([]byte{SET_REGISTER_0, 0x11, SET_REGISTER_1, 0x12, 0x65, 0xCC})
	suite.vm.Run()

	suite.Equal(byte(0x11), suite.vm.registers[0])
	suite.Equal(byte(0x12), suite.vm.registers[1])
	suite.Equal(byte(0xCC), suite.vm.registers[5])
}

func (suite *Chip8TestSuite) executeInstruction(data []byte) {
	m := mockDisplay{false, drawPatternValues{}, true}
	suite.vm = NewChip8VM(&m)

	suite.vm.Load(data)
	suite.vm.Run()
}

func (suite *Chip8TestSuite) TestAddToRegister() {
	suite.executeInstruction([]byte{ADD_REGISTER_0, 0x0A})
	suite.Equal(byte(0x0A), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetAndAddToRegister() {
	suite.executeInstruction([]byte{SET_REGISTER_0, 0x01, ADD_REGISTER_0, 0x0A})
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
	suite.vm.Load([]byte{0x00, 0xE0})
	suite.vm.Run()

	suite.True(suite.m.screenCleared)
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
	instructions1 := setRegisterOpcode(0x5, 0x14)
	instructions2 := setRegisterOpcode(0xA, 0x1E)
	indexInstruction := setIndexRegisterOpcode(0x050)
	drawInstruction := drawOpcode(0x5, 0xA, 5)

	result := append(instructions1, instructions2...)
	result = append(result, indexInstruction...)
	result = append(result, drawInstruction...)
	suite.vm.Load(result)

	suite.vm.Run()

	suite.Equal(byte(20), suite.m.values.x)
	suite.Equal(byte(30), suite.m.values.y)
	suite.Equal(uint16(0x50), suite.m.values.address)
	suite.Equal(byte(5), suite.m.values.numberOfBytes)
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
		SET_REGISTER_0, 0x0F, // Set register 0 to 0x0F
		SET_REGISTER_1, 0xF0, // Set register 1 to 0xF0
		0x80, 0x11, // Set register 0 to what's in register 0 & 1 ORd together
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryANDofVXVY() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0b00001111, // Set register 0 to ...
		SET_REGISTER_1, 0b00110011, // Set register 1 to ...
		0x80, 0x12, // Set register 0 to what's in register 0 & 1 ANDd together
	})

	suite.Equal(byte(0b00000011), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryXORofVXVY() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0b00001111, // Set register 0 to ...
		SET_REGISTER_1, 0b00110011, // Set register 1 to ...
		0x80, 0x13, // Set register 0 to what's in register 0 & 1 ANDd together
	})

	suite.Equal(byte(0b00111100), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestAddWithNoCarry() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x0A, // Set register 0 to ...
		SET_REGISTER_1, 0x0A, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestAddWithCarry() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0xFF, // Set register 0 to ...
		SET_REGISTER_1, 0x01, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})

	suite.Equal(byte(0x00), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestCarryFlagIsSetTo0AfterPreviousCarry() {
	suite.vm.registers[15] = 1

	suite.vm.Load([]byte{
		SET_REGISTER_0, 0x0A, // Set register 0 to ...
		SET_REGISTER_1, 0x0A, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})
	suite.vm.Run()

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVY() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x0A, // Set register 0 to 10
		SET_REGISTER_1, 0x01, // Set register 1 to 1
		0x80, 0x15, // Set VX to 10 - 1
	})

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVYUnderflow() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x0A, // Set register 0 to 10
		SET_REGISTER_1, 0x0B, // Set register 1 to 1
		0x80, 0x15, // Set VX to 10 - 11
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVX() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x01, // Set register 0 to 1
		SET_REGISTER_1, 0x0A, // Set register 1 to 10
		0x80, 0x17, // Set VX to 10 - 1
	})

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVXUnderflow() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x0B, // Set register 0 to 11
		SET_REGISTER_1, 0x0A, // Set register 1 to 10
		0x80, 0x17, // Set VX to 10 - 1
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

//8XY6
func (suite *Chip8TestSuite) TestVXShiftRight() {
	suite.executeInstruction([]byte{
		SET_REGISTER_1, 0b11111110, // Set register 1 to 10
		0x80, 0x16, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b11111110), suite.vm.registers[1])
	suite.Equal(byte(0b01111111), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

//8XY6
func (suite *Chip8TestSuite) TestVXShiftRightWithOverflow() {
	suite.executeInstruction([]byte{
		SET_REGISTER_1, 0b00110001, // Set register 1 to 10
		0x80, 0x16, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b00110001), suite.vm.registers[1])
	suite.Equal(byte(0b00011000), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXShiftLeft() {
	suite.executeInstruction([]byte{
		SET_REGISTER_1, 0b01111110, // Set register 1 to 10
		0x80, 0x1E, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b01111110), suite.vm.registers[1])
	suite.Equal(byte(0b11111100), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXShiftLeftWithOverflow() {
	suite.executeInstruction([]byte{
		SET_REGISTER_1, 0b11111100, // Set register 1 to 10
		0x80, 0x1E, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b11111100), suite.vm.registers[1])
	suite.Equal(byte(0b11111000), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter0() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x00, // Set register 0 to 0
		FONT_REGISTER_0, 0x29, // Set Index to point to character 0
	})

	suite.Equal(uint16(FONT_MEMORY), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter1() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x01, // Set register 0 to 1
		FONT_REGISTER_0, 0x29, // Set Index to point to character 1
	})

	suite.Equal(uint16(FONT_MEMORY+5), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter2() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x02, // Set register 0 to 2
		FONT_REGISTER_0, 0x29, // Set Index to point to character 2
	})

	suite.Equal(uint16(FONT_MEMORY+10), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacterF() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x0F, // Set register 0 to F
		FONT_REGISTER_0, 0x29, // Set Index to point to character F
	})

	suite.Equal(uint16(FONT_MEMORY+(0x0F*5)), suite.vm.indexRegister)
}

type MockRandom struct {
	fakeRandom byte
}

func (mockRandom MockRandom) Generate() byte {
	return mockRandom.fakeRandom
}

func (suite *Chip8TestSuite) TestRandomNumber() {
	suite.verifyRandomIsStoredInRegister(0xC0, 0xFF, 0x55, 0x55, 0)
	suite.verifyRandomIsStoredInRegister(0xC1, 0xFF, 0x55, 0x55, 1)
	suite.verifyRandomIsStoredInRegister(0xC2, 0xEF, 0x55, 0x45, 2)
}

func (suite *Chip8TestSuite) verifyRandomIsStoredInRegister(instruction byte, bitmask byte, fakeRandom byte, expected int, expectedRegister int) {
	m := mockDisplay{false, drawPatternValues{}, true}
	suite.vm = NewChip8VM(&m)

	r := MockRandom{fakeRandom}
	var random Random = r

	suite.vm.SetRandom(random)

	suite.vm.Load([]byte{
		instruction, bitmask, // Random number into register 0, ANDed with 0xFF
	})
	suite.vm.Run()

	suite.Equal(byte(expected), suite.vm.registers[expectedRegister])
}

func (suite *Chip8TestSuite) TestDecimalConversion() {
	suite.executeInstruction([]byte{
		SET_REGISTER_0, 0x7B, // Set register 0 to 123
		0xA4, 0x00, // Set index register to point to address 0x400 (1024)
		0xF0, 0x33, // Convert number in register 0 and store in index register
	})

	suite.Equal(byte(1), suite.vm.Memory[1024])
	suite.Equal(byte(2), suite.vm.Memory[1025])
	suite.Equal(byte(3), suite.vm.Memory[1026])
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
