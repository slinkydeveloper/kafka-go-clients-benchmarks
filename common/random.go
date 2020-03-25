package common

import (
	"math/rand"
	"time"
)

func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenKey(rand *rand.Rand) []byte {
	return gen(rand, KeySize)
}

func GenPayload(rand *rand.Rand) []byte {
	return gen(rand, PayloadSize)
}

func gen(rand *rand.Rand, size int) []byte {
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}