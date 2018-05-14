package main

import (
	"protoMsg"
	"zeus/linmath"
	"zeus/msgdef"
)

// EntityMsgProc Entity消息处理函数
type EntityMsgProc struct {
	e *Entity
}

func (proc *EntityMsgProc) MsgProc_PropsSyncClient(content msgdef.IMsg) {
	msg := content.(*msgdef.PropsSyncClient)
	proc.e.CastMsgToAllClient(msg)
}

func (proc *EntityMsgProc) MsgProc_SyncUserState(content msgdef.IMsg) {
	proc.e.Debug("MsgProc_SyncUserState")
	//msg := content.(*msgdef.SyncUserState)
	//proc.e.syncClientUserState(msg)
}

func (proc *EntityMsgProc) MsgProc_SessClosed(content interface{}) {
	// log.Info("SessClosed ", proc.e)
	proc.e.SetClient(nil)

	// proc.e.LeaveCell()
}

//MsgProc_MoveReq 客户端更新坐标
func (proc *EntityMsgProc) MsgProc_MoveReq(content msgdef.IMsg) {
	msgRecv := content.(*protoMsg.MoveReq)
	//proc.e.Debug("MoveReq:", msgRecv, ", isGhost: ", proc.e.IsGhost())

	pos := linmath.Vector3{
		msgRecv.GetPos().GetX(),
		msgRecv.GetPos().GetY(),
		msgRecv.GetPos().GetZ()}

	rota := linmath.Vector3{
		msgRecv.GetRota().GetX(),
		msgRecv.GetRota().GetY(),
		msgRecv.GetRota().GetZ()}

	proc.e.SetPos(pos)
	proc.e.SetRota(rota)

	//只有real entity才会转发消息
	if !proc.e.IsGhost() {
		sendMsg := &protoMsg.MoveUpdate{
			EntityID: proc.e.GetID(),
			Pos:      msgRecv.GetPos(),
			Rota:     msgRecv.GetRota(),
			Stoped:   msgRecv.Stoped,
		}

		proc.e.CastMsgToAllClientExceptMe(sendMsg)

		proc.e.SyncToGhosts(content)
	}

}

// 通知real切换
func (proc *EntityMsgProc) MsgProc_NewRealNotify(content msgdef.IMsg) {
	msg := content.(*msgdef.NewRealNotify)

	proc.e.Debug("NewRealNotify")

	if !proc.e.IsGhost() {
		proc.e.Error("NewRealNotify, not ghost")
		return
	}

	proc.e.SetRealServerID(msg.RealServerID)
	proc.e.SetRealCellID(msg.RealCellID)
}
