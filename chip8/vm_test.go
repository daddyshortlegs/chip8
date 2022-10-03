package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Chip8TestSuite struct {
	suite.Suite
	vm          *VM
	mockDisplay mockDisplay
	mockRandom  MockRandom
}

const FontMemory = 0x50
const SetRegister0 = 0x60
const SetRegister2 = 0x62
const SetRegister3 = 0x63
const SetRegister4 = 0x64
const AddRegister0 = 0x70
const FontRegister0 = 0xF0

func (suite *Chip8TestSuite) SetupTest() {
	suite.mockDisplay = mockDisplay{false, drawPatternValues{}, KeyboardEvent, 4}
	suite.mockRandom = MockRandom{55}
	suite.vm = NewVM(&suite.mockDisplay, suite.mockRandom)
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
	verifyRegisterSet(suite, []byte{SetRegister0, 0xFF}, 0, 0xFF)
	verifyRegisterSet(suite, []byte{0x61, 0xEE}, 1, 0xEE)
	verifyRegisterSet(suite, []byte{SetRegister2, 0xDD}, 2, 0xDD)
	verifyRegisterSet(suite, []byte{SetRegister3, 0xCC}, 3, 0xCC)
	verifyRegisterSet(suite, []byte{SetRegister4, 0xBB}, 4, 0xBB)
}

func verifyRegisterSet(suite *Chip8TestSuite, instruction []byte, register int, result int) {
	suite.executeInstruction2(instruction)
	suite.Equal(byte(result), suite.vm.registers[register])
}

func (suite *Chip8TestSuite) TestFetchAndSetAllRegisters() {
	suite.vm.Load([]byte{SetRegister0, 0x11, 0x61, 0x12, 0x65, 0xCC})
	suite.vm.Run()

	suite.Equal(byte(0x11), suite.vm.registers[0])
	suite.Equal(byte(0x12), suite.vm.registers[1])
	suite.Equal(byte(0xCC), suite.vm.registers[5])
}

func (suite *Chip8TestSuite) executeInstruction2(data []byte) {
	m := mockDisplay{false, drawPatternValues{}, QuitEvent, 0}
	r := MockRandom{55}
	suite.vm = NewVM(&m, r)

	suite.executeInstruction(data)
}

func (suite *Chip8TestSuite) executeInstruction(data []byte) {
	suite.vm.Load(data)
	suite.vm.Run()
}

func (suite *Chip8TestSuite) TestAddToRegister() {
	suite.executeInstruction([]byte{AddRegister0, 0x0A})
	suite.Equal(byte(0x0A), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetAndAddToRegister() {
	suite.executeInstruction([]byte{SetRegister0, 0x01, AddRegister0, 0x0A})
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

	suite.True(suite.mockDisplay.screenCleared)
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

	suite.Equal(byte(20), suite.mockDisplay.values.x)
	suite.Equal(byte(30), suite.mockDisplay.values.y)
	suite.Equal(uint16(0x50), suite.mockDisplay.values.address)
	suite.Equal(byte(5), suite.mockDisplay.values.numberOfBytes)
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
		SetRegister0, 0x0F, // Set register 0 to 0x0F
		0x61, 0xF0, // Set register 1 to 0xF0
		0x80, 0x11, // Set register 0 to what's in register 0 & 1 ORd together
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryANDofVXVY() {
	suite.executeInstruction([]byte{
		SetRegister0, 0b00001111, // Set register 0 to ...
		0x61, 0b00110011, // Set register 1 to ...
		0x80, 0x12, // Set register 0 to what's in register 0 & 1 ANDd together
	})

	suite.Equal(byte(0b00000011), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryXORofVXVY() {
	suite.executeInstruction([]byte{
		SetRegister0, 0b00001111, // Set register 0 to ...
		0x61, 0b00110011, // Set register 1 to ...
		0x80, 0x13, // Set register 0 to what's in register 0 & 1 ANDd together
	})

	suite.Equal(byte(0b00111100), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestAddWithNoCarry() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x0A, // Set register 0 to ...
		0x61, 0x0A, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestAddWithCarry() {
	suite.executeInstruction([]byte{
		SetRegister0, 0xFF, // Set register 0 to ...
		0x61, 0x01, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})

	suite.Equal(byte(0x00), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestCarryFlagIsSetTo0AfterPreviousCarry() {
	suite.vm.registers[15] = 1

	suite.vm.Load([]byte{
		SetRegister0, 0x0A, // Set register 0 to ...
		0x61, 0x0A, // Set register 1 to ...
		0x80, 0x14, // Set register 0 to what's in register 0 + 1
	})
	suite.vm.Run()

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVY() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x0A, // Set register 0 to 10
		0x61, 0x01, // Set register 1 to 1
		0x80, 0x15, // Set VX to 10 - 1
	})

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVYUnderflow() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x0A, // Set register 0 to 10
		0x61, 0x0B, // Set register 1 to 1
		0x80, 0x15, // Set VX to 10 - 11
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVX() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x01, // Set register 0 to 1
		0x61, 0x0A, // Set register 1 to 10
		0x80, 0x17, // Set VX to 10 - 1
	})

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVXUnderflow() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x0B, // Set register 0 to 11
		0x61, 0x0A, // Set register 1 to 10
		0x80, 0x17, // Set VX to 10 - 1
	})

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

// 8XY6
func (suite *Chip8TestSuite) TestVXShiftRight() {
	suite.executeInstruction([]byte{
		0x61, 0b11111110, // Set register 1 to 10
		0x80, 0x16, // Set VX to VY and shift right
	})

	suite.Equal(byte(0b11111110), suite.vm.registers[1])
	suite.Equal(byte(0b01111111), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

// 8XY6
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

func (suite *Chip8TestSuite) TestIndexPointsToCharacter0() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x00, // Set register 0 to 0
		FontRegister0, 0x29, // Set Index to point to character 0
	})

	suite.Equal(uint16(FontMemory), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter1() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x01, // Set register 0 to 1
		FontRegister0, 0x29, // Set Index to point to character 1
	})

	suite.Equal(uint16(FontMemory+5), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter2() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x02, // Set register 0 to 2
		FontRegister0, 0x29, // Set Index to point to character 2
	})

	suite.Equal(uint16(FontMemory+10), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacterF() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x0F, // Set register 0 to F
		FontRegister0, 0x29, // Set Index to point to character F
	})

	suite.Equal(uint16(FontMemory+(0x0F*5)), suite.vm.indexRegister)
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
	m := mockDisplay{false, drawPatternValues{}, QuitEvent, 0}
	r := MockRandom{fakeRandom}

	suite.vm = NewVM(&m, r)

	suite.vm.Load([]byte{
		instruction, bitmask, // Random number into register 0, ANDed with 0xFF
	})
	suite.vm.Run()

	suite.Equal(byte(expected), suite.vm.registers[expectedRegister])
}

func (suite *Chip8TestSuite) TestDecimalConversion() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x7B, // Set register 0 to 123
		0xA4, 0x00, // Set index register to point to address 0x400 (1024)
		0xF0, 0x33, // Convert number in register 0 and store in index register
	})

	suite.Equal(byte(1), suite.vm.Memory[1024])
	suite.Equal(byte(2), suite.vm.Memory[1025])
	suite.Equal(byte(3), suite.vm.Memory[1026])
}

func (suite *Chip8TestSuite) TestDoesSkipIfEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		0x30, 0x00})
	suite.Equal(uint16(0x206), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesSkipIfEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x11,
		0x30, 0x11})
	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesNotSkipIfNotEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x11,
		0x30, 0x22})
	suite.Equal(uint16(0x206), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesSkipIfNotEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x11,
		0x40, 0x22})
	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesNotSkipIfEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x33,
		0x40, 0x33})
	suite.Equal(uint16(0x206), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestSkipsWhenVxAndVyAreEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x15,
		0x61, 0x15,
		0x50, 0x10})
	suite.Equal(uint16(0x20A), suite.vm.pc)
}
func (suite *Chip8TestSuite) TestDoesNotSkipWhenVxAndVyAreNotEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x15,
		0x61, 0x20,
		0x50, 0x10})
	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestSkipsWhenVxAndVyAreNotEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x15,
		0x61, 0x25,
		0x90, 0x10})
	suite.Equal(uint16(0x20A), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesNotSkipWhenVxAndVyAreEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)
	suite.executeInstruction([]byte{
		SetRegister0, 0x25,
		0x61, 0x25,
		0x90, 0x10})
	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestJumpToSubroutineUpdatesProgramCounterAndPushesToStack() {
	suite.executeInstruction([]byte{0x23, 0x45})
	suite.Equal(uint16(0x345), suite.vm.pc)
	value, _ := suite.vm.theStack.Pop()
	suite.Equal(uint16(0x345), value)
}

func (suite *Chip8TestSuite) TestReturnFromSubroutine() {
	suite.vm.theStack.Push(0xA12)
	suite.vm.Load([]byte{0x00, 0xEE})
	suite.vm.Run()
	suite.Equal(uint16(0xA12), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestGetKey() {
	suite.mockDisplay = mockDisplay{false, drawPatternValues{}, KeyboardEvent, 55}
	suite.mockRandom = MockRandom{55}
	suite.mockDisplay.SetKey(55)
	suite.vm = NewVM(&suite.mockDisplay, suite.mockRandom)

	suite.vm.Load([]byte{0xF3, 0x0A})
	suite.vm.Run()
	suite.Equal(byte(55), suite.vm.registers[3])
}

func (suite *Chip8TestSuite) TestRegister0AddToIndex() {
	suite.executeInstruction([]byte{
		SetRegister0, 0x25,
		0xF0, 0x1E, // Add value in VX to index register
	})

	suite.Equal(uint16(0x25), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestRegister1AddToIndex() {
	suite.executeInstruction([]byte{
		0x61, 0x55,
		0xF1, 0x1E, // Add value in VX to index register
	})

	suite.Equal(uint16(0x55), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestStoreSingleRegisterToMemory() {
	suite.executeInstruction([]byte{
		0xA2, 0x00, // Set Index register to 200
		SetRegister0, 0x69,
		0xF0, 0x55,
	})
	suite.Equal(uint8(0x69), suite.vm.Memory[0x200])
	suite.Equal(uint8(0x00), suite.vm.Memory[0x201])
}

func (suite *Chip8TestSuite) TestStoreTwoRegistersToMemory() {
	suite.executeInstruction([]byte{
		0xA2, 0x00, // Set Index register to 200
		SetRegister0, 0x69,
		0x61, 0x58,
		0x62, 0x47,
		0x63, 0x47,
		0x64, 0x36,
		0x65, 0x44,
		0x66, 0x33,
		0x67, 0x22,
		0x68, 0x11,
		0x69, 0x00,
		0x6A, 0xAA,
		0x6B, 0xBB,
		0x6C, 0xCC,
		0x6D, 0xDD,
		0x6E, 0xEE,
		0x6F, 0xFF,
		0xFF, 0x55, // Store registers up to F
	})
	suite.Equal(uint8(0x69), suite.vm.Memory[0x200])
	suite.Equal(uint8(0x58), suite.vm.Memory[0x201])
	suite.Equal(uint8(0x47), suite.vm.Memory[0x202])
	suite.Equal(uint8(0x47), suite.vm.Memory[0x203])
	suite.Equal(uint8(0x36), suite.vm.Memory[0x204])
	suite.Equal(uint8(0x44), suite.vm.Memory[0x205])
	suite.Equal(uint8(0x33), suite.vm.Memory[0x206])
	suite.Equal(uint8(0x22), suite.vm.Memory[0x207])
	suite.Equal(uint8(0x11), suite.vm.Memory[0x208])
	suite.Equal(uint8(0x00), suite.vm.Memory[0x209])
	suite.Equal(uint8(0xAA), suite.vm.Memory[0x20A])
	suite.Equal(uint8(0xBB), suite.vm.Memory[0x20B])
	suite.Equal(uint8(0xCC), suite.vm.Memory[0x20C])
	suite.Equal(uint8(0xDD), suite.vm.Memory[0x20D])
	suite.Equal(uint8(0xEE), suite.vm.Memory[0x20E])
	suite.Equal(uint8(0xFF), suite.vm.Memory[0x20F])
}

func (suite *Chip8TestSuite) TestLoadSingleRegisterFromMemory() {
	data := []byte{
		0xA3, 0x00, // Set Index register to 200
		0xF0, 0x65, // Load from memory
	}

	suite.vm.Load(data)
	suite.vm.Memory[0x300] = 0x88
	suite.vm.Run()
	suite.Equal(uint8(0x88), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestLoadMultipleRegistersFromMemory() {
	data := []byte{
		0xA3, 0x00, // Set Index register to 300
		0xFF, 0x65, // Load from memory
	}

	suite.vm.Load(data)
	suite.vm.Memory[0x300] = 0x88
	suite.vm.Memory[0x301] = 0x99
	suite.vm.Memory[0x302] = 0x77
	suite.vm.Memory[0x30F] = 0x69
	suite.vm.Run()
	suite.Equal(uint8(0x88), suite.vm.registers[0])
	suite.Equal(uint8(0x99), suite.vm.registers[1])
	suite.Equal(uint8(0x77), suite.vm.registers[2])
	suite.Equal(uint8(0x00), suite.vm.registers[10])
	suite.Equal(uint8(0x69), suite.vm.registers[15])
}

/*
TODO:



BNNN: Jump with offset
EX9E and EXA1: Skip if key
FX07, FX15 and FX18: Timers
FX55 and FX65: Store and load memory

*/

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
