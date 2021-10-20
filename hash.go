package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

type HashType string

const (
	MD5  HashType = "MD5"
	SHA1 HashType = "SHA1"
)

func Hash(input string, method HashType) ([]byte, error) {
	var hash []byte
	if method == MD5 {
		tmp := md5.Sum([]byte(input))
		hash = tmp[:]
	} else if method == SHA1 {
		tmp := sha1.Sum([]byte(input))
		hash = tmp[:]
	} else {
		return nil, errors.New("Unknown hash method")
	}
	return hash, nil
}

func HashToString(input string, method HashType) (string, error) {
	hash, err := Hash(input, method)
	return hex.EncodeToString(hash[:]), err
}

func (a *Alphabet) I2c(input uint64) string {
	var size int
	index := a.min - 1
	for size = a.min; size <= a.max; size++ {
		if a.levelledPossibilities[index] >= input {
			break
		}
		index++
	}
	if size > a.max {
		size = a.max
	}

	return a.i2cSameSize(input, size)
}

func (a *Alphabet) i2cSameSize(input uint64, size int) string {
	coeff := uint64(a.length)
	if a.min >= 2 {
		for i := 0; i < a.min-1; i++ {
			input += a.levelledPossibilities[i]
		}
	}
	for input >= coeff {
		input -= coeff
		coeff *= uint64(a.length)
	}

	var strBuilder string
	for i := 0; i < size; i++ {
		letter := input % uint64(a.length)
		strBuilder = string(a.alphabet[letter]) + strBuilder
		input /= uint64(a.length)
	}

	return strBuilder
}

func (a *Alphabet) H2i(hash []byte, y uint64) uint64 {
	buf := hash[:8]
	out := binary.LittleEndian.Uint64(buf)
	return (out + y) % uint64(a.possibilities)
}

func (a *Alphabet) I2i(input, y uint64, hashMethod HashType) uint64 {
	str := a.I2c(input)
	hash, _ := Hash(str, hashMethod)
	return a.H2i(hash, y)
}

func (a *Alphabet) NewChain(idx, width uint64, hashMethod HashType) uint64 {
	for i := uint64(1); i < width; i++ {
		idx = a.I2i(idx, i, hashMethod)
	}
	return idx
}
