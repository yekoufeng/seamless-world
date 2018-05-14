package events

import (
	"fmt"
	"testing"
)

type forTest1 struct {
}

func (f *forTest1) Proc1(args string, flag bool, index uint64) {
	fmt.Println("Proc1:", args, flag, index)
}

func (f *forTest1) Proc2(args []byte) {
	fmt.Println("Proc2:", string(args))
}

type forTest2 struct {
}

func (f *forTest2) Proc3(args []byte) {
	fmt.Println("Proc3:", string(args))
}

func (f *forTest2) Proc4(args []byte) {
	fmt.Println("Proc4:", string(args))
}

func TestEvents(t *testing.T) {
	e := NewEventsInst()
	obj := new(forTest1)

	e.AddListener("Testing", obj, "Proc1")
	e.FireEvent("Testing")
}
