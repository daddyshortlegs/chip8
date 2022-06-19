package chip8

import (
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
	theStack.push(uint(0x00cc))
	suite.Equal(1, theStack.length())
}

func TestStackTestSuite(t *testing.T) {
	suite.Run(t, new(StackTestSuite))
}
