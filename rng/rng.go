package rng

import (
	"math/rand"
	"time"
)

/*
	Too many s̶o̶r̶c̶e̶r̶e̶r̶s̶ repetitions in code. DRY.
*/

// Intn the recreation of the timestamp on each call assures a different seed.
func Intn(ceil int) int {
	var source = rand.NewSource(time.Now().UnixNano())
	var random = rand.New(source)
	return random.Intn(ceil)
}
