package main

import (
	"fmt"
	"time"
	"zeus/linmath"
	"zeus/space"
)

type SpaceNPC struct {
	space.Entity
}

// NPCDef 玩家实体定义

func (user *SpaceNPC) Init() {
	fmt.Println("SpaceNPC inited", user.GetDBID())

	user.SetMarker(true)
	user.SetWatcher(true)

	step := 0
	user.GetEntities().RegTimerByObj(user, func() {
		xPos := float32(step)
		step += int(user.GetID() % 4)
		if xPos >= 100 {
			xPos = 0
			step = 0
		}
		//zPos := r.Float32() * 100

		//fmt.Println("Space user", user.GetDBID(), "moveto", xPos, zPos)

		user.SetPos(linmath.Vector3{
			X: xPos,
			//Z: zPos,
		})
	}, 3*time.Second)

	// GetSrvInst().RegTimerByObj(user, func() {
	// 	user.CastRPCToAllClient("doTest", true, "hello")
	// }, time.Duration(user.GetID()%4+1)*time.Second)
}

func (user *SpaceNPC) Destroy() {
	fmt.Println("SpaceNPC destroy", user.GetDBID())
	user.GetEntities().UnregTimerByObj(user)
	//user.LeaveSpace()
}

func (user *SpaceNPC) OnMarkerEnterAOI(m space.IMarker) {
	// distance := user.GetPos().Sub(m.GetPos())
	// fmt.Println("Space user", user.GetDBID(),
	// 	"OnMarkerEnterAOI, MyPOS:", user.GetPos().X, user.GetPos().Z,
	// 	"Marker POS:", m.GetPos().X, m.GetPos().Z,
	// 	"Distance:", distance.X, distance.Z)
}

func (user *SpaceNPC) OnMarkerLeaveAOI(m space.IMarker) {
	// distance := user.GetPos().Sub(m.GetPos())
	// fmt.Println("Space user", user.GetDBID(),
	// 	"OnMarkerLeaveAOI, MyPOS:", user.GetPos().X, user.GetPos().Z,
	// 	"Marker POS:", m.GetPos().X, m.GetPos().Z,
	// 	"Distance:", distance.X, distance.Z)
}
