package main

import "math"

type Alphabet struct {
	alphabet      string
	min           int
	max           int
	possibilities int
}

func GenerateAlphabet(alphabet string, min, max int) Alphabet {
	var result float64
	for i := min; i <= max; i++ {
		result += math.Pow(float64(len(alphabet)), float64(i))
	}
	return Alphabet{
		alphabet:      alphabet,
		min:           min,
		max:           max,
		possibilities: int(result),
	}
}
