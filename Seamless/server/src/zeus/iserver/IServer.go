package iserver

import (
	"time"
	"zeus/admin"
)

const (
	// BroadcastChannel 广播消息频道
	BroadcastChannel = "broadcast"

	// RPCChannel 频道
	RPCChannel = "rpc"
)

// IServer 基础Server提供的接口
// 接口都是非线程安全，注意调用线程
type IServer interface {
	ISrvNet
	IEntities
	IEntityProto
	IUIDFetcher
	admin.IAdmin

	GetFrameDeltaTime() time.Duration
	GetStartupTime() time.Time
	Run()
	GetLoad() int

	IsSrvValid() bool
	HandlerSrvInvalid(entityID uint64)

	OnServerConnect(srvID uint64, serverType uint8)
}

var srvInst IServer

// GetSrvInst 获取当前服务器
func GetSrvInst() IServer {
	return srvInst
}

// SetSrvInst 设置单例
func SetSrvInst(srv IServer) {
	srvInst = srv
}
