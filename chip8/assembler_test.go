package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssemblerTestSuite struct {
	suite.Suite
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

func (suite *AssemblerTestSuite) TestSubroutine() {
	theAssembler := NewAssembler()
	theAssembler.Sub(0x643)
	suite.Equal([]byte{0x26, 0x43}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestReturn() {
	theAssembler := NewAssembler()
	theAssembler.Return()
	suite.Equal([]byte{0x00, 0xEE}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSkipIfEqual() {
	theAssembler := NewAssembler()
	theAssembler.SkipIfEqual(1, 0xDD)
	suite.Equal([]byte{0x31, 0xDD}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSkipIfNotEqual() {
	theAssembler := NewAssembler()
	theAssembler.SkipIfNotEqual(2, 0xCD)
	suite.Equal([]byte{0x42, 0xCD}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSkipIfRegistersEqual() {
	theAssembler := NewAssembler()
	theAssembler.SkipIfRegistersEqual(3, 4)
	suite.Equal([]byte{0x53, 0x40}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSkipIfRegistersNotEqual() {
	theAssembler := NewAssembler()
	theAssembler.SkipIfRegistersNotEqual(2, 7)
	suite.Equal([]byte{0x92, 0x70}, theAssembler.Assemble())
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

func (suite *AssemblerTestSuite) TestAddToRegister() {
	theAssembler := NewAssembler()
	theAssembler.AddToRegister(0, 0x33)
	suite.Equal([]byte{0x70, 0x33}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSet() {
	theAssembler := NewAssembler()
	theAssembler.Set(6, 1)
	suite.Equal([]byte{0x86, 0x10}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestOr() {
	theAssembler := NewAssembler()
	theAssembler.Or(2, 3)
	suite.Equal([]byte{0x82, 0x31}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestAnd() {
	theAssembler := NewAssembler()
	theAssembler.And(4, 5)
	suite.Equal([]byte{0x84, 0x52}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestXor() {
	theAssembler := NewAssembler()
	theAssembler.Xor(2, 6)
	suite.Equal([]byte{0x82, 0x63}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestAdd() {
	theAssembler := NewAssembler()
	theAssembler.Add(1, 5)
	suite.Equal([]byte{0x81, 0x54}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSubtract() {
	theAssembler := NewAssembler()
	theAssembler.Subtract(6, 5)
	suite.Equal([]byte{0x86, 0x55}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSubtractLast() {
	theAssembler := NewAssembler()
	theAssembler.SubtractLast(1, 2)
	suite.Equal([]byte{0x81, 0x27}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestShiftRight() {
	theAssembler := NewAssembler()
	theAssembler.ShiftRight(0, 2)
	suite.Equal([]byte{0x80, 0x26}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestShiftLeft() {
	theAssembler := NewAssembler()
	theAssembler.ShiftLeft(1, 2)
	suite.Equal([]byte{0x81, 0x2E}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestIndexRegister() {
	theAssembler := NewAssembler()
	theAssembler.SetIndexRegister(0x123)
	suite.Equal([]byte{0xA1, 0x23}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestJumpWithOffset() {
	theAssembler := NewAssembler()
	theAssembler.SetJumpWithOffset(0x39A)
	suite.Equal([]byte{0xB3, 0x9A}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestRandom() {
	theAssembler := NewAssembler()
	theAssembler.Random(2, 0x66)
	suite.Equal([]byte{0xC2, 0x66}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestDisplay() {
	theAssembler := NewAssembler()
	theAssembler.Display(6, 4, 5)
	suite.Equal([]byte{0xD6, 0x45}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSkipIfKeyPressed() {
	theAssembler := NewAssembler()
	theAssembler.SkipIfKeyPressed(2)
	suite.Equal([]byte{0xE2, 0x9E}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSkipIfKeyNotPressed() {
	theAssembler := NewAssembler()
	theAssembler.SkipIfKeyNotPressed(4)
	suite.Equal([]byte{0xE4, 0xA1}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestGetDelayTimer() {
	theAssembler := NewAssembler()
	theAssembler.GetDelayTimer(4)
	suite.Equal([]byte{0xF4, 0x07}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSetDelayTimer() {
	theAssembler := NewAssembler()
	theAssembler.SetDelayTimer(1)
	suite.Equal([]byte{0xF1, 0x15}, theAssembler.Assemble())
}

func (suite *AssemblerTestSuite) TestSetSoundTimer() {
	theAssembler := NewAssembler()
	theAssembler.SetSoundTimer(3)
	suite.Equal([]byte{0xF3, 0x18}, theAssembler.Assemble())
}

// TODO: Add to index

func (suite *AssemblerTestSuite) TestGetKey() {
	theAssembler := NewAssembler()
	theAssembler.GetKey(5)
	suite.Equal([]byte{0xF5, 0x0A}, theAssembler.Assemble())
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
