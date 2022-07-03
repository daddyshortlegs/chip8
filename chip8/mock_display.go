package chip8

type mockDisplay struct {
	screenCleared bool
	values        drawPatternValues
}

type drawPatternValues struct {
	x             byte
	y             byte
	address       uint16
	numberOfBytes byte
}

func (m *mockDisplay) DrawPattern(address uint16, numberOfBytes byte, x byte, y byte) {
	values := drawPatternValues{
		x, y, address, numberOfBytes,
	}
	m.values = values
}

func (m *mockDisplay) ClearScreen() {
	m.screenCleared = true
}
