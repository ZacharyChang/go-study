package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Person struct {
	XMLName Name    `xml:"person"`
	Name    Name    `xml:"name"`
	Email   []Email `xml:"email"`
}

type Name struct {
	Family   string `xml:"family"`
	Personal string `xml:"personal"`
}

type Email struct {
	Type    string `xml:"type,attr"`
	Address string `xml:",chardata"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "file")
		os.Exit(1)
	}
	file := os.Args[1]
	bytes, err := ioutil.ReadFile(file)
	checkErr(err)

	var person Person
	err = xml.Unmarshal(bytes, &person)
	checkErr(err)

	fmt.Println(person.String())
}

func (person *Person) String() string {
	var buf bytes.Buffer
	buf.WriteString("Family Name: " + person.Name.Family + "\n")
	buf.WriteString("Personal Name: " + person.Name.Personal + "\n")
	for i, v := range person.Email {
		buf.WriteString("[Email " + strconv.Itoa(i+1) + "]\n")
		buf.WriteString("Type: " + v.Type + "\n")
		buf.WriteString("Address: " + strings.TrimSpace(v.Address) + "\n")
	}
	return buf.String()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}
