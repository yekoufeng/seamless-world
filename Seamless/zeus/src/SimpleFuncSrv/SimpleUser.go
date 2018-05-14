package main

import (
	//"entitydef"
	"fmt"
	//"protoMsg"
	"time"
	"zeus/entity"

	log "github.com/cihub/seelog"
)

type SimpleUser struct {
	entitydef.PlayerDef
	entity.Entity
}

func (user *SimpleUser) Init() {
	fmt.Println("Simple user init", user.GetDBID())
	user.RegMsgProc(&SimpleUserMsgProc{user: user})
	// user.SetName("SimpleUser")

	// user.SetHP(50)

	// user.SetBackPackProp(&protoMsg.BackPackProp{Baseid: 1024})

	GetSimpleSrvInst().AddDelayCall(func() { user.EnterSpace(50000) }, time.Second)

	GetSimpleSrvInst().RegTimerByObj(user, func() {
		// if err := user.RPC(iserver.ServerTypeClient, user.GetID(), "TestRPC", int16(-13), uint32(1111), true, "hello", &protoMsg.BackPackProp{Baseid: 1024}); err != nil {
		// 	log.Error(err)
		// }

		if err := user.RPC(4, user.GetID(), "TestRPC", int16(-13), uint32(1111), true, "hello", &protoMsg.BackPackProp{Baseid: 1024}); err != nil {
			log.Error(err)
		}

	}, 3*time.Second)
}

func (user *SimpleUser) Destroy() {
	GetSimpleSrvInst().UnregTimerByObj(user)
	fmt.Println("Simple user destroy", user.GetDBID())
}

type SimpleUserMsgProc struct {
	user *SimpleUser
}

func (proc *SimpleUserMsgProc) RPC_TestRPC(i16 int16, u32 uint32, flag bool, str string, prop *protoMsg.BackPackProp) {
	fmt.Println(i16, u32, flag, str, prop)

	// if err := proc.user.RPC(iserver.ServerTypeClient, proc.user.GetID(), "TestRPC", i16, u32, flag, str, prop); err != nil {
	// 	log.Error(err)
	// }
}

func (proc *SimpleUserMsgProc) RPC_TestRPCNoArgs() {
	fmt.Println("RPC_TestRPCNoArgs")
}
