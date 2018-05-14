package events

import (
	"fmt"
	"testing"
	"time"
	"zeus/env"
)

func TestGlobalEvents(t *testing.T) {
	if !env.Load("..\\..\\..\\server.json") {
		fmt.Println("加载配置文件失败")
		return
	}

	e := NewGlobalEventsInst()
	obj := new(forTest1)

	e.AddListener("Testing", obj, "Proc1")
	e.FireEvent("Testing", "Testing Add and Fire", true, uint64(15))
	time.Sleep(time.Second)
	e.HandleEvent()
}

func ExampleGlobalEvents() {
	if !env.Load("..\\..\\..\\server.json") {
		fmt.Println("加载配置文件失败")
		return
	}

	e := NewGlobalEventsInst()
	obj := new(forTest1)

	h1 := e.AddListener("Testing", obj, "Proc1")
	e.FireEvent("Testing", []byte("Testing Add and Fire"))
	e.AddListener("TestingOne", obj, "Proc2")
	e.FireEvent("TestingOne", []byte("Testing Add and Fire"))
	time.Sleep(time.Second)
	e.HandleEvent()

	objOne := new(forTest2)
	e.AddListener("Testing", objOne, "Proc3")
	e.FireEvent("Testing", []byte("Testing Add and Fire"))
	e.AddListener("TestingOne", objOne, "Proc4")
	e.FireEvent("TestingOne", []byte("Testing Add and Fire"))
	time.Sleep(time.Second)
	e.HandleEvent()

	e.RemoveListenerByEvtHandle(h1)
	e.FireEvent("Testing", []byte("Testing Add and Fire"))
	time.Sleep(time.Second)
	e.HandleEvent()
	e.RemoveListenerByObjInst(objOne)
	e.FireEvent("Testing", []byte("Testing Add and Fire"))
	e.FireEvent("TestingOne", []byte("Testing Add and Fire"))
	time.Sleep(time.Second)
	e.HandleEvent()
	// Output:
	// Proc1: Testing Add and Fire
	// Proc2: Testing Add and Fire
	// Proc1: Testing Add and Fire
	// Proc3: Testing Add and Fire
	// Proc2: Testing Add and Fire
	// Proc4: Testing Add and Fire
	// Proc3: Testing Add and Fire
	// Proc2: Testing Add and Fire
}
