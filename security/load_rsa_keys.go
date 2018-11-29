package main

import (
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	var key rsa.PrivateKey
	loadKey("key/private.key", &key)

	fmt.Println("Private key primes: ", key.Primes[0].String(), key.Primes[1].String())
	fmt.Println("Private key exponent: ", key.D.String())

	var publicKey rsa.PublicKey
	loadKey("key/public.key", &publicKey)

	fmt.Println("Public key modulus: ", publicKey.N.String())
	fmt.Println("Public key exponent", publicKey.E)
}

func loadKey(filename string, key interface{}) {
	infile, err := os.Open(filename)
	checkErr(err)

	decoder := gob.NewDecoder(infile)
	err = decoder.Decode(key)
	checkErr(err)

	infile.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
