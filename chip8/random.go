package chip8

import "math/rand"

type Random interface {
	Generate() byte
}

type PseudoRandom struct {
	Seed int64
}

func (pseudoRandom PseudoRandom) Generate() byte {
	return byte(rand.Intn(256))
}
