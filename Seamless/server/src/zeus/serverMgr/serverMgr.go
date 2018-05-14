package serverMgr

import (
	"errors"
	"zeus/dbservice"
	"zeus/iserver"

	log "github.com/cihub/seelog"
)

/*
 服务器管理, 目前由redis实现

 redis中的数据结构: 哈希表
 Key: "server:" + ServerID
 ServerID    	服务器ID
 Type   		服务器类型
 OuterAddress,	服务器外网地址
 InnerAddress,	服务器内网地址
 Load,			服务器当前负载
 Token,			Token

*/

// ServerMgr 服务器管理类
type ServerMgr struct {
	*LoadBalancer
}

var mgr *ServerMgr

// GetServerMgr 获取服务器管理实例
func GetServerMgr() *ServerMgr {
	if mgr == nil {
		mgr = &ServerMgr{}
		mgr.LoadBalancer = NewLoadBalancer()
	}
	return mgr
}

// Unregister 将服务器信息从redis中删除
func (mgr *ServerMgr) Unregister(server *iserver.ServerInfo) error {
	if err := dbservice.ServerUtil(server.ServerID).Delete(); err != nil {
		return err
	}
	//iserver.GetSrvInst().FireEvent(serverInfoChannel)
	return nil
}

var errParamNil = errors.New("Param is nil")

// RegState 注册服务器信息
func (mgr *ServerMgr) RegState(server *iserver.ServerInfo) error {
	if server == nil {
		return errParamNil
	}

	util := dbservice.ServerUtil(server.ServerID)

	if util.IsExist() {
		log.Error("Server ID is duplicate!!!!! ", server.ServerID)
		// panic("server is is dupblicate")
	}

	if err := util.Register(server); err != nil {
		return err
	}
	return nil
}

// Update 更新服务器信息
func (mgr *ServerMgr) Update(server *iserver.ServerInfo) error {
	if server == nil {
		return errParamNil
	}

	return dbservice.ServerUtil(server.ServerID).Update(server)
}

// VerifyServer 验证服务器连接有效性
func (mgr *ServerMgr) VerifyServer(uid uint64, token string) bool {
	t, err := dbservice.ServerUtil(uid).GetToken()
	if err != nil {
		log.Error(err)
		return false
	}

	if t != token {
		return false
	}

	return true
}
