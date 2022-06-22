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
	theStack.push(uint16(0x00cc))
	suite.Equal(1, theStack.length())
}

func (suite *StackTestSuite) TestAnotherPush() {
	theStack := stack{}
	theStack.push(uint16(0x00cc))
	theStack.push(uint16(0x00cc))
	suite.Equal(2, theStack.length())
}

func (suite *StackTestSuite) TestPop() {
	theStack := stack{}
	theStack.push(uint16(0x1111))
	theStack.push(uint16(0x2222))
	result, _ := theStack.pop()
	suite.Equal(uint16(0x2222), result)
}

func (suite *StackTestSuite) TestPopTwice() {
	theStack := stack{}
	theStack.push(uint16(0x1111))
	theStack.push(uint16(0x2222))
	theStack.pop()
	result, _ := theStack.pop()
	suite.Equal(uint16(0x1111), result)
}

func (suite *StackTestSuite) TestBlowStack() {
	theStack := stack{}
	for i := 0; i < 16; i++ {
		theStack.push(uint16(0x1111))
	}

	err := theStack.push(uint16(0x1111))
	suite.Equal(errors.New("stack overflow"), err)
}

func (suite *StackTestSuite) TestPopEmptyStack() {
	theStack := stack{}
	_, err := theStack.pop()
	suite.Equal(errors.New("stack empty"), err)
}

func TestStackTestSuite(t *testing.T) {
	suite.Run(t, new(StackTestSuite))
}
