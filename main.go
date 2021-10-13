package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("'Salut' in MD5: " + HashMD5("Salut"))
	fmt.Println("'Salut' in SHA1: " + HashSHA1("Salut"))
}
