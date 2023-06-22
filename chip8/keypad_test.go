package chip8

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeyPadTestSuite struct {
	suite.Suite
}

func (suite *KeyPadTestSuite) TestKeyCodesAreMappedToCorrectValues() {
	suite.Equal(byte(0x1), KeyCodeToValue(49))
	suite.Equal(byte(0x2), KeyCodeToValue(50))
	suite.Equal(byte(0x3), KeyCodeToValue(51))
}

func TestKeyPadTestSuite(t *testing.T) {
	suite.Run(t, new(KeyPadTestSuite))
}
