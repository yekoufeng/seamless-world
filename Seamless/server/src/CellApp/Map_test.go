package main

import (
	"fmt"
	"testing"
)

type ii interface {
	test()
}

type AA struct {
	a int
}

func (a *AA) test() {

}

func Test(t *testing.T) {

	var a interface{}

	b := 1
	c := "XXXX"

	a = b

	fmt.Println(a)
	fmt.Println(a == c)

	a = c

	fmt.Println(a)
	fmt.Println(a == c)
}
