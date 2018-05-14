package msghandler

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"zeus/safecontainer"
	"zeus/serializer"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

/*
	MsgHandler 作为底层通讯层与上层应用层之间逻辑传递的桥梁
*/

// IMsgHandlers 消息处理模块的接口
type IMsgHandlers interface {
	RegMsgProc(proc interface{})
	FireMsg(name string, content interface{})
	FireRPC(methodName string, data []byte)

	DoNormalMsg(string, interface{}) error
	DoRPCMsg(string, []byte) error

	DoMsg()
	SetEnable(enable bool)
}

// NewMsgHandlers 创建一个新的消息处理器
func NewMsgHandlers() IMsgHandlers {

	return &MsgHandlers{
		msgFuncs: &sync.Map{},
		rpcFuncs: &sync.Map{},

		enable: true,
		//msgCount:     make(map[string][]int64),
		msgFireInfo:  safecontainer.NewSafeList(),
		rpcFireInfo:  safecontainer.NewSafeList(),
		defaultFuncs: make([]reflect.Value, 10),
	}
}

//msgFireInfo 触发一个事件时必须的信息
type msgFireInfo struct {
	name    string
	content interface{}
}

type rpcFireInfo struct {
	methodName string
	data       []byte
}

// MsgHandlers 消息处理中心
type MsgHandlers struct {
	msgFuncs *sync.Map
	rpcFuncs *sync.Map

	enable       bool
	defaultFuncs []reflect.Value

	msgFireInfo *safecontainer.SafeList
	rpcFireInfo *safecontainer.SafeList
}

// SetEnable 临时方法
func (handlers *MsgHandlers) SetEnable(enable bool) {
	handlers.enable = enable
}

// RegMsgProc 注册消息处理对象
// 其中 proc 是一个对象，包含是类似于 MsgProc_XXXXX的一系列函数，分别用来处理不同的消息
func (handlers *MsgHandlers) RegMsgProc(proc interface{}) {

	v := reflect.ValueOf(proc)
	t := reflect.TypeOf(proc)

	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		msgName, msgHandler, err := handlers.getMsgHandler(methodName, v.MethodByName(methodName))
		if err == nil {
			handlers.addMsgHandler(msgName, msgHandler)
			continue
		}

		// 判断是否是RPC处理函数
		msgName, msgHandler, err = handlers.getRPCHandler(methodName, v.MethodByName(methodName))
		if err == nil {
			handlers.addRPCHandler(msgName, msgHandler)
		}
	}

	handlers.addDefaultFunc(proc)
}

func (handlers *MsgHandlers) addDefaultFunc(proc interface{}) {

	v := reflect.ValueOf(proc)
	var defaultFunc = v.MethodByName("MsgProc_DefaultMsgHandler")
	if defaultFunc.IsValid() {
		handlers.defaultFuncs = append(handlers.defaultFuncs, defaultFunc)
	}
}

// FireMsg 触发消息, 保证不被挂起
func (handlers *MsgHandlers) FireMsg(name string, content interface{}) {
	handlers.msgFireInfo.Put(&msgFireInfo{name, content})
}

// FireRPC 触发消息, 保证不被挂起
func (handlers *MsgHandlers) FireRPC(method string, data []byte) {
	handlers.rpcFireInfo.Put(&rpcFireInfo{method, data})
}

// DoMsg 将缓冲的消息一次性处理
func (handlers *MsgHandlers) DoMsg() {
	if !handlers.enable {
		return
	}

	var procName string
	defer func() {
		if err := recover(); err != nil {
			log.Error(err, procName)
			if viper.GetString("Config.Recover") == "0" {
				panic(fmt.Sprintln(err, procName))
			}
		}
	}()

	for {
		info, err := handlers.msgFireInfo.Pop()
		if err != nil {
			break
		}

		msg := info.(*msgFireInfo)
		e := handlers.DoNormalMsg(msg.name, msg.content)
		if e != nil {
			log.Error(e)
		}
	}

	for {
		info, err := handlers.rpcFireInfo.Pop()
		if err != nil {
			break
		}

		msg := info.(*rpcFireInfo)
		handlers.DoRPCMsg(msg.methodName, msg.data)
	}
}

// DoRPCMsg RPC执行
func (handlers *MsgHandlers) DoRPCMsg(methodName string, data []byte) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Do RPCMsg error !!! ", err, methodName)
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ifuncs, ok := handlers.rpcFuncs.Load(methodName)
	if !ok {
		log.Error("Method ", methodName, " Can't Find")
		return fmt.Errorf("Method %s can't find", methodName)
	}

	funcs, ok := ifuncs.([]reflect.Value)

	args := serializer.UnSerialize(data)
	callArgs := []reflect.Value{}
	for _, arg := range args {
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}

	for _, rpcFunc := range funcs {
		rpcFunc.Call(callArgs)
	}

	return nil
}

func (handlers *MsgHandlers) addMsgHandler(msgName string, msgHandler reflect.Value) {
	ifuncs, ok := handlers.msgFuncs.Load(msgName)

	var funcs []reflect.Value
	if !ok {
		funcs = make([]reflect.Value, 0, 10)
	} else {
		funcs = ifuncs.([]reflect.Value)
	}

	funcs = append(funcs, msgHandler)
	handlers.msgFuncs.Store(msgName, funcs)
}

func (handlers *MsgHandlers) addRPCHandler(msgName string, msgHandler reflect.Value) {
	ifuncs, ok := handlers.rpcFuncs.Load(msgName)

	var funcs []reflect.Value
	if !ok {
		funcs = make([]reflect.Value, 0, 10)
	} else {
		funcs = ifuncs.([]reflect.Value)
	}

	funcs = append(funcs, msgHandler)
	handlers.rpcFuncs.Store(msgName, funcs)
}

func (handlers *MsgHandlers) getMsgHandler(methodName string, v reflect.Value) (string, reflect.Value, error) {
	methodHead := "MsgProc_"
	methodHeadLen := len(methodHead)

	if len(methodName) < methodHeadLen+1 {
		return "", reflect.ValueOf(nil), fmt.Errorf("")
	}

	if methodName[0:methodHeadLen] != methodHead {
		return "", reflect.ValueOf(nil), fmt.Errorf("")
	}

	msgName := methodName[methodHeadLen:]

	//此处应该检查该函数是否是MsgHanderFunc类型的参数
	return msgName, v, nil
}

func (handlers *MsgHandlers) getRPCHandler(methodName string, v reflect.Value) (string, reflect.Value, error) {
	if strings.Contains(methodName, "RPC_") {
		return methodName[4:], v, nil
	}

	return "", reflect.ValueOf(nil), fmt.Errorf("")
}

// DoNormalMsg 立即触发消息
func (handlers *MsgHandlers) DoNormalMsg(name string, content interface{}) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error("DoNormal error !!! ", err, name)
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ifuncs, ok := handlers.msgFuncs.Load(name)
	if !ok {
		if len(handlers.defaultFuncs) == 0 {
			return fmt.Errorf("Cant find handler %s", name)
		}

		for _, f := range handlers.defaultFuncs {
			if f.IsValid() {
				f.Call([]reflect.Value{reflect.ValueOf(name), reflect.ValueOf(content)})
			}
		}

		return nil
	}

	for _, msgFunc := range ifuncs.([]reflect.Value) {
		msgFunc.Call([]reflect.Value{reflect.ValueOf(content)})
	}

	return nil
}
