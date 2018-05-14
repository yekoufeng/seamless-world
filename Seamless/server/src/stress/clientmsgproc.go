package main

import (
	"common"
	"fmt"
	"math/rand"
	"protoMsg"
	"time"
	"zeus/msgdef"
)

type ClientMsgProc struct {
	c *Client
}

func (proc *ClientMsgProc) MsgProc_ClientVertifySucceedRet(content interface{}) {
	verifyRet := &TranRet{}
	verifyRet.Action = "UserVerity"
	verifyRet.Err = nil
	verifyRet.Cost = time.Now().Sub(proc.c.verifyTime).Nanoseconds() / 1000000
	proc.c.retChan <- verifyRet

	proc.c.enterGatewaySig <- true
}

func (proc *ClientMsgProc) MsgProc_RPCMsg(content msgdef.IMsg) {
	msg := content.(*msgdef.RPCMsg)
	// 从客户端收到的消息要判断要调用的服务器类型
	proc.c.msgClient.FireRPC(msg.MethodName, msg.Data)
}

func (proc *ClientMsgProc) MsgProc_EnterSpace(content msgdef.IMsg) {
	msg := content.(*msgdef.EnterCell)
	proc.c.enterSpaceMsg = msg
	proc.c.enterSpaceSig <- true

	fmt.Println(msg)
}

func (proc *ClientMsgProc) RPC_MatchSuccess(mapid uint32, skybox uint32) {
	matchRet := &TranRet{}
	matchRet.Action = "MatchSuccess"
	matchRet.Err = nil
	matchRet.Cost = time.Now().Sub(proc.c.queueTime).Nanoseconds() / 1000000
	fmt.Println(matchRet)
	proc.c.retChan <- matchRet

	proc.c.matchSuccSig <- true
}

func (proc *ClientMsgProc) RPC_SyncApplyList(msg *protoMsg.SyncFriendApplyList) {

}

func (proc *ClientMsgProc) RPC_SyncTeamInfoRet(msg *protoMsg.SyncTeamInfoRet) {

}

func (proc *ClientMsgProc) RPC_NotifyWaitingNums(uint32) {}

func (proc *ClientMsgProc) RPC_SyncFriendList(*protoMsg.SyncFriendList) {}

func (proc *ClientMsgProc) RPC_InitOwnGoodsInfo(*protoMsg.OwnGoodsItem) {}

func (proc *ClientMsgProc) RPC_OnlineCheckMatchOpen(bool, bool, bool) {}

func (proc *ClientMsgProc) RPC_ZoneNotify(*protoMsg.ZoneNotify) {

}

func (proc *ClientMsgProc) RPC_TestRPC(i16 int16, u32 uint32, flag bool, str string, prop *protoMsg.BackPackProp) {
	fmt.Println(i16, u32, flag, str, prop)
}

func (proc *ClientMsgProc) RPC_TransResult(action string, success bool, cost int64) {
	// fmt.Println("TransResult: ", action, success, cost)
	ret := &TranRet{}
	ret.Action = action
	ret.Cost = cost
	if !success {
		ret.Err = fmt.Errorf("Failed")
	}
	proc.c.retChan <- ret
}

func (proc *ClientMsgProc) RPC_ExpectTime(t uint64) {
	// fmt.Println("RPC_ExpectTime", t)
}

func (proc *ClientMsgProc) RPC_JumpAir(set uint64) {

}

func (proc *ClientMsgProc) RPC_PickupItem(id uint64) {

}

func (proc *ClientMsgProc) RPC_AddNewMail() {

}

func (proc *ClientMsgProc) MsgProc_EntityAOIS(content interface{}) {
	// fmt.Println("EntityAOIS ", content)
}

// func (proc *ClientMsgProc) MsgProc_MRolePropsSyncClient(content msgdef.IMsg) {
// 	msg := content.(*msgdef.MRolePropsSyncClient)
// 	fmt.Println("MsgProc_PropsSyncClient", msg)
// }

func (proc *ClientMsgProc) RPC_SyncUserTeamid(teamID uint64) {
	fmt.Println("SyncUserTeamid", teamID)
}

func (proc *ClientMsgProc) RPC_SetMatchNumber(num uint64) {
}

func (proc *ClientMsgProc) RPC_AirLine(x, y, z, u float64) {

}

func (proc *ClientMsgProc) RPC_ParachutePos() {
	proc.c.RPCCall(common.ServerTypeRoom, 0, "ReleaseParachute")

	go func() {
		time.Sleep(15 * time.Second)
		proc.c.setBaseState(RoomPlayerBaseState_Stand)
		proc.c.inGame = true
	}()
}

func (proc *ClientMsgProc) RPC_RefreshPackCellNotify(*protoMsg.RefreshPackCellNotify) {

}

func (proc *ClientMsgProc) RPC_ShrinkRefresh(id uint64) {

}

func (proc *ClientMsgProc) RPC_BombDam(x, y, z, r float32) {

}

func (proc *ClientMsgProc) RPC_UpdateTotalNum(uint32) {

}

func (proc *ClientMsgProc) RPC_UpdateAirLeft(uint32) {

}

func (proc *ClientMsgProc) RPC_BombRefresh(x, y, z, r float32) {

}

func (proc *ClientMsgProc) RPC_MapCharacterResultNotify(*protoMsg.MapCharacterResultNotify) {}

func (proc *ClientMsgProc) RPC_OnDamage(typ uint64) {

}

func (proc *ClientMsgProc) RPC_BombDisapear() {

}

func (proc *ClientMsgProc) RPC_UpdateKillNum(num uint32) {
}

func (proc *ClientMsgProc) RPC_UpdateAliveNum(num uint32) {
}

func (proc *ClientMsgProc) RPC_AllowParachute() {
	proc.c.setBaseState(RoomPlayerBaseState_Inplane)

	go func() {
		delay := rand.Intn(10) + 3
		time.Sleep(time.Duration(delay) * time.Second)
		proc.c.setBaseState(RoomPlayerBaseState_Glide)
	}()
}

func (proc *ClientMsgProc) RPC_InitNotifyMysqlDbAddr(addr string) {

}

func (proc *ClientMsgProc) MsgProc_AdjustUserState(content msgdef.IMsg) {
	msg := content.(*msgdef.AdjustUserState)

	proc.c.curState.Combine(msg.Data)

	// fmt.Println("AdjustUserState", msg.Data, proc.c.curState.String())
}

// func (proc *ClientMsgProc) MsgProc_AOIPosChange(content interface{}) {
// 	msg := content.(*msgdef.AOIPosChange)
// 	log.Error("w11111133333", msg)
// }

// func (proc *ClientMsgProc) MsgProc_EnterAOI(content interface{}) {
// 	msg := content.(*msgdef.EnterAOI)
// 	log.Error("w111111", msg)
// }

type ClientRoomMsgProc struct {
	c *Client
}

func (proc *ClientRoomMsgProc) MsgProc_ClientVertifySucceedRet(content interface{}) {
	enterRoomRet := &TranRet{}
	enterRoomRet.Action = "EnterRoom"
	enterRoomRet.Err = nil
	enterRoomRet.Cost = time.Now().Sub(proc.c.queueTime).Nanoseconds() / 1000000

	fmt.Println(enterRoomRet)

	proc.c.retChan <- enterRoomRet

	proc.c.RPCCall(common.ServerTypeRoom, 0, "ParachuteReady")
	proc.c.enterRoomSig <- true
}
