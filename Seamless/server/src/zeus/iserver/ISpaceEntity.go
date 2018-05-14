package iserver

import (
	"zeus/linmath"
	"zeus/msgdef"
)

// ICellEntity 代表空间中的一个实体
type ICellEntity interface {
	IEntity
	IPos
	IClientBroadcaster
	IAOIEntity
	SetClient(ISess)
}

// ICellEntityPropsGetter 获取属性数据
type ICellEntityPropsGetter interface {
	GetAOIProp() (int, []byte)
}

// IPos 拥有位置信息的接口
type IPos interface {
	SetPos(pos linmath.Vector3)
	GetPos() linmath.Vector3

	SetRota(linmath.Vector3)
	GetRota() linmath.Vector3
}

/*
// IPosValidate 位置验证
type IPosValidate interface {
	IsPosValid(linmath.Vector3) bool
}
*/

// IPosChange 位置变化回调
type IPosChange interface {
	OnPosChange(curPos, origPos linmath.Vector3)
}

// IMover 寻路移动能力
type IMover interface {
	SetSpeed(speed float32)
	Move(destPos linmath.Vector3) bool
	StopMove()
	IsMoving() bool

	PauseNav()
	ResumeNav()
}

// IClientBroadcaster AOI广播
type IClientBroadcaster interface {
	CastMsgToAllClient(msgdef.IMsg)
	CastMsgToMe(msgdef.IMsg)
	CastMsgToAllClientExceptMe(msgdef.IMsg)
	CastMsgToRangeExceptMe(center *linmath.Vector3, radius int, msg msgdef.IMsg)
	CastMsgToCenterExceptMe(center *linmath.Vector3, radius int, msg msgdef.IMsg)

	CastRPCToAllClient(methodName string, args ...interface{})
	CastRPCToMe(methodName string, args ...interface{})
	CastRPCToAllClientExceptMe(methodName string, args ...interface{})

	BroadcastEvent(event string, args ...interface{})
	BroadcastEventExceptMe(event string, args ...interface{})
}

// IAOIEntity  AOI实体类型查询
type IAOIEntity interface {
	IsWatcher() bool
	// IsMarker() bool
}

type ISyncToGhosts interface {
	SyncToGhosts(msgdef.IMsg)
}

type ISendMsgToReal interface {
	SendMsgToReal(msg msgdef.IMsg)
}
