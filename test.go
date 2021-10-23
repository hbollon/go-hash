package gohash

import "github.com/sirupsen/logrus"

// All tests case prepared for the TUI's test menu

func TestHash() {
	logrus.Info("This test will execute Hash function on \"Salut\", \"Bob\" and \"Cryptologie\" with both MD5 and SHA1 methods.")
	logrus.Info("The result will be displayed in the terminal.\n\n")

	cases := []string{"Salut", "Bob", "Cryptologie"}
	for i := 0; i < len(cases); i++ {
		hashStr, err := HashToString(cases[i], MD5)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("%s in MD5: %s\n", cases[i], hashStr)

		hashStr, err = HashToString(cases[i], SHA1)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("%s in SHA1: %s\n\n", cases[i], hashStr)
	}
}

func TestI2c() {
	logrus.Info("This test will execute I2c function on 150106454, 75324, 1651 and 4173921 using: \"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz\" alphabet, 1 as min size and 6 as max size.")
	logrus.Info("The result will be displayed in the terminal.\n\n")

	cases := []uint64{150106454, 75324, 1651, 4173921}
	alph := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 1, 6)
	for i := 0; i < len(cases); i++ {
		res := alph.I2c(cases[i])
		logrus.Infof("%d's i2c result: %s\n", cases[i], res)
	}
}

func TestH2i() {
	logrus.Info("This test will execute H2i function on \"72eb471fb3bd65c03d29f2fcbb9984d6\" (\"oups\" MD5), 75324, 1651 and 4173921 using: \"abcdefghijklmnopqrstuvwxyz\" alphabet, 4 as min size and 5 as max size.")
	logrus.Info("The result will be displayed in the terminal.\n\n")

	alph := GenerateAlphabet("abcdefghijklmnopqrstuvwxyz", 4, 5)
	hash, err := Hash("oups", MD5)
	if err != nil {
		logrus.Error(err)
	}
	res := alph.H2i(hash, 1)
	logrus.Infof("\"oups\" MD5's h2i result: %d\n", res)
}

func TestI2i() {
	logrus.Info("This test will execute I2i function (with MD5 hash method) on 1 using: \"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz\" alphabet, 4 as min size and 5 as max size.")
	logrus.Info("The result will be displayed in the terminal.\n\n")

	alph := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 4, 5)
	res := alph.I2i(1, 1, MD5)
	logrus.Infof("1's i2i result: %d\n", res)
}

func TestNewChain() {
	logrus.Info("This test will execute NewChain (successive i2i calls) function on 1 using: \"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz\" alphabet, 4 as min size, 5 as max size and 100 as string width.")
	logrus.Info("The result will be displayed in the terminal.\n\n")

	alph := GenerateAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 4, 5)
	res := alph.NewChain(1, 100, MD5)
	logrus.Infof("1's NewChain result: %d\n", res)
}
