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
	asm         *Assembler
}

const FontMemory = 0x50
const programStart = 0x200

func (suite *Chip8TestSuite) SetupTest() {
	suite.mockDisplay = mockDisplay{false, drawPatternValues{}, KeyboardEvent, 4}
	suite.mockRandom = MockRandom{55}
	suite.vm = NewVM(&suite.mockDisplay, suite.mockRandom)
	suite.asm = NewAssembler()
}

func (suite *Chip8TestSuite) executeInstructions() {
	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()
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

func (suite *Chip8TestSuite) TestFetchAndSetRegisters() {
	suite.asm.SetRegister(0, 0x011)
	suite.asm.SetRegister(1, 0x012)
	suite.asm.SetRegister(5, 0x0CC)

	suite.executeInstructions()

	suite.Equal(byte(0x11), suite.vm.registers[0])
	suite.Equal(byte(0x12), suite.vm.registers[1])
	suite.Equal(byte(0xCC), suite.vm.registers[5])
}

func (suite *Chip8TestSuite) TestAddToRegister() {
	suite.asm.AddToRegister(0, 0x0A)

	suite.executeInstructions()
	suite.Equal(byte(0x0A), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestAddToRegisterMultipleTimes() {
	suite.asm.AddToRegister(0, 0x01)
	suite.asm.AddToRegister(0, 0x01)
	suite.asm.AddToRegister(0, 0x01)

	suite.executeInstructions()
	suite.Equal(byte(0x03), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetAndAddToRegister() {
	suite.asm.SetRegister(0, 0x01)
	suite.asm.AddToRegister(0, 0x0A)

	suite.executeInstructions()
	suite.Equal(byte(0x0B), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestSetIndexRegister() {
	suite.asm.SetIndexRegister(0x0A)

	suite.executeInstructions()
	suite.Equal(uint16(0x0A), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestSetIndexRegisterWith12BitValue() {
	suite.asm.SetIndexRegister(0xFFF)

	suite.executeInstructions()
	suite.Equal(uint16(0xFFF), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestSetJumpToAddress() {
	suite.asm.Jump(0x300)

	suite.executeInstructions()
	suite.Equal(uint16(0x302), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestClearScreen() {
	suite.asm.ClearScreen()

	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()

	suite.True(suite.mockDisplay.screenCleared)
}

func (suite *Chip8TestSuite) TestGetCoordinatesFromRegisters_whenDraw() {
	suite.vm.registers[5] = 20
	suite.vm.registers[10] = 30

	suite.asm.Display(5, 0xA, 0)
	suite.vm.Load(suite.asm.Assemble())

	suite.vm.Run()

	suite.Equal(byte(20), suite.vm.getXCoordinate())
	suite.Equal(byte(30), suite.vm.getYCoordinate())
}

func (suite *Chip8TestSuite) TestCoordinatesShouldWrap() {
	suite.vm.registers[5] = 64
	suite.vm.registers[10] = 32
	suite.vm.indexRegister = 0x200

	suite.asm.Display(5, 0xA, 5)

	suite.vm.Load(suite.asm.Assemble())

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
	suite.asm.SetRegister(0x5, 0x14)
	suite.asm.SetRegister(0xA, 0x1E)
	suite.asm.SetIndexRegister(0x50)
	suite.asm.Display(0x5, 0xA, 5)

	suite.vm.Load(suite.asm.Assemble())

	suite.vm.Run()

	suite.Equal(byte(20), suite.mockDisplay.values.x)
	suite.Equal(byte(30), suite.mockDisplay.values.y)
	suite.Equal(uint16(0x50), suite.mockDisplay.values.address)
	suite.Equal(byte(5), suite.mockDisplay.values.numberOfBytes)
}

func (suite *Chip8TestSuite) TestVXIsSetToVY() {
	suite.asm.SetRegister(0x5, 0x14)
	suite.asm.Set(0x0, 0x5)

	suite.executeInstructions()

	suite.Equal(byte(0x14), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryORofVXVY() {
	suite.asm.SetRegister(0x0, 0x0F)
	suite.asm.SetRegister(0x1, 0xF0)
	suite.asm.Or(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0xFF), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryANDofVXVY() {
	suite.asm.SetRegister(0x0, 0b00001111)
	suite.asm.SetRegister(0x1, 0b00110011)
	suite.asm.And(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0b00000011), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestVXIsSetToBinaryXORofVXVY() {
	suite.asm.SetRegister(0x0, 0b00001111)
	suite.asm.SetRegister(0x1, 0b00110011)
	suite.asm.Xor(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0b00111100), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestAddWithNoCarry() {
	suite.asm.SetRegister(0x0, 0x0A)
	suite.asm.SetRegister(0x1, 0x0A)
	suite.asm.Add(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestAddWithCarry() {
	suite.asm.SetRegister(0x0, 0xFF)
	suite.asm.SetRegister(0x1, 0x01)
	suite.asm.Add(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0x00), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestCarryFlagIsSetTo0AfterPreviousCarry() {
	suite.vm.registers[15] = 1

	suite.asm.SetRegister(0x0, 0x0A)
	suite.asm.SetRegister(0x1, 0x0A)
	suite.asm.Add(0x0, 0x1)

	suite.executeInstructions()

	suite.vm.Run()

	suite.Equal(byte(0x14), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVY() {
	suite.asm.SetRegister(0x0, 0x0A)
	suite.asm.SetRegister(0x1, 0x01)
	suite.asm.Subtract(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXSubtractVYUnderflow() {
	suite.asm.SetRegister(0x0, 0x0A)
	suite.asm.SetRegister(0x1, 0x0B)
	suite.asm.Subtract(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVX() {
	suite.asm.SetRegister(0x0, 0x01)
	suite.asm.SetRegister(0x1, 0x0A)
	suite.asm.SubtractLast(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0x09), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVYSubtractVXUnderflow() {
	suite.asm.SetRegister(0x0, 0x0B)
	suite.asm.SetRegister(0x1, 0x0A)
	suite.asm.SubtractLast(0x0, 0x1)

	suite.executeInstructions()

	suite.Equal(byte(0xFF), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

// 8XY6
func (suite *Chip8TestSuite) TestVXShiftRight() {
	suite.asm.SetRegister(1, 0b11111110)
	suite.asm.ShiftRight(0, 1)

	suite.executeInstructions()

	suite.Equal(byte(0b11111110), suite.vm.registers[1])
	suite.Equal(byte(0b01111111), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

// 8XY6
func (suite *Chip8TestSuite) TestVXShiftRightWithOverflow() {
	suite.asm.SetRegister(1, 0b00110001)
	suite.asm.ShiftRight(0, 1)

	suite.executeInstructions()

	suite.Equal(byte(0b00110001), suite.vm.registers[1])
	suite.Equal(byte(0b00011000), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXShiftLeft() {
	suite.asm.SetRegister(1, 0b01111110)
	suite.asm.ShiftLeft(0, 1)

	suite.executeInstructions()

	suite.Equal(byte(0b01111110), suite.vm.registers[1])
	suite.Equal(byte(0b11111100), suite.vm.registers[0])
	suite.Equal(byte(0), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestVXShiftLeftWithOverflow() {
	suite.asm.SetRegister(1, 0b11111100)
	suite.asm.ShiftLeft(0, 1)

	suite.executeInstructions()

	suite.Equal(byte(0b11111100), suite.vm.registers[1])
	suite.Equal(byte(0b11111000), suite.vm.registers[0])
	suite.Equal(byte(1), suite.vm.registers[15])
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter0() {
	suite.asm.SetRegister(0, 0x0)
	suite.asm.FontChar(0)

	suite.executeInstructions()

	suite.Equal(uint16(FontMemory), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter1() {
	suite.asm.SetRegister(0, 0x1)
	suite.asm.FontChar(0)

	suite.executeInstructions()

	suite.Equal(uint16(FontMemory+5), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacter2() {
	suite.asm.SetRegister(0, 0x2)
	suite.asm.FontChar(0)

	suite.executeInstructions()

	suite.Equal(uint16(FontMemory+10), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestIndexPointsToCharacterF() {
	suite.asm.SetRegister(0, 0xF)
	suite.asm.FontChar(0)

	suite.executeInstructions()

	suite.Equal(uint16(FontMemory+(0x0F*5)), suite.vm.indexRegister)
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
	suite.asm.SetRegister(0, 0x7B)
	suite.asm.SetIndexRegister(0x400)
	suite.asm.BCD(0)

	suite.executeInstructions()

	suite.Equal(byte(1), suite.vm.Memory[1024])
	suite.Equal(byte(2), suite.vm.Memory[1025])
	suite.Equal(byte(3), suite.vm.Memory[1026])
}

func (suite *Chip8TestSuite) TestDoesSkipIfEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SkipIfEqual(0, 0x00)
	suite.executeInstructions()
	suite.Equal(uint16(0x206), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesSkipIfEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x11)
	suite.asm.SkipIfEqual(0, 0x11)
	suite.executeInstructions()

	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesNotSkipIfNotEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x11)
	suite.asm.SkipIfEqual(0, 0x22)
	suite.executeInstructions()

	suite.Equal(uint16(0x206), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesSkipIfNotEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x11)
	suite.asm.SkipIfNotEqual(0, 0x22)
	suite.executeInstructions()

	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesNotSkipIfEqualToRegister() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x33)
	suite.asm.SkipIfNotEqual(0, 0x33)
	suite.executeInstructions()

	suite.Equal(uint16(0x206), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestSkipsWhenVxAndVyAreEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x15)
	suite.asm.SetRegister(1, 0x15)
	suite.asm.SkipIfRegistersEqual(0, 1)
	suite.executeInstructions()

	suite.Equal(uint16(0x20A), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestProgramCounterIncrementsAfterJump() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x15)
	// 0x202
	suite.asm.SetRegister(1, 0x15)
	// 0x204
	suite.asm.SkipIfRegistersEqual(0, 1)
	// 0x208
	suite.asm.SetRegister(2, 0x69) // This operation gets skipped so register won't be set
	suite.asm.SetRegister(3, 0x77) // but skip moves PC to this instruction
	suite.executeInstructions()

	suite.Equal(uint16(0x20C), suite.vm.pc)
	suite.Equal(byte(0x15), suite.vm.registers[0])
	suite.Equal(byte(0x15), suite.vm.registers[1])
	suite.Equal(byte(0x00), suite.vm.registers[2]) // this is zero as set instruction skipped
	suite.Equal(byte(0x77), suite.vm.registers[3]) // this is zero as set instruction skipped
}

func (suite *Chip8TestSuite) TestDoesNotSkipWhenVxAndVyAreNotEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x15)
	suite.asm.SetRegister(1, 0x20)
	suite.asm.SkipIfRegistersEqual(0, 1)
	suite.executeInstructions()

	suite.Equal(uint16(0x208), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestSkipsWhenVxAndVyAreNotEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x15)
	suite.asm.SetRegister(1, 0x20)
	suite.asm.SkipIfRegistersNotEqual(0, 1)
	suite.executeInstructions()

	suite.Equal(uint16(0x20A), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestDoesNotSkipWhenVxAndVyAreEqual() {
	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x25)
	suite.asm.SetRegister(1, 0x25)
	suite.asm.SkipIfRegistersNotEqual(0, 1)
	suite.executeInstructions()

	suite.Equal(uint16(0x208), suite.vm.pc)
}

//func (suite *Chip8TestSuite) TestJumpToSubroutineUpdatesProgramCounterAndPushesToStack() {
//	suite.asm.Sub(0x345)
//
//	suite.executeInstructions()
//	suite.Equal(uint16(0x345), suite.vm.pc)
//	value, _ := suite.vm.theStack.Pop()
//	suite.Equal(uint16(programStart+2), value)
//}

func (suite *Chip8TestSuite) TestReturnFromSubroutine() {
	suite.vm.theStack.Push(0xA12)

	suite.asm.Return()

	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()
	suite.Equal(uint16(0xA14), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestJumpAndReturnFromSubroutine() {
	const programStart = 0x200
	suite.asm.Sub(0x20C)
	suite.asm.SetRegister(0, 0x11)
	suite.asm.SetRegister(1, 0x22)
	suite.asm.SetRegister(2, 0x33)
	suite.asm.SetRegister(3, 0x44)
	suite.asm.Data([]byte{0x00, 0x00})
	suite.asm.SetRegister(5, 0xAA) // subroutine
	suite.asm.Return()

	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()
	suite.Equal(byte(0x11), suite.vm.registers[0])
	suite.Equal(byte(0x22), suite.vm.registers[1])
	suite.Equal(byte(0x33), suite.vm.registers[2])
	suite.Equal(byte(0x44), suite.vm.registers[3])
	suite.Equal(byte(0xAA), suite.vm.registers[5])
}

func (suite *Chip8TestSuite) TestGetKey() {
	suite.mockDisplay = mockDisplay{false, drawPatternValues{}, KeyboardEvent, 55}
	suite.mockRandom = MockRandom{55}
	suite.mockDisplay.SetKey(55)
	suite.vm = NewVM(&suite.mockDisplay, suite.mockRandom)

	suite.asm.GetKey(3)
	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()
	suite.Equal(byte(55), suite.vm.registers[3])
}

func (suite *Chip8TestSuite) TestRegister0AddToIndex() {
	suite.asm.SetRegister(0, 0x25)
	suite.asm.AddToIndex(0)

	suite.executeInstructions()

	suite.Equal(uint16(0x25), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestRegister1AddToIndex() {
	suite.asm.SetRegister(1, 0x55)
	suite.asm.AddToIndex(1)

	suite.executeInstructions()

	suite.Equal(uint16(0x55), suite.vm.indexRegister)
}

func (suite *Chip8TestSuite) TestStoreSingleRegisterToMemory() {
	suite.asm.SetIndexRegister(0x200)
	suite.asm.SetRegister(0, 0x69)
	suite.asm.Store(0)
	suite.executeInstructions()
	suite.Equal(uint8(0x69), suite.vm.Memory[0x200])
	suite.Equal(uint8(0x00), suite.vm.Memory[0x201])
}

func (suite *Chip8TestSuite) TestStoreTwoRegistersToMemory() {
	suite.asm.SetIndexRegister(0x200)
	suite.asm.SetRegister(0, 0x69)
	suite.asm.SetRegister(1, 0x58)
	suite.asm.SetRegister(2, 0x47)
	suite.asm.SetRegister(3, 0x47)
	suite.asm.SetRegister(4, 0x36)
	suite.asm.SetRegister(5, 0x44)
	suite.asm.SetRegister(6, 0x33)
	suite.asm.SetRegister(7, 0x22)
	suite.asm.SetRegister(8, 0x11)
	suite.asm.SetRegister(9, 0x00)
	suite.asm.SetRegister(10, 0xAA)
	suite.asm.SetRegister(11, 0xBB)
	suite.asm.SetRegister(12, 0xCC)
	suite.asm.SetRegister(13, 0xDD)
	suite.asm.SetRegister(14, 0xEE)
	suite.asm.SetRegister(15, 0xFF)
	suite.asm.Store(0xF)

	suite.executeInstructions()

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
	suite.asm.SetIndexRegister(0x300)
	suite.asm.Load(0)

	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Memory[0x300] = 0x88
	suite.vm.Run()
	suite.Equal(uint8(0x88), suite.vm.registers[0])
}

func (suite *Chip8TestSuite) TestLoadMultipleRegistersFromMemory() {
	suite.asm.SetIndexRegister(0x300)
	suite.asm.Load(0xF)

	suite.vm.Load(suite.asm.Assemble())

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

func (suite *Chip8TestSuite) TestJumpWithoutOffset() {
	suite.asm.Jump(0x345)
	suite.executeInstructions()
	suite.Equal(uint16(0x347), suite.vm.pc)
}

func (suite *Chip8TestSuite) TestJumpWithOffset() {
	data := []byte{0xB3, 0x45}
	suite.vm.Load(data)
	suite.vm.registers[0] = 0x10
	suite.vm.Run()
	suite.Equal(uint16(0x355), suite.vm.pc)
}

/*
TODO:


// FX18 sets sound timer to value in VX

EX9E and EXA1: Skip if key
FX18: Timers

*/

func (suite *Chip8TestSuite) TestSkipIfKeyPressed() {
	suite.mockDisplay = mockDisplay{false, drawPatternValues{}, KeyboardEvent, 55}
	suite.mockRandom = MockRandom{55}
	suite.mockDisplay.SetKey(0xC)
	suite.vm = NewVM(&suite.mockDisplay, suite.mockRandom)

	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0xC)
	suite.asm.SkipIfKeyPressed(0)
	suite.asm.SetRegister(0xA, 0x22) // instead this registr should NOT be set
	suite.asm.SetRegister(0xB, 0x69) // instead this registr should be set
	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()

	suite.executeInstructions()

	suite.NotEqual(byte(0x22), suite.vm.registers[0xA])
	suite.Equal(byte(0x69), suite.vm.registers[0xB])
}

func (suite *Chip8TestSuite) TestSkipIfKeyNotPressed() {
	suite.mockDisplay = mockDisplay{false, drawPatternValues{}, KeyboardEvent, 55}
	suite.mockRandom = MockRandom{55}
	suite.mockDisplay.SetKey(0xC)
	suite.vm = NewVM(&suite.mockDisplay, suite.mockRandom)

	suite.Equal(uint16(0x200), suite.vm.pc)

	suite.asm.SetRegister(0, 0x9)
	suite.asm.SkipIfKeyNotPressed(0)
	suite.asm.SetRegister(0xA, 0x22) // This register should NOT be set as it is skipped
	suite.asm.SetRegister(0xB, 0x69) // instead this registr should be set
	suite.vm.Load(suite.asm.Assemble())
	suite.vm.Run()

	suite.executeInstructions()

	suite.NotEqual(byte(0x22), suite.vm.registers[0xA])
	suite.Equal(byte(0x69), suite.vm.registers[0xB])
}

func TestChip8TestSuite(t *testing.T) {
	suite.Run(t, new(Chip8TestSuite))
}
