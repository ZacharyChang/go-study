package main

import (
	"fmt"
	"os"

	"github.com/zacharychang/go-study/websocket/xml/codec"
	"golang.org/x/net/websocket"
)

func main() {
	svc := "ws://localhost:1200"

	conn, err := websocket.Dial(svc, "", "http://localhost")
	checkErr(err)

	person := codec.Person{
		Name: "Jack",
		Emails: []string{
			"jack@gmail.com",
			"jack@outlook.com",
		},
	}
	err = codec.XMLCodec.Send(conn, person)
	checkErr(err)

	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}
