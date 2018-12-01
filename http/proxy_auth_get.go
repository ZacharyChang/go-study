package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const auth = "user:password"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "http://host:port", "http://proxy-host:proxy:port")
		os.Exit(1)
	}

	targetUrl, err := url.Parse(os.Args[1])
	checkErr(err, 2)

	proxyStr := os.Args[2]
	proxyUrl, err := url.Parse(proxyStr)
	checkErr(err, 3)

	// encode the auth
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	req, err := http.NewRequest("GET", targetUrl.String(), nil)
	// add auth header to proxy server
	req.Header.Add("Proxy-Authorization", authHeader)

	dump, _ := httputil.DumpRequest(req, false)
	fmt.Println(string(dump))

	resp, err := client.Do(req)
	checkErr(err, 4)

	fmt.Println("Read done")
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		os.Exit(5)
	}

	var buf [512]byte
	for {
		n, err := resp.Body.Read(buf[0:])
		fmt.Print(string(buf[0:n]))
		if err != nil {
			fmt.Println()
			os.Exit(0)
		}
	}
}

func checkErr(err error, code int) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(code)
	}
}
