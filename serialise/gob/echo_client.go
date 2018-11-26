package main

import (
	"encoding/gob"
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
	person := Person{
		Name: Name{
			Family: "Smith",
			Last:   "John",
		},
		Email: []Email{
			{
				Kind:    "home",
				Address: "john@home.org",
			},
			{
				Kind:    "work",
				Address: "john@gmail.com",
			},
		},
	}
	svc := "localhost:1200"

	conn, err := net.Dial("tcp", svc)
	checkErr(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for n := 0; n < 10; n++ {
		encoder.Encode(person)
		var newPerson Person
		decoder.Decode(&newPerson)
		fmt.Println(n, " ", newPerson.String())
	}

	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
