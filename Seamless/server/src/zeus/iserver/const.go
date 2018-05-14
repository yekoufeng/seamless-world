package iserver

const (
	_ = iota
	// ServerTypeGateway 网关服
	ServerTypeGateway

	// ServerTypeClient 客户端
	ServerTypeClient

	// ServerTypeSpace  带场景的服务器
	ServerTypeSpace
)

const (
	//MaxServerID 最大的服务器ID号
	MaxServerID = 1<<27 - 1
)
