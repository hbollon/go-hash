package main

import (
	"math"
)

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
	var index int
	levelledResult := make([]uint64, (max-min)+1)
	for i := min; i <= max; i++ {
		result += uint64(math.Pow(float64(len(alphabet)), float64(i)))
		levelledResult[index] = uint64(math.Pow(float64(len(alphabet)), float64(i)))
		index++
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
