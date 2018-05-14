package main

import (
	"fmt"
	//"protoMsg"
	"time"
	"zeus/space"

	log "github.com/cihub/seelog"
)

type TestSpace struct {
	space.Space
}

func (sp *TestSpace) Init() {
	fmt.Println("Test space inited", sp.GetID())

	sp.AddDelayCall(func() {
		if err := sp.RPC(3, sp.GetID(), "TestRPC", true, uint16(1024), "hello"); err != nil {
			log.Error(err)
		}
		msg := &protoMsg.Vector3{}
		msg.X = 3.1415
		msg.Y = 9786
		msg.Z = -3.1415926
		if err := sp.RPC(3, sp.GetID(), "TestRPCProto", msg); err != nil {
			log.Error(err)
		}
	}, 1*time.Second)
}

func (sp *TestSpace) Destroy() {
	fmt.Println("Test space destroy", sp.GetID())
}
