package main

import (
	"fmt"
	//"protoMsg"
	"zeus/entity"
)

type TestSpace struct {
	entity.Entity
}

func (sp *TestSpace) Init() {
	fmt.Println("Test space inited", sp.GetID())
	sp.RegMsgProc(&TestSpaceMsgProc{})
}

func (sp *TestSpace) Destroy() {
	fmt.Println("Test space destroy", sp.GetID())
}

type TestSpaceMsgProc struct {
}

func (proc *TestSpaceMsgProc) RPC_TestRPC(flag bool, index uint16, str string) {
	fmt.Println(flag, index, str)
}

func (proc *TestSpaceMsgProc) RPC_TestRPCProto(msg *protoMsg.Vector3) {
	fmt.Println(msg)
}
