package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/zacharychang/go-study/websocket/xml/codec"
	"golang.org/x/net/websocket"
)

func ReceivePerson(ws *websocket.Conn) {
	var person codec.Person
	err := codec.XMLCodec.Receive(ws, &person)
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
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}
