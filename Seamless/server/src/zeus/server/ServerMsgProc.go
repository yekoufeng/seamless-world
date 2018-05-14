package server

import (
	"zeus/iserver"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
)

type serverMsgProc struct {
	srv *Server
}

func (p *serverMsgProc) getEntities(cellID uint64) iserver.IEntities {

	if cellID == 0 {
		return p.srv
	}

	//e := p.srv.GetEntity(cellID)
	e := p.srv.GetEntities(cellID)
	if e == nil {
		return nil
	}

	return e.(iserver.IEntities)
}

// 把消息投递到一个实体里
func (p *serverMsgProc) MsgProc_EntityMsgTransport(content msgdef.IMsg) {
	msg := content.(*msgdef.EntityMsgTransport)
	if es := p.getEntities(msg.CellID); es != nil {
		es.FireMsg(msg.Name(), msg)
	} else {
		log.Error("EntityMsgTransport GetEntities failed, CellID:", msg.CellID)
	}
}

// 把消息投递到一个空间里
func (p *serverMsgProc) MsgProc_SrvMsgTransport(content msgdef.IMsg) {
	msg := content.(*msgdef.SrvMsgTransport)

	if es := p.getEntities(msg.CellID); es != nil {
		es.FireMsg(msg.Name(), msg)
	} else {
		log.Error("SrvMsgTransport GetEntities failed, CellID:", msg.CellID)
	}
}
