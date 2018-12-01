package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host:port")
		os.Exit(1)
	}
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		os.Exit(3)
	}
	dump, _ := httputil.DumpResponse(resp, false)
	fmt.Println(string(dump))

	contentType := resp.Header["Content-Type"]
	if !acceptableCharset(contentType) {
		fmt.Println("Cannot handle: ", contentType)
	}

	var buf [512]byte
	reader := resp.Body
	for {
		n, err := reader.Read(buf[0:])
		if err != nil {
			fmt.Println(err.Error())
		}
		if n == 0 {
			os.Exit(0)
		}
		fmt.Println(string(buf[0:n]))
	}
}

func acceptableCharset(contentType []string) bool {
	for _, cType := range contentType {
		if strings.Contains(cType, "UTF-8") {
			return true
		}
	}
	return false
}
