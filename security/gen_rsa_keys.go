package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	reader := rand.Reader
	bitSize := 512
	key, err := rsa.GenerateKey(reader, bitSize)
	checkErr(err)

	fmt.Println("Private key primes", key.Primes[0].String(), key.Primes[1].String())
	fmt.Println("Private key exponent", key.D.String())

	publicKey := key.PublicKey
	fmt.Println("Public key modulus", publicKey.N.String())
	fmt.Println("Public key exponent", publicKey.E)

	saveGobKey("key/private.key", key)
	saveGobKey("key/public.key", publicKey)
	savePemKey("key/private.pem", key)
}

func savePemKey(filename string, key *rsa.PrivateKey) {
	outfile, err := os.Create(filename)
	checkErr(err)

	privateKey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	pem.Encode(outfile, privateKey)
	outfile.Close()
}

func saveGobKey(filename string, key interface{}) {
	outfile, err := os.Create(filename)
	checkErr(err)

	encoder := gob.NewEncoder(outfile)
	err = encoder.Encode(key)
	checkErr(err)

	outfile.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
