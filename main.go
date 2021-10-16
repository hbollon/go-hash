package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Welcome to go-hash project!")

	hash, err := Hash("Salut", MD5)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("'Salut' in MD5: %s\n", hash)

	hash, err = Hash("Salut", SHA1)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("'Salut' in SHA1: %s\n", hash)

	alphabetInstance := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4, 4)
	fmt.Printf("Alphabet: %+v\n", alphabetInstance)

	fmt.Printf("I2C of 142678997: %s\n", alphabetInstance.I2c(142678997))
}
