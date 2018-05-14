package main

import (
	"zeus/msgdef"
	"zeus/serializer"

	log "github.com/cihub/seelog"
)

// CellMsgProc RoomUser的消息处理函数
type CellMsgProc struct {
	cell *Cell
}

// 创建ghost
func (proc *CellMsgProc) MsgProc_CreateGhostReq(content msgdef.IMsg) {
	msg := content.(*msgdef.CreateGhostReq)

	log.Info("CreateGhostReq, entityID: ", msg.EntityID, ", dbid: ", msg.DBID, ", CellID: ", proc.cell.GetID())
	params := serializer.UnSerialize(msg.InitParam)
	if len(params) < 1 {
		log.Error("CreateGhostReq: Unmarshal initparam error ", msg.InitParam)
		return
	}

	entity, err := proc.cell.AddEntity(msg.EntityType, msg.EntityID, msg.DBID, params[0], true, msg.RealServerID, msg.RealCellID)
	if err != nil {
		log.Error("CreateGhostReq: Add entity error ", err, msg)
		return
	}

	entity.ReflushFromMsg(int(msg.PropNum), msg.Props)

	log.Info("AddEntity 玩家初始位置, X: ", msg.Pos.X, ", Y:", msg.Pos.Y, ", Z: ", msg.Pos.Z)
	entity.SetPos(msg.Pos)
	entity.updatePosCoord(entity.GetPos())
}

// 删除ghost
func (proc *CellMsgProc) MsgProc_DeleteGhostReq(content msgdef.IMsg) {
	msg := content.(*msgdef.DeleteGhostReq)

	log.Info("DeleteGhostReq, entityID: ", msg.EntityID)

	proc.cell.DestroyEntity(msg.EntityID)
}

// real和ghost切换
func (proc *CellMsgProc) MsgProc_TransferRealReq(content msgdef.IMsg) {
	msg := content.(*msgdef.TransferRealReq)

	log.Info("TransferRealReq, entityID: ", msg.EntityID)

	proc.cell.GhostToReal(msg)

}

// MsgProc_EnterCellReq 进入cell
func (proc *CellMsgProc) MsgProc_EnterCellReq(content msgdef.IMsg) {
	msg := content.(*msgdef.EnterCellReq)

	log.Debug("EnterCellReq, entityID: ", msg.EntityID, ", dbid: ", msg.DBID)
	params := serializer.UnSerialize(msg.InitParam)
	if len(params) < 1 {
		log.Error("EnterCellReq: Unmarshal initparam error ", msg.InitParam)
		return
	}

	if proc.cell.GetEntity(msg.EntityID) != nil {
		log.Debug("duplicate create celluser")
		return
	}

	entity, err := proc.cell.AddEntity(msg.EntityType, msg.EntityID, msg.DBID, params[0], true, 0, 0)
	if err != nil {
		log.Error("EnterCellReq: Add entity error ", err, msg)
		return
	}

	log.Info("AddEntity 玩家初始位置, X: ", msg.Pos.X, ", Y:", msg.Pos.Y, ", Z: ", msg.Pos.Z)
	entity.SetPos(msg.Pos)
	entity.updatePosCoord(entity.GetPos())
}
