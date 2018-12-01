package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	random := rand.Reader

	var key rsa.PrivateKey
	loadKey("key/private.key", &key)

	now := time.Now()
	then := now.Add(60 * 60 * 24 * 365 * 1000 * 1000 * 1000)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "jan.newmarch.name",
			Organization: []string{
				"Jan Newmarch",
			},
		},
		NotBefore:             now,
		NotAfter:              then,
		SubjectKeyId:          []byte{1, 2, 3, 4},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:     true,
		DNSNames: []string{"jan.newmarch.name", "localhost"},
	}
	derBytes, err := x509.CreateCertificate(random, &template, &template, &key.PublicKey, &key)
	checkErr(err)

	certCerFile, err := os.Create("cert/jan.newmarch.name.cer")
	checkErr(err)
	_, _ = certCerFile.Write(derBytes)
	certCerFile.Close()

	certPemFile, err := os.Create("cert/jan.newmarch.name.pem")
	checkErr(err)
	_ = pem.Encode(certPemFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	certPemFile.Close()

	keyPemFile, err := os.Create("key/private.pem")
	checkErr(err)
	_ = pem.Encode(keyPemFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(&key),
	})
	keyPemFile.Close()
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
