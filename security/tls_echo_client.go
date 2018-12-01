package main

import (
	"crypto/tls"
	"fmt"
	"os"
)

func main() {
	svc := "localhost:1200"
	conn, err := tls.Dial("tcp", svc, &tls.Config{
		// skip the ca verify
		InsecureSkipVerify: true,
	})
	checkErr(err)

	for n := 0; n < 10; n++ {
		fmt.Println("Writing...")
		_, _ = conn.Write([]byte("Hello " + string(n+48)))

		var buf [512]byte
		n, err := conn.Read(buf[0:])
		checkErr(err)
		fmt.Println(string(buf[0:n]))
	}
	fmt.Println("Writing done.")
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
