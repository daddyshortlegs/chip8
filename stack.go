package chip8

type stack struct {
	address [16]uint
	index   int
}

func (s stack) length() int {
	return s.index
}

func (s *stack) push(value uint) {
	s.address[s.index] = value
	s.index++
}
