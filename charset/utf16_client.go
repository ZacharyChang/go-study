package main

import (
	"fmt"
	"net"
	"os"
	"unicode/utf16"
)

//const BOM = '\ufffe'

func main() {
	svc := "localhost:1200"

	conn, err := net.Dial("tcp", svc)
	checkErr(err)

	shorts := readShors(conn)
	ints := utf16.Decode(shorts)
	str := string(ints)
	fmt.Println(str)
	os.Exit(0)
}

func readShors(conn net.Conn) []uint16 {
	var buf [512]byte
	n, err := conn.Read(buf[0:2])
	for true {
		m, err := conn.Read(buf[n:])
		if m == 0 || err != nil {
			break
		}
		n += m
	}
	checkErr(err)

	var shorts []uint16
	shorts = make([]uint16, n/2)
	if buf[0] == 0xff && buf[1] == 0xfe {
		// big endian
		for i := 2; i < n; i += 2 {
			shorts[i/2] = uint16(buf[i])<<8 + uint16(buf[i+1])
		}
	} else if buf[1] == 0xff && buf[0] == 0xfe {
		// little endian
		for i := 2; i < n; i += 2 {
			shorts[i/2] = uint16(buf[i+1])<<8 + uint16(buf[i])
		}
	} else {
		// unknown byte order
		fmt.Println("Unknown byte order")
	}
	return shorts
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
