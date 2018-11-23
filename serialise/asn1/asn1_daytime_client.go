package main

import (
	"bytes"
	"encoding/asn1"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	svc := "localhost:1200"

	conn, err := net.Dial("tcp", svc)
	checkErr(err)

	res, err := readFully(conn)
	checkErr(err)

	var newtime time.Time
	_, err1 := asn1.Unmarshal(res, &newtime)
	checkErr(err1)

	fmt.Println("After marshal/unmarshal: ", newtime.String())

	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	res := bytes.NewBuffer(nil)
	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])
		res.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return res.Bytes(), nil
}
