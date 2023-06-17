package chip8

import (
	"errors"
	"fmt"
)

type stack struct {
	address [16]uint16
	index   int
}

func (s *stack) length() int {
	return s.index
}

func (s *stack) Push(value uint16) error {
	if s.index >= len(s.address) {
		return errors.New("stack overflow")
	}
	s.address[s.index] = value
	fmt.Printf("stack push %X\n", value)

	s.index++
	return nil
}

func (s *stack) Pop() (uint16, error) {
	s.index--
	if s.index < 0 {
		return 0, errors.New("stack empty")
	}
	value := s.address[s.index]
	fmt.Printf("stack pop %X\n", value)

	return value, nil
}
