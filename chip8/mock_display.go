package chip8

type mockDisplay struct {
	screenCleared bool
	values        drawPatternValues
	eventType     EventType
	fakeKey       int
}

func (k *mockDisplay) GetKey() int {
	return k.fakeKey
}

func (k *mockDisplay) SetKey(key int) {
	k.fakeKey = key
}

type drawPatternValues struct {
	x             byte
	y             byte
	address       uint16
	numberOfBytes byte
}

func (m *mockDisplay) DrawSprite(chip8 *VM, address uint16, numberOfBytes byte, x byte, y byte) {
	values := drawPatternValues{
		x, y, address, numberOfBytes,
	}
	m.values = values
}

func (m *mockDisplay) ClearScreen() {
	m.screenCleared = true
}

func (m *mockDisplay) PollEvents() EventType {
	return m.eventType
}

func (m *mockDisplay) setPollEvents(eventType EventType) {
	m.eventType = eventType
}
