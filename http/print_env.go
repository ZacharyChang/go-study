package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fileServer := http.FileServer(http.Dir("/Users/zachary/Downloads"))
	http.Handle("/", fileServer)

	http.HandleFunc("/cgi-bin/printenv", printEnv)

	err := http.ListenAndServe(":1200", nil)
	checkErr(err)
}

func printEnv(writer http.ResponseWriter, req *http.Request) {
	env := os.Environ()
	_, _ = writer.Write([]byte("<h1>Environment</h1>\n<pre>"))
	for _, v := range env {
		_, _ = writer.Write([]byte(v + "\n"))
	}
	_, _ = writer.Write([]byte("</pre>"))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
