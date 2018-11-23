package main

import (
	"encoding/asn1"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	svc := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", svc)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		daytime := time.Now()

		mdata, _ := asn1.Marshal(daytime)
		fmt.Println("Send time: ", daytime.String())
		conn.Write(mdata)
		conn.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
