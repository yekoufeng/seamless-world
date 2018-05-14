package iserver

import "zeus/msgdef"
import "zeus/msghandler"

// ISess 代表一个网络连接
type ISess interface {
	msghandler.IMsgHandlers

	SetMsgHandler(handler msghandler.IMsgHandlers)
	Send(msgdef.IMsg)
	SendRaw([]byte)

	Start()
	Close()
	IsClosed() bool

	RemoteAddr() string

	IsVertified() bool
	SetVertify()

	SetID(uint64)
	GetID() uint64

	SetServerType(uint8)
	GetServerType() uint8

	Touch()

	FetchBacklog(ISess)
	FlushBacklog()
}
