package chip8

import (
	"math/rand"
	"time"
)

type Random interface {
	Generate() byte
}

type PseudoRandom struct {
	Seed int64
}

func NewRandom() *PseudoRandom {
	r := new(PseudoRandom)
	r.Seed = time.Now().UnixNano()
	return r
}

func (pseudoRandom PseudoRandom) Generate() byte {
	return byte(rand.Intn(256))
}
