package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage", os.Args[0], "http://proxy-host:proxy-port", " http://host:port")
		os.Exit(1)
	}
	proxyStr := os.Args[1]
	proxyUrl, err := url.Parse(proxyStr)
	checkErr(err, 2)

	targetUrlStr := os.Args[2]
	targetUrl, err := url.Parse(targetUrlStr)
	checkErr(err, 3)

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", targetUrl.String(), nil)
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
