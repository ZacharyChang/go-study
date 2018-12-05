package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

type Person struct {
	Name   string
	Emails []string
}

func ReceivePerson(ws *websocket.Conn) {
	var person Person
	err := websocket.JSON.Receive(ws, &person)
	if err != nil {
		fmt.Println("Can't receive:", err.Error())
		return
	}
	fmt.Println("Name:", person.Name)
	for _, v := range person.Emails {
		fmt.Println("Email:", v)
	}
}

func main() {
	http.Handle("/", websocket.Handler(ReceivePerson))
	err := http.ListenAndServe("localhost:1200", nil)
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}
