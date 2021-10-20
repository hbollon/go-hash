package main

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Alphabet struct {
	alphabet              string
	length                int
	min                   int
	max                   int
	possibilities         uint64
	levelledPossibilities []uint64
}

func GenerateAlphabet(alphabet string, min, max int) Alphabet {
	var result uint64
	levelledResult := make([]uint64, max)
	for i := 1; i <= max; i++ {
		if i >= min {
			result += uint64(math.Pow(float64(len(alphabet)), float64(i)))
		}
		levelledResult[i-1] = uint64(math.Pow(float64(len(alphabet)), float64(i)))
	}
	return Alphabet{
		alphabet:              alphabet,
		length:                len(alphabet),
		min:                   min,
		max:                   max,
		possibilities:         uint64(result),
		levelledPossibilities: levelledResult,
	}
}

func (a *Alphabet) RandomIndex() uint64 {
	return rand.Uint64() % a.possibilities
}
