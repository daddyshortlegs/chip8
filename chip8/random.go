package chip8

type Random interface {
	Generate() byte
}
