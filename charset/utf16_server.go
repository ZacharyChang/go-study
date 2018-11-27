package main

import (
	"fmt"
	"net"
	"os"
	"unicode/utf16"
)

const BOM = '\ufffe'

func main() {
	svc := "localhost:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", svc)
	checkErr(err)

	listner, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}
		str := "你好,hello,こんにちは,สวัสดี,привет"
		shorts := utf16.Encode([]rune(str))
		writeShorts(conn, shorts)
		conn.Close()
	}
}

func writeShorts(conn net.Conn, shorts []uint16) {
	var bytes [2]byte
	bytes[0] = BOM >> 8
	bytes[1] = BOM & 255
	_, err := conn.Write(bytes[0:])
	if err != nil {
		return
	}
	for _, v := range shorts {
		bytes[0] = byte(v >> 8)
		bytes[1] = byte(v & 255)
		_, err = conn.Write(bytes[0:])
		if err != nil {
			return
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
