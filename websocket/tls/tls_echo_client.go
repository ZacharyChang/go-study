package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	svc := "wss://localhost:1200"
	config, err := websocket.NewConfig(svc, "http://localhost")
	checkErr(err)

	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := websocket.DialConfig(config)
	checkErr(err)
	var msg string

	for {
		conn.Config().TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Couldn't receive msg:", err.Error())
			break
		}
		fmt.Println("Received msg:", msg)
		err = websocket.Message.Send(conn, msg)
		if err != nil {
			fmt.Println("Couldn't send msg")
			break
		}
	}
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}
