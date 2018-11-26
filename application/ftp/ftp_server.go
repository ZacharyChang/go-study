package main

import (
	"fmt"
	"net"
	"os"
)

const (
	DIR = "DIR"
	CD  = "CD"
	PWD = "PWD"
)

func main() {
	svc := "localhost:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", svc)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			conn.Close()
			return
		}
		s := string(buf[0:n])
		if s[0:2] == CD {
			chdir(conn, s[3:])
		} else if s[0:3] == DIR {
			dirList(conn)
		} else if s[0:3] == PWD {
			pwd(conn)
		}
	}
}

func chdir(conn net.Conn, dir string) {
	if os.Chdir(dir) == nil {
		conn.Write([]byte("OK"))
	} else {
		conn.Write([]byte("ERROR"))
	}
}

func dirList(conn net.Conn) {
	defer conn.Write([]byte("\r\n"))

	dir, err := os.Open(".")
	if err != nil {
		return
	}
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, name := range names {
		conn.Write([]byte(name + "\r\n"))
	}
}

func pwd(conn net.Conn) {
	s, err := os.Getwd()
	if err != nil {
		conn.Write([]byte("ERROR"))
		return
	}
	conn.Write([]byte(s))
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
