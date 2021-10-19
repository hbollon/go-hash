package main

import (
	"encoding/hex"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Welcome to go-hash project!")

	hashStr, err := HashToString("Salut", MD5)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("'Salut' in MD5: %s\n", hashStr)

	hashStr, err = HashToString("Salut", SHA1)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("'Salut' in SHA1: %s\n", hashStr)

	alphabetInstance := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 4, 5)
	logrus.Infof("Alphabet: %+v\n", alphabetInstance)
	logrus.Infof("I2C of 142678997: %s\n", alphabetInstance.I2c(142678997))

	alphabetInstance = GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4, 4)
	logrus.Infof("Alphabet: %+v\n", alphabetInstance)
	logrus.Infof("I2C of 1234: %s\n", alphabetInstance.I2c(1234))

	alphabetInstance = GenerateAlphabet("abcdefghijklmnopqrstuvwxyz", 4, 5)
	logrus.Infof("Alphabet: %+v\n", alphabetInstance)
	hash, err := Hash("oups", MD5)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("'oups' in MD5: %s\n", hex.EncodeToString(hash))
	logrus.Infof("H2I of \"oups\": %d\n", alphabetInstance.H2i(hash, 1))
	logrus.Infof("New string with width of 1000 and idx of 1234: %d\n", alphabetInstance.NewChain(1234, 1000))

	table := CreateRaindowTable(10, 5, alphabetInstance, MD5)
	logrus.Infof("Generate random rainbow table of 20*10: %s\n", spew.Sdump(table))
	err = table.Export("test.txt")
	if err != nil {
		logrus.Fatal(err)
	}
	table = CreateRaindowTable(20, 10, alphabetInstance, MD5)
	err = table.Import("test.txt")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Loaded table: %s\n", spew.Sdump(table))
}
