package main

import (
	"fmt"
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

func greet(name string) {
	fmt.Printf("hello, %s!\n", name)
}

func main() {
	vm := otto.New()

	data, err := ioutil.ReadFile("mermaid.min.js")
	if err != nil {
		panic(err)
	}
	strData := string(data)

	this, err := vm.Run(strData)

	fmt.Println(this, err)
}
