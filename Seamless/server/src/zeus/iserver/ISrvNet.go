package iserver

import (
	"zeus/msgdef"
)

// ISrvNet 服务器的网状结构
type ISrvNet interface {
	GetSrvType() uint8
	GetSrvID() uint64
	GetSrvAddr() string

	PostMsgToSrv(srvID uint64, msg msgdef.IMsg) error
	PostMsgToCell(srvID uint64, cellID uint64, msg msgdef.IMsg) error

	// GetSrvInfo(srvID uint64) *ServerInfo
	GetCurSrvInfo() *ServerInfo
	GetSrvIDBySrvType(srvType uint8) (uint64, error)
}
