package main

import (
	"encoding/asn1"
	"fmt"
	"os"
)

func main() {
	mdata, err := asn1.Marshal(13)
	checkErr(err)

	var n int
	_, err1 := asn1.Unmarshal(mdata, &n)
	checkErr(err1)

	fmt.Println("After marshal/unmarshal: ", n)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
