package iserver

import (
	"zeus/dbservice"
	"zeus/msgdef"
)

// IEntity 实体接口
type IEntity interface {
	// 消息传递
	Post(srvType uint8, msg msgdef.IMsg) error
	// 快速消息调用
	RPC(srvType uint8, methodName string, args ...interface{}) error
	//RPCOther(srvType uint8, srcEntityID uint64, methodName string, args ...interface{}) error

	// 获取实体类型
	GetType() string
	GetID() uint64
	GetDBID() uint64
	GetRealPtr() interface{}
	GetInitParam() interface{}
	GetSrvIDS() map[uint8]*dbservice.EntitySrvInfo
	IsGhost() bool

	GetProxy() IEntityProxy

	LeaveCell()
	GetCellID() uint64
	IsOwnerCellEntity() bool
	IsCellEntity() bool
}

const (
	Entity_State_Init    = 0
	Entity_State_Loop    = 1
	Entity_State_Destroy = 2
	Entity_State_InValid = 3
)

// IEntityStateGetter 获取Entity状态
type IEntityStateGetter interface {
	GetEntityState() uint8
}
