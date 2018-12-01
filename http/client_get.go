package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "http://host:port")
		os.Exit(1)
	}
	inputUrl, err := url.Parse(os.Args[1])
	checkErr(err, 2)

	client := &http.Client{}
	req, err := http.NewRequest("GET", inputUrl.String(), nil)
	req.Header.Add("Accept-Charset", "UTF-8;q=1,ISO-8859-1;q=0")
	checkErr(err, 3)

	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		os.Exit(4)
	}

	chSet := getCharset(resp)
	fmt.Println("Parse response charset: ", chSet)
	if chSet != "UTF-8" {
		fmt.Println("Not UTF-8 encoding", chSet)
		os.Exit(5)
	}

	var buf [512]byte
	for {
		n, err := resp.Body.Read(buf[0:])
		if err != nil {
			continue
		}
		if n == 0 {
			os.Exit(0)
		}
		fmt.Print(string(buf[0:n]))
	}
}

func getCharset(resp *http.Response) string {
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// default ad utf-8
		return "UTF-8"
	}
	idx := strings.Index(contentType, "charset:")
	if idx == -1 {
		return "UTF-8"
	}
	return strings.Trim(contentType[idx:], " ")
}

func checkErr(err error, code int) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(code)
	}
}
