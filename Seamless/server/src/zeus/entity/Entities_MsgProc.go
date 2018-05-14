package entity

import (
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/serializer"
	"zeus/sess"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

//MsgProc_EntityMsgTransport Entity之间转发消息
func (es *Entities) MsgProc_EntityMsgTransport(content msgdef.IMsg) {
	msg := content.(*msgdef.EntityMsgTransport)
	ie := es.GetEntity(msg.EntityID)
	if ie == nil {
		return
	}

	//如果是Gateway转过来的消息，但当前实体为ghost，需要转给real
	if msg.IsGateway && ie.IsGhost() {
		if se, ok := ie.(iserver.ISendMsgToReal); ok {
			innerMsg, err := sess.DecodeMsg(msg.MsgContent[3], msg.MsgContent[4:])
			if err != nil {
				log.Error("Decode innerMsg failed", err, msg)
				return
			}

			se.SendMsgToReal(innerMsg)
			return
		}
	}

	te := ie.(iEntityCtrl)
	if msg.SrvType == iserver.ServerTypeClient {
		if sess := te.GetClientSess(); sess != nil {
			sess.SendRaw(msg.MsgContent)
		}
	} else {
		innerMsg, err := sess.DecodeMsg(msg.MsgContent[3], msg.MsgContent[4:])
		if err != nil {
			log.Error("Decode innerMsg failed", err, msg)
			return
		}

		if !iserver.GetSrvInst().IsSrvValid() {
			iserver.GetSrvInst().HandlerSrvInvalid(ie.GetID())
		} else {
			te.FireMsg(innerMsg.Name(), innerMsg)
		}
	}
}

//MsgProc_CreateEntityReq 请求创建Entity
func (es *Entities) MsgProc_CreateEntityReq(content msgdef.IMsg) {
	maxLoad := viper.GetInt("Config.MaxLoad")
	if maxLoad != 0 {
		load := iserver.GetSrvInst().GetLoad()
		if load > maxLoad {
			log.Warnf("Overload Current:%d MaxLoad:%d Msg:%s", load, maxLoad, content)
			return
		}
	}

	msg := content.(*msgdef.CreateEntityReq)
	params := serializer.UnSerialize(msg.InitParam)
	if len(params) != 1 {
		log.Error("Init param error ", params)
		return
	}
	_, err := es.CreateEntity(msg.EntityType, msg.EntityID, msg.DBID, 0, params[0], false, 0)
	if err != nil {
		log.Error("Create entity failed ", err, msg)
	}
}

//MsgProc_DestroyEntityReq 请求删除Entity
func (es *Entities) MsgProc_DestroyEntityReq(content msgdef.IMsg) {
	msg := content.(*msgdef.DestroyEntityReq)
	log.Info("MsgProc_DestroyEntityReq entityID ", msg.EntityID)

	err := es.DestroyEntity(msg.EntityID)
	if err != nil {
		log.Error("Delete entity failed ", err, msg)
	}

}

func (es *Entities) MsgProc_SrvMsgTransport(content msgdef.IMsg) {
	msg := content.(*msgdef.SrvMsgTransport)
	innerMsg, err := sess.DecodeMsg(msg.MsgContent[3], msg.MsgContent[4:])

	if err != nil {
		log.Error(err, "msg", sess.GetMsgID(msg.MsgContent[4:]))
		return
	}

	es.FireMsg(innerMsg.Name(), innerMsg)
}
