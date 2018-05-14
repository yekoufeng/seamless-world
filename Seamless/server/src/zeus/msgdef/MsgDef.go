package msgdef

import (
	"fmt"
	"reflect"

	"runtime/debug"

	log "github.com/cihub/seelog"
)

// MsgInfo 消息信息
type MsgInfo struct {
	id      uint16
	name    string
	msgType reflect.Type
}

// MsgDef 消息定义管理器，包含是消息与消息号的映射结构，客户端初始化的时候，由服务器下发下去
type MsgDef struct {
	id2Info   map[uint16]*MsgInfo
	name2Info map[string]*MsgInfo
}

// Init 初始消息定义管理器
func (def *MsgDef) Init() {
	def.InitBytesMsg()
}

// GetMsgInfo 根据消息号，获得消息的类型，名称，如果是protobuf消息，获得proto消息的容器
func (def *MsgDef) GetMsgInfo(msgID uint16) (msgName string, msgContent IMsg, err error) {

	//此处会被多线程调用，不确定会不会有问题
	info, ok := def.id2Info[msgID]

	if !ok {
		return "", nil, fmt.Errorf("不存在消息号 ID: %d", msgID)
	}

	return info.name, reflect.New(info.msgType.Elem()).Interface().(IMsg), nil
}

// IsMsgExist 消息是否存在
func (def *MsgDef) IsMsgExist(msgID uint16) bool {
	_, ok := def.id2Info[msgID]
	return ok
}

// GetMsgIDByName 通过名字获取ID号
func (def *MsgDef) GetMsgIDByName(msgName string) (uint16, error) {

	info, ok := def.name2Info[msgName]
	if !ok {
		log.Debug(string(debug.Stack()))
		return 0, fmt.Errorf("不存在消息号 Name: %s", msgName)
	}

	return info.id, nil
}

// RegMsg 注册消息
func (def *MsgDef) RegMsg(msgID uint16, msgBody interface{}) {
	_, ok := def.id2Info[msgID]

	if ok {
		log.Warn("消息ID已经存在 ", msgID)
		return
	}

	msgName := reflect.TypeOf(msgBody).Elem().Name()

	_, ok = def.name2Info[msgName]

	if ok {
		log.Warn("消息名称已经存在  ", msgName)
		return
	}

	_, ok = msgBody.(IMsg)

	if !ok {
		log.Warn("注册的消息对象未实现 Imsg接口 ", msgName)
		return
	}

	info := &MsgInfo{
		msgID,
		msgName,
		reflect.TypeOf(msgBody),
	}

	def.id2Info[msgID] = info
	def.name2Info[msgName] = info
}

var msgDefInst *MsgDef

// GetMsgDef 获取消息定义对象的全局实例
func GetMsgDef() *MsgDef {

	if msgDefInst == nil {
		msgDefInst = &MsgDef{
			make(map[uint16]*MsgInfo),
			make(map[string]*MsgInfo),
		}
		msgDefInst.Init()
	}

	return msgDefInst
}

// Init 初始化
func Init() {
	if msgDefInst == nil {
		msgDefInst = &MsgDef{
			make(map[uint16]*MsgInfo),
			make(map[string]*MsgInfo),
		}
		msgDefInst.Init()
	}
}
