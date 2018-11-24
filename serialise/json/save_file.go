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
	saveJSON("person.json", person)

}

func saveJSON(filename string, data interface{}) {
	out, err := os.Create(filename)
	checkErr(err)

	encoder := json.NewEncoder(out)
	err = encoder.Encode(data)
	checkErr(err)
	out.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
