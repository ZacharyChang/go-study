package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

func Echo(ws *websocket.Conn) {
	fmt.Println("Echoing...")

	for i := 0; i < 10; i++ {
		msg := "Hello " + string(i+48)
		fmt.Println("Sending to client: " + msg)
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can't send")
			break
		}
		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive")
		}
		fmt.Println("Received from client:", reply)
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))
	err := http.ListenAndServe("localhost:1200", nil)
	if err != nil {
		fmt.Println("Server error")
		os.Exit(1)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
