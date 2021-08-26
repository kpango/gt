package main

import (
	"fmt"

	"github.com/kpango/gt"
)

type Dummy struct {
	Text   string  `table:"text"`
	Number int     `table:"number"`
	Float  float64 `table:"float"`
	Child  *Dummy  `table:"child"`
}

type Child Dummy

var (
	dummyStruct = Dummy{
		Text:   "sample",
		Number: 12345,
		Float:  0.9876,
		Child: &Dummy{
			Text:   "I am child",
			Number: 12345,
			Float:  0.9876,
		},
	}
	dummyStructs = []Dummy{
		dummyStruct,
		dummyStruct,
		dummyStruct,
		dummyStruct,
		dummyStruct,
		dummyStruct,
	}
)

func main() {
	gt.Print(dummyStruct)
	fmt.Println("fuck")
	gt.Print(dummyStruct)
}
