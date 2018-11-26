package main

import (
	"encoding/gob"
	"fmt"
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
	var person Person
	loadGob("person.gob", &person)
	fmt.Println("Person", person.String())
}

func loadGob(filename string, data interface{}) {
	infile, err := os.Open(filename)
	checkErr(err)

	decoder := gob.NewDecoder(infile)
	err = decoder.Decode(data)
	checkErr(err)

	infile.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
