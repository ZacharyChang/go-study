package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Daytime server
func main() {
	svc := ":1248"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", svc)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		daytime := time.Now().String()
		conn.Write([]byte(daytime + "\n"))
		conn.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
