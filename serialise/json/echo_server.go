package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Person struct {
	Name  Name
	Email []Email
}

type Name struct {
	Family string
	Last   string
}

type Email struct {
	Kind    string
	Address string
}

func (p Person) String() string {
	s := p.Name.Last + " " + p.Name.Family
	for _, v := range p.Email {
		s += "\n" + v.Kind + ":" + v.Address
	}
	return s
}

func main() {
	svc := "localhost:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", svc)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		encoder := json.NewEncoder(conn)
		decoder := json.NewDecoder(conn)

		for n := 0; n < 10; n++ {
			var person Person
			decoder.Decode(&person)
			fmt.Println(n, " ", person.String())
			encoder.Encode(person)
		}
		conn.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
