package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host:port")
		os.Exit(1)
	}
	url := os.Args[1]

	resp, err := http.Head(url)
	checkErr(err)

	fmt.Println("[response code]: ", resp.Status)
	for k, v := range resp.Header {
		fmt.Println(k, " : ", v)
	}
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
