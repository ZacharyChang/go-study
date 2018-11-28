package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	hash := md5.New()
	bytes := []byte("hello\n")
	hash.Write(bytes)
	hashVal := hash.Sum(nil)
	hashSize := hash.Size()

	for n := 0; n < hashSize; n += 4 {
		val := uint32(hashVal[n])<<24 + uint32(hashVal[n+1])<<16 + uint32(hashVal[n+2])<<8 + uint32(hashVal[n+3])
		fmt.Printf("%x", val)
	}
	fmt.Println()
}
