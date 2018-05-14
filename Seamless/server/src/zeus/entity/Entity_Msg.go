package entity

import (
	"errors"
	"fmt"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/pool"
	"zeus/serializer"
	"zeus/sess"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// PostGateRaw 投递一个gateway原始消息
func (e *Entity) PostGateRaw(srvType uint8, rawMsg []byte) error {
	if srvType == iserver.ServerTypeClient || srvType == iserver.ServerTypeGateway {
		return nil
	}

	if e.state != iserver.Entity_State_Loop {
		return fmt.Errorf("error entity state %d ", e.state)
	}

	if rawMsg == nil {
		return fmt.Errorf("Messsage is nil")
	}

	srvID, cellID, err := e.getEntitySrvID(srvType)
	if err != nil {
		return err
	}

	if srvID == iserver.GetSrvInst().GetSrvID() {
		return nil
	}

	var packMsg msgdef.IMsg
	if packMsg, err = e.packGateRawMsg(srvType, cellID, rawMsg); err != nil {
		return err
	}

	return iserver.GetSrvInst().PostMsgToSrv(srvID, packMsg)

}

// PostGateRawToCell 投递一个原始消息给它的SpaceEntity部分
func (e *Entity) PostGateRawToCell(rawMsg []byte) error {
	if e.cellSrvType == 0 {
		return errors.New("PostGateRawToCell not exist space part")
	}

	return e.PostGateRaw(e.cellSrvType, rawMsg)
}

// Post 投递一个消息给指定的部分
// 发送消息可能失败，在极端情况下
func (e *Entity) Post(srvType uint8, msg msgdef.IMsg) error {
	if msg == nil {
		return fmt.Errorf("Message is nil")
	}

	//如果投递给客户端且本地就有客户端连接话，直接投递
	if srvType == iserver.ServerTypeClient && e.GetClientSess() != nil {
		e.GetClientSess().Send(msg)
		return nil
	}

	srvID, cellID, err := e.getEntitySrvID(srvType)
	if err != nil {
		return err
	}

	var packMsg msgdef.IMsg
	if packMsg, err = e.packMsg(srvType, cellID, msg); err != nil {
		return err
	}

	//如果是投递给自己的消息，直接处理
	if srvID == iserver.GetSrvInst().GetSrvID() {
		e.IMsgHandlers.FireMsg(packMsg.Name(), packMsg)
		return nil
	}

	return iserver.GetSrvInst().PostMsgToSrv(srvID, packMsg)
}

// DelayPost 延迟发送消息, 在帧尾发送
func (e *Entity) DelayPost(srvType uint8, msg msgdef.IMsg) error {
	if msg == nil {
		return fmt.Errorf("Message is nil")
	}

	e.delayedMsgs = append(e.delayedMsgs, &delayedSendMsg{
		srvType: srvType,
		msg:     msg,
	})

	return nil
}

// 发送所有延迟发送的消息
func (e *Entity) FlushDelayedMsgs() {
	if len(e.delayedMsgs) == 0 {
		return
	}

	for _, dm := range e.delayedMsgs {
		e.Post(dm.srvType, dm.msg)
	}

	e.delayedMsgs = e.delayedMsgs[0:0]
}

// PostToCell 投递一个消息给它的SpaceEntity部分
func (e *Entity) PostToCell(msg msgdef.IMsg) error {
	if e.cellSrvType == 0 {
		return fmt.Errorf("PostToCell not exist space part")
	}

	return e.Post(e.cellSrvType, msg)
}

// RPC 快速消息调用
func (e *Entity) RPC(srvType uint8, methodName string, args ...interface{}) error {
	data := serializer.Serialize(args...)
	msg := &msgdef.RPCMsg{}
	msg.ServerType = srvType
	msg.SrcEntityID = e.entityID
	msg.MethodName = methodName
	msg.Data = data

	return e.DelayPost(srvType, msg)
}

// RPCOther 快速触发一个其它实体的方法
/*func (e *Entity) RPCOther(srvType uint8, srcEntityID uint64, methodName string, args ...interface{}) error {
	data := serializer.Serialize(args...)
	msg := &msgdef.RPCMsg{}
	msg.ServerType = srvType
	msg.SrcEntityID = srcEntityID
	msg.MethodName = methodName
	msg.Data = data

	return e.DelayPost(srvType, msg)
}*/

// SpaceRPC 快速调用spaceentity的rpc方法
/*func (e *Entity) SpaceRPC(methodName string, args ...interface{}) error {
	if e.cellSrvType == 0 {
		return fmt.Errorf("not exist space part")
	}
	return e.RPC(e.cellSrvType, methodName, args)
}*/

func (e *Entity) ExportPackMsg(srvType uint8, cellID uint64, msg msgdef.IMsg) (msgdef.IMsg, error) {
	return e.packMsg(srvType, cellID, msg)
}
func (e *Entity) packMsg(srvType uint8, cellID uint64, msg msgdef.IMsg) (msgdef.IMsg, error) {
	//buf := make([]byte, sess.MaxMsgBuffer)
	buf := pool.Get(sess.MaxMsgBuffer)

	var encBuf []byte
	var err error

	if viper.GetBool("Config.EncryptEnabled") && srvType == iserver.ServerTypeClient {
		encBuf, err = sess.EncodeMsgWithEncrypt(msg, buf, true, true)
	} else {
		encBuf, err = sess.EncodeMsg(msg, buf, true)
	}
	if err != nil {
		return nil, err
	}

	msgContent := make([]byte, len(encBuf))
	copy(msgContent, encBuf)
	pool.Put(buf)

	return &msgdef.EntityMsgTransport{
		SrvType:    srvType,
		CellID:     cellID,
		EntityID:   e.entityID,
		MsgContent: msgContent,
	}, nil
}

func (e *Entity) packGateRawMsg(srvType uint8, cellID uint64, rawMsg []byte) (msgdef.IMsg, error) {

	return &msgdef.EntityMsgTransport{
		SrvType:    srvType,
		EntityID:   e.entityID,
		CellID:     cellID,
		IsGateway:  true,
		MsgContent: rawMsg,
	}, nil
}

/////////////////////////////////////////////////////////////////////////////
// 此处为Entity相关消息
////////////////////////////////////////////////////////////////////////////

// MsgProc_PropsSync 实体间同步数据的消息
func (e *Entity) MsgProc_PropsSync(content msgdef.IMsg) {
	msg := content.(*msgdef.PropsSync)

	if msg.Num > 100 {
		log.Error("prop num is exceed 100 , is right ?")
		return
	}

	e.ReflushFromMsg(int(msg.Num), msg.Data)
}

// MsgProc_RPCMsg 实体间RPC调用消息处理
func (e *Entity) MsgProc_RPCMsg(content msgdef.IMsg) {
	msg := content.(*msgdef.RPCMsg)
	// 从客户端收到的消息要判断要调用的服务器类型

	if msg.ServerType == iserver.GetSrvInst().GetSrvType() {
		if msg.SrcEntityID == e.GetID() || msg.SrcEntityID == 0 {
			e.FireRPC(msg.MethodName, msg.Data)
		}
	} else {
		e.Post(msg.ServerType, msg)
	}
}
