package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Welcome to go-hash project!")
	// alphabetInstance := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 4, 5)
	// spew.Dump(alphabetInstance)
	// input := uint64(1)
	// for i := 1; i < 100; i++ {
	// 	input = alphabetInstance.I2i(input, uint64(i))
	// 	fmt.Println(input)
	// }

	// hashStr, err := HashToString("Salut", MD5)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// logrus.Infof("'Salut' in MD5: %s\n", hashStr)

	// hashStr, err = HashToString("Salut", SHA1)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// logrus.Infof("'Salut' in SHA1: %s\n", hashStr)

	// alphabetInstance := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 4, 5)
	// logrus.Infof("Alphabet: %+v\n", alphabetInstance)
	// logrus.Infof("I2C of 142678997: %s\n", alphabetInstance.I2c(142678997))

	// alphabetInstance = GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4, 4)
	// logrus.Infof("Alphabet: %+v\n", alphabetInstance)
	// logrus.Infof("I2C of 1234: %s\n", alphabetInstance.I2c(1234))

	// alphabetInstance = GenerateAlphabet("abcdefghijklmnopqrstuvwxyz", 4, 5)
	// logrus.Infof("Alphabet: %+v\n", alphabetInstance)
	// hash, err := Hash("oups", MD5)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// logrus.Infof("'oups' in MD5: %s\n", hex.EncodeToString(hash))
	// logrus.Infof("H2I of \"oups\": %d\n", alphabetInstance.H2i(hash, 1))
	// logrus.Infof("New string with width of 1000 and idx of 1234: %d\n", alphabetInstance.NewChain(1234, 1000))

	// table := CreateRaindowTable(10, 5, alphabetInstance, MD5)
	// logrus.Infof("Generate random rainbow table of 20*10: %s\n", spew.Sdump(table))
	// err = table.Export("test.txt")
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// table = CreateRaindowTable(20, 10, alphabetInstance, MD5)
	// err = table.Import("test.txt")
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// logrus.Info("Loaded table:")
	// table.Print()

	alphabetInstance := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4, 4)
	table := CreateRaindowTable(100000, 1000, alphabetInstance, MD5)
	spew.Dump(alphabetInstance)
	table.Print()
	//logrus.Infof("Generate random rainbow table of 100000*1000: %s\n", spew.Sdump(table))
	hash, _ := Hash("ABCD", MD5)
	if out, err := table.Invert(hash); out == "ABCD" && err == nil {
		logrus.Info("Invert success!")
	} else {
		logrus.Errorf("Failed, output: %s, error: %s", out, err)
	}
	table.Stats()
}
