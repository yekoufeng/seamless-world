package entity

import (
	"fmt"
	"zeus/global"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/serializer"
	"zeus/sess"
)

// EntityProxy 实体的一个代理，可以方便的传递
type EntityProxy struct {
	SrvID    uint64
	CellID   uint64
	EntityID uint64
}

// NewEntityProxy 创建一个新的实体代理
func NewEntityProxy(srvID uint64, cellID uint64, entityID uint64) *EntityProxy {
	return &EntityProxy{
		SrvID:    srvID,
		CellID:   cellID,
		EntityID: entityID,
	}
}

// RegEntityProxy 注册Proxy信息至global变量
func RegEntityProxy(name string, proxy *EntityProxy) {
	global.GetGlobalInst().SetGlobalEntityProxy(name, proxy.EntityID, proxy.SrvID, proxy.CellID)
}

// GetEntityProxy 从global变量获取Proxy
func GetEntityProxy(name string) *EntityProxy {
	entityID, srvID, cellID := global.GetGlobalInst().GetGlobalEntityProxy(name)
	if entityID == 0 {
		return nil
	}

	return NewEntityProxy(srvID, cellID, entityID)
}

func (e *EntityProxy) String() string {
	return fmt.Sprintf("%+v", *e)
}

// Post 将某个消息投递给该实体
func (e *EntityProxy) Post(msg msgdef.IMsg) error {

	if msg == nil {
		return fmt.Errorf("Message is nil")
	}

	var packMsg msgdef.IMsg
	var err error
	if packMsg, err = e.packMsg(msg); err != nil {
		return err
	}

	return iserver.GetSrvInst().PostMsgToSrv(e.SrvID, packMsg)
}

// RPC 快速消息调用，可以快速触发一个其它实体的方法, 服务器间使用
func (e *EntityProxy) RPC(srvType uint8, methodName string, args ...interface{}) error {
	data := serializer.Serialize(args...)
	msg := &msgdef.RPCMsg{}
	msg.SrcEntityID = e.EntityID
	msg.ServerType = srvType
	msg.MethodName = methodName
	msg.Data = data

	return e.Post(msg)
}

func (e *EntityProxy) packMsg(msg msgdef.IMsg) (msgdef.IMsg, error) {

	buf := make([]byte, sess.MaxMsgBuffer)
	encBuf, err := sess.EncodeMsg(msg, buf, true)
	if err != nil {
		return nil, err
	}
	msgContent := make([]byte, len(encBuf))
	copy(msgContent, encBuf)

	return &msgdef.EntityMsgTransport{
		SrvType:    0,
		CellID:     e.CellID,
		EntityID:   e.EntityID,
		MsgContent: msgContent,
	}, nil
}
