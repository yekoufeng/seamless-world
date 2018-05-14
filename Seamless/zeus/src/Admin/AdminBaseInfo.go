package main

import (
	"zeus/serverMgr"

	"zeus/iserver"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

/*
 服务器基本信息查询
*/

// ServerList /servers处理函数
func (srv *AdminServer) ServerList(w rest.ResponseWriter, r *rest.Request) {
	if err := srv.initBaseInfo(); err != nil {
		log.Error(err)
	}
	w.WriteJson(srv.serverMap)
}

func (srv *AdminServer) initBaseInfo() error {
	serverList, err := serverMgr.GetServerMgr().GetServerList()
	if err != nil {
		return err
	}

	srv.serverMapRW.Lock()
	defer srv.serverMapRW.Unlock()
	srv.serverMap = make(map[uint64]*iserver.ServerInfo)
	for _, v := range serverList {
		srv.serverMap[v.ServerID] = v
	}

	return nil
}
