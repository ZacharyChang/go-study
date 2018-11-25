package main

import (
	"encoding/json"
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
	loadJSON("person.json", &person)
	fmt.Println("Loading from file: \n", person)
}

func loadJSON(filename string, data interface{}) {
	inFile, err := os.Open(filename)
	checkErr(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(data)
	checkErr(err)

	inFile.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Print("Fatal error ", err.Error())
		os.Exit(1)
	}
}
