package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"os"
)

func main() {
	svc := "ws://localhost:1200"
	conn, err := websocket.Dial(svc, "", "http://localhost")
	checkErr(err)
	var msg string
	for {
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
