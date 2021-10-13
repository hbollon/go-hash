package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"math"
)

type HashType int

const (
	MD5 HashType = iota
	SHA1
)

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

func Hash(input string, method HashType) (string, error) {
	var hash []byte
	if method == MD5 {
		tmp := md5.Sum([]byte(input))
		hash = tmp[:]
	} else if method == SHA1 {
		tmp := sha1.Sum([]byte(input))
		hash = tmp[:]
	} else {
		return "", errors.New("Unknown hash method")
	}
	return hex.EncodeToString(hash[:]), nil
}
