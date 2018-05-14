package iserver

import "zeus/msgdef"

// IEntityProxy 实体代理接口
type IEntityProxy interface {
	Post(msg msgdef.IMsg) error
	RPC(srvType uint8, methodName string, args ...interface{}) error
}
