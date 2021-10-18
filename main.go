package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Welcome to go-hash project!")

	hashStr, err := HashToString("Salut", MD5)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("'Salut' in MD5: %s\n", hashStr)

	hashStr, err = HashToString("Salut", SHA1)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("'Salut' in SHA1: %s\n", hashStr)

	alphabetInstance := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 4, 5)
	fmt.Printf("Alphabet: %+v\n", alphabetInstance)
	fmt.Printf("I2C of 142678997: %s\n", alphabetInstance.I2c(142678997))

	alphabetInstance = GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4, 4)
	fmt.Printf("Alphabet: %+v\n", alphabetInstance)
	fmt.Printf("I2C of 1234: %s\n", alphabetInstance.I2c(1234))

	alphabetInstance = GenerateAlphabet("abcdefghijklmnopqrstuvwxyz", 4, 5)
	fmt.Printf("Alphabet: %+v\n", alphabetInstance)
	hash, err := Hash("oups", MD5)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("'oups' in MD5: %s\n", hex.EncodeToString(hash))
	fmt.Printf("H2I of \"oups\": %d\n", alphabetInstance.H2i(hash, 1))
}
