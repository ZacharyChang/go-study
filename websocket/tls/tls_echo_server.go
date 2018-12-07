package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/", websocket.Handler(Echo))
	err := http.ListenAndServeTLS(":1200", "cert/jan.newmarch.name.pem", "key/private.pem", nil)
	checkErr(err)
}

func Echo(ws *websocket.Conn) {
	fmt.Println("Echoing")

	for n := 0; n < 10; n++ {
		msg := "hello " + string(n+48)
		fmt.Println("Sending:", msg)
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can't read:", err.Error())
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive:", err.Error())
			break
		}
		fmt.Println("Received:", reply)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
