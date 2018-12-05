package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"
)

type Person struct {
	Name   string
	Emails []string
}

func main() {
	svc := "ws://localhost:1200"

	conn, err := websocket.Dial(svc, "", "http://localhost")
	checkErr(err)

	person := Person{
		Name: "Jack",
		Emails: []string{
			"jack@gmail.com",
			"jack@outlook.com",
		},
	}
	err = websocket.JSON.Send(conn, person)
	checkErr(err)

	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}
