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

func (a *Alphabet) I2c(input uint64) string {
	var size, index int
	for size = a.min; size <= a.max; size++ {
		if a.levelledPossibilities[index] >= input {
			break
		}
		index++
	}

	return a.i2cSameSize(input, size)
}

func (a *Alphabet) i2cSameSize(input uint64, size int) string {
	coeff := uint64(a.length)
	for input >= coeff {
		input -= coeff
		coeff *= uint64(a.length)
	}

	var strBuilder string
	for i := 0; i < size; i++ {
		letter := input % uint64(a.length)
		strBuilder = string(a.alphabet[letter]) + strBuilder
		input = uint64(math.RoundToEven(float64(input) / float64(a.length)))
	}

	return strBuilder
}
