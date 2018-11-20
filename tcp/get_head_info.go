package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

// Send the HEAD request to target host, and print the response
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	svc := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", svc)
	checkErr(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkErr(err)

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkErr(err)

	result, err := readFully(conn)
	checkErr(err)

	fmt.Println(string(result))
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
