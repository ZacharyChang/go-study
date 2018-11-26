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
	saveGob("person.gob", person)
}

func saveGob(filename string, data interface{}) {
	outfile, err := os.Create(filename)
	checkErr(err)

	encoder := gob.NewEncoder(outfile)
	err = encoder.Encode(data)
	checkErr(err)

	outfile.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
