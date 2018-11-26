package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	uiDir  = "dir"
	uiCd   = "cd"
	uiPwd  = "pwd"
	uiQuit = "quit"
)

const (
	DIR = "DIR"
	CD  = "CD"
	PWD = "PWD"
)

func main() {
	svc := "localhost:1200"
	conn, err := net.Dial("tcp", svc)
	checkErr(err)

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, "\t\r\n")
		if err != nil {
			break
		}
		// split into cammnd + arg
		strs := strings.SplitN(line, " ", 2)
		// decoder request
		switch strs[0] {
		case uiDir:
			dirRequest(conn)
		case uiCd:
			if len(strs) != 2 {
				fmt.Println("cd <dir>")
				continue
			}
			fmt.Println("CD \"", strs[1], "\"")
			cdRequest(conn, strs[1])
		case uiPwd:
			pwdRequest(conn)
		case uiQuit:
			conn.Close()
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
		}
	}
}

func dirRequest(conn net.Conn) {
	conn.Write([]byte(DIR + " "))
	var buf [512]byte
	res := bytes.NewBuffer(nil)
	for {
		n, _ := conn.Read(buf[0:])
		res.Write(buf[0:n])
		len := res.Len()
		contents := res.Bytes()
		if string(contents) == "\r\n" {
			fmt.Println("Empty folder")
			return
		}
		if string(contents[len-4:]) == "\r\n\r\n" {
			fmt.Println(string(contents[0 : len-4]))
			return
		}
	}
}

func cdRequest(conn net.Conn, dir string) {
	conn.Write([]byte(CD + " " + dir))
	var resp [512]byte
	n, _ := conn.Read(resp[0:])
	s := string(resp[0:n])
	if s != "OK" {
		fmt.Println("Failed to change dir")
	}
}

func pwdRequest(conn net.Conn) {
	conn.Write([]byte(PWD))
	var resp [512]byte
	n, _ := conn.Read(resp[0:])
	s := string(resp[0:n])
	fmt.Println("Current dir \"", s, "\"")
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
