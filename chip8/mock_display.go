package chip8

type mockDisplay struct {
	screenCleared bool
	values        drawPatternValues
	poll          bool
}

type drawPatternValues struct {
	x             byte
	y             byte
	address       uint16
	numberOfBytes byte
}

func (m *mockDisplay) DrawSprite(chip8 *Chip8VM, address uint16, numberOfBytes byte, x byte, y byte) {
	values := drawPatternValues{
		x, y, address, numberOfBytes,
	}
	m.values = values
}

func (m *mockDisplay) ClearScreen() {
	m.screenCleared = true
}

func (m *mockDisplay) PollEvents() bool {
	return m.poll
}

func (m *mockDisplay) setPollEvents(poll bool) {
	m.poll = poll
}
