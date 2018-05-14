package main

import (
	"zeus/msgdef"

	log "github.com/cihub/seelog"

	"protoMsg"
)

// LobbySrvMsgProc Lobby服务器消息处理类
type LobbySrvMsgProc struct {
	srv *LobbySrv
}

//MsgProc_CellInfoRet 收到cell信息，开始进入地图
func (p *LobbySrvMsgProc) MsgProc_CellInfoRet(content msgdef.IMsg) {
	msg, ok := content.(*protoMsg.CellInfoRet)
	if !ok {
		log.Error("MsgProc_CellInfoRet, 消息解析错误")
		return
	}

	entity := p.srv.GetEntity(msg.GetEntityID())
	if entity == nil {
		return
	}

	if user, ok := entity.(*LobbyUser); ok {

		//开始进入场景
		user.EnterCell(msg)
	}

	log.Debug("MsgProc_CellInfoRet spaceID: ", msg.GetSpaceID(), ", cellID: ", msg.GetCellID(), ", entityID: ", msg.GetEntityID())
}
