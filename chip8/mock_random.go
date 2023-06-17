package chip8

type MockRandom struct {
	fakeRandom byte
}

func (mockRandom MockRandom) Generate() byte {
	return mockRandom.fakeRandom
}
