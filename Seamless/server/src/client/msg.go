package main

import (
	"common"
	"protoMsg"
	"sync/atomic"
	"time"

	"github.com/cihub/seelog"
)

func (cli *Client) TryEnterCell() {
	lastEnterCellTime := atomic.LoadInt64(&cli.lastEnterCellTime)
	now := time.Now().UnixNano() / 1e6
	if now-lastEnterCellTime < 2000 {
		return
	}
	cli.RPCCall(common.ServerTypeLobby, "EnterCellReq")
	atomic.StoreInt64(&cli.lastEnterCellTime, now)
}

func (cli *Client) SendStopMoveMsg() {
	us := &cli.User.UserState
	cli.SendMsgToCell(&protoMsg.MoveReq{
		Pos: &protoMsg.Vector3{
			X: us.pos.X / MapRate,
			Z: us.pos.Z / MapRate,
		},
		Rota: &protoMsg.Vector3{
			X: us.rota.X,
			Z: us.rota.Z,
		},
		Stoped: true,
	})
}

func (cli *Client) NormalAttack(targetID uint64) {
	msg := &protoMsg.AttackReq{
		SkillID:  1001,
		TargetID: targetID,
	}
	cli.RPCCall(common.ServerTypeCellApp, "DoSkill", msg)
	seelog.Debug("DoSkill:", msg)
}

func (cli *Client) DoSkill1() {
	p := GetMainView().GetMousePos()
	rect := cli.User.UserView.us.GetViewMapRect()
	msg := &protoMsg.AttackReq{
		SkillID:  1002,
		TargetID: 0,
		X:        float32(p.X+rect.X) / MapRate,
		Z:        float32(p.Y+rect.Y) / MapRate,
	}
	cli.RPCCall(common.ServerTypeCellApp, "DoSkill", msg)
	seelog.Debug("DoSkill:", msg)
}
