package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// the key and cert file must be matched
	cert, err := tls.LoadX509KeyPair("cert/jan.newmarch.name.pem", "key/private.pem")
	checkErr(err)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	now := time.Now()
	config.Time = func() time.Time {
		return now
	}
	config.Rand = rand.Reader

	svc := "localhost:1200"
	listener, err := tls.Listen("tcp", svc, &config)
	checkErr(err)
	fmt.Println("Listening on: ", svc)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Accepted")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		fmt.Println("Reading...")
		n, err := conn.Read(buf[0:])
		if err != nil {
			_ = fmt.Errorf("error: %s\n", err.Error())
		}
		if n == 0 {
			break
		}
		fmt.Println(string(buf[0:n]))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			_ = fmt.Errorf("error: %s\n", err2.Error())
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
