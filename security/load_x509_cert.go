package main

import (
	"crypto/x509"
	"fmt"
	"os"
)

func main() {
	certCerFile, err := os.Open("cert/jan.newmarch.name.cer")
	checkErr(err)

	derBytes := make([]byte, 1000)
	count, err := certCerFile.Read(derBytes)
	checkErr(err)
	certCerFile.Close()

	cert, err := x509.ParseCertificate(derBytes[0:count])
	checkErr(err)

	fmt.Printf("Name %s\n", cert.Subject.CommonName)
	fmt.Printf("Not before %s\n", cert.NotBefore.String())
	fmt.Printf("Not after %s\n", cert.NotAfter.String())
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
