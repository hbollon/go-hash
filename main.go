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
}
