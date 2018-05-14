package main

import (
	//"entitydef"
	"fmt"
	//"protoMsg"
	"time"
	"zeus/linmath"
	"zeus/space"
	//	log "github.com/cihub/seelog"
)

type SpaceUser struct {
	entitydef.PlayerDef
	space.Entity
}

func (user *SpaceUser) Init() {
	fmt.Println("Space user inited", user.GetDBID())
	user.RegMsgProc(&SpaceUserMsgProc{user: user})
	user.SetWatcher(true)

	user.SetPos(linmath.Vector3_Zero())

	// log.Error(user.GetName())
	// log.Error(user.GetHP())
	// log.Error(user.GetBackPackProp())
	// log.Error(user.GetGunSight())

	user.GetEntities().AddDelayCall(func() {
		user.RPC(3, user.GetID(), "TestRPCNoArgs")
	}, time.Second)

	// user.GetEntities().RegTimer(func() {
	// 	user.SetName("SpaceUser")

	// 	i := rand.Intn(100)

	// 	user.SetHP(uint32(i))

	// 	user.SetGunSight(uint64(i))

	// 	user.SetBackPackProp(&protoMsg.BackPackProp{Baseid: uint32(i)})
	// }, time.Second)
}

func (user *SpaceUser) Destroy() {
	fmt.Println("Space user destroy", user.GetDBID())
}

func (user *SpaceUser) OnMarkerEnterAOI(m space.IMarker) {
}

func (user *SpaceUser) OnMarkerLeaveAOI(m space.IMarker) {
}

type SpaceUserMsgProc struct {
	user *SpaceUser
}

func (proc *SpaceUserMsgProc) RPC_TestRPC(i16 int16, u32 uint32, flag bool, str string, prop *protoMsg.BackPackProp) {
	fmt.Println(i16, u32, flag, str, prop)
}
