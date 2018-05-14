package global

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"zeus/entity"
	"zeus/env"
	"zeus/events"
)

func TestSetGetGlobalStr(t *testing.T) {
	if !env.Load("..\\..\\server.json") {
		t.Fatalf("加载配置文件失败")
		return
	}

	g := GetGlobalInst()
	if g == nil {
		t.Fatalf("获取全局变量实例失败")
	}

	e := events.NewGlobalEventsInst()
	if e == nil {
		t.Fatalf("获取全局事件监听实例失败")
	}

	// Set
	g.SetGlobalStr("TestingStr", "Hello World!")
	time.Sleep(time.Second)
	e.HandleEvent()
	// Get
	str := g.GetGlobalStr("TestingStr")
	if str != "Hello World!" {
		t.Fatalf("GetGlobalStr(TestingStr) returned (%s), want (Hello World!)", str)
	}
	// Remove
	g.RemoveGlobal("TestingStr")
	time.Sleep(time.Second)
	e.HandleEvent()
	str = g.GetGlobalStr("TestingStr")
	if str != "" {
		t.Fatalf("After Remove, GetGlobalStr(TestingStr) returned (%s)", str)
	}

	fmt.Println("TestSetGetGlobalStr PASS")
}

func TestSetGetGlobalInt(t *testing.T) {
	if !env.Load("..\\..\\server.json") {
		t.Fatalf("加载配置文件失败")
		return
	}

	g := GetGlobalInst()
	if g == nil {
		t.Fatalf("获取全局变量实例失败")
	}

	e := events.NewGlobalEventsInst()
	if e == nil {
		t.Fatalf("获取全局事件监听实例失败")
	}

	// Set
	g.SetGlobalInt("TestingInt", 123)
	time.Sleep(time.Second)
	e.HandleEvent()
	// Get
	val := g.GetGlobalInt("TestingInt")
	if val != 123 {
		t.Fatalf("GetGlobalInt(TestingInt) returned (%d), want (123)", val)
	}
	// Remove
	g.RemoveGlobal("TestingInt")
	time.Sleep(time.Second)
	e.HandleEvent()
	val = g.GetGlobalInt("TestingStr")
	if val != 0 {
		t.Fatalf("After Remove, GetGlobalInt(TestingInt) returned (%d)", val)
	}
	fmt.Println("TestSetGetGlobalInt PASS")
}

func TestSetGetGlobalEntityProxy(t *testing.T) {
	if !env.Load("..\\..\\server.json") {
		t.Fatalf("加载配置文件失败")
		return
	}

	g := GetGlobalInst()
	if g == nil {
		t.Fatalf("获取全局变量实例失败")
	}

	e := events.NewGlobalEventsInst()
	if e == nil {
		t.Fatalf("获取全局事件监听实例失败")
	}

	entity2Set := new(entity.EntityProxy)
	entity2Set.SrvID = 1000
	entity2Set.EntityID = 100
	// Set
	g.SetGlobalEntityProxy("TestingEntityProxy", entity2Set)
	time.Sleep(time.Second)
	e.HandleEvent()
	// Get
	entity2Get := g.GetGlobalEntityProxy("TestingEntityProxy")
	if !reflect.DeepEqual(entity2Set, entity2Get) {
		t.Fatalf("GetGlobalEntityProxy(TestingEntityProxy) return (%v), want (%v)", entity2Get, entity2Set)
	}
	// Remove
	g.RemoveGlobal("TestingEntityProxy")
	time.Sleep(time.Second)
	e.HandleEvent()
	entity2Get = g.GetGlobalEntityProxy("TestingEntityProxy")
	if entity2Get != nil {
		t.Fatalf("After Remove, GetGlobalEntityProxy(TestingEntityProxy) return (%v)", entity2Get)
	}

	fmt.Println("TestSetGetGlobalEntityProxy PASS")
}
