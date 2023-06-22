package chip8

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StackTestSuite struct {
	suite.Suite
}

func (suite *StackTestSuite) TestEmptyStack() {
	theStack := stack{}
	suite.Equal(0, theStack.length())
}

func (suite *StackTestSuite) TestPush() {
	theStack := stack{}
	theStack.Push(uint16(0x00cc))
	suite.Equal(1, theStack.length())
}

func (suite *StackTestSuite) TestAnotherPush() {
	theStack := stack{}
	theStack.Push(uint16(0x00cc))
	theStack.Push(uint16(0x00cc))
	suite.Equal(2, theStack.length())
}

func (suite *StackTestSuite) TestPop() {
	theStack := stack{}
	theStack.Push(uint16(0x1111))
	theStack.Push(uint16(0x2222))
	result, _ := theStack.Pop()
	suite.Equal(uint16(0x2222), result)
}

func (suite *StackTestSuite) TestPopTwice() {
	theStack := stack{}
	theStack.Push(uint16(0x1111))
	theStack.Push(uint16(0x2222))
	result, _ := theStack.Pop()
	suite.Equal(uint16(0x2222), result)
	result, _ = theStack.Pop()
	suite.Equal(uint16(0x1111), result)
}

//TODO: We should not be able to pop 3 times
//func (suite *StackTestSuite) TestPopThrice() {
//	theStack := stack{}
//	theStack.Push(uint16(0x1111))
//	theStack.Push(uint16(0x2222))
//	result, _ := theStack.Pop()
//	suite.Equal(uint16(0x2222), result)
//	result, _ = theStack.Pop()
//	suite.Equal(uint16(0x1111), result)
//	result, _ = theStack.Pop()
//	suite.Equal(uint16(0x1111), result)
//}

func (suite *StackTestSuite) TestBlowStack() {
	theStack := stack{}
	for i := 0; i < 16; i++ {
		theStack.Push(uint16(0x1111))
	}

	err := theStack.Push(uint16(0x1111))
	suite.Equal(errors.New("stack overflow"), err)
}

func (suite *StackTestSuite) TestPopEmptyStack() {
	theStack := stack{}
	_, err := theStack.Pop()
	suite.Equal(errors.New("stack empty"), err)
}

func TestStackTestSuite(t *testing.T) {
	suite.Run(t, new(StackTestSuite))
}
