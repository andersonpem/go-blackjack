package rng

import (
	"crypto/rand"
	"math/big"
)

// RandomInteger generates a cryptographically secure random number in the range [0, ceil).
// It uses crypto/rand for better security and panics on error.
func RandomInteger(ceil int) int {
	if ceil <= 0 {
		panic("ceil must be positive")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(ceil)))
	if err != nil {
		panic(err)
	}

	return int(n.Int64())
}
