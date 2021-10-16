package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
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

func (a *Alphabet) I2c(input int) string {
	coeff := len(a.alphabet)
	for input >= coeff {
		input = input - coeff
		coeff *= len(a.alphabet)
	}

	var str string
	for input > len(a.alphabet) {
		letter := input % len(a.alphabet)
		fmt.Println(letter)
		str = string(a.alphabet[letter]) + str
		input = input / len(a.alphabet)
	}

	return str
}
