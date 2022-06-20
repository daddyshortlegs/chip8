package chip8

import "errors"

type stack struct {
	address [16]uint16
	index   int
}

func (s stack) length() int {
	return s.index
}

func (s *stack) push(value uint16) error {
	if s.index >= len(s.address) {
		return errors.New("Stack overflow")
	}
	s.address[s.index] = value
	s.index++
	return nil
}

func (s *stack) pop() (uint16, error) {
	s.index--
	if s.index < 0 {
		return 0, errors.New("Stack empty")
	}
	value := s.address[s.index]
	return value, nil
}
