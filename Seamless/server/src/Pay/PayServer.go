package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

var srvInst *Server

// GetSrvInst 获取服务器全局实例
func GetSrvInst() *Server {
	if srvInst == nil {
		srvInst = &Server{}
		srvInst.addr = "0.0.0.0"
		srvInst.port = "3400"
	}

	return srvInst
}

// Server 中心服务器
type Server struct {
	addr string
	port string
}

// Start 启动服务器
func (srv *Server) Start() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	router, err := rest.MakeRouter(
		rest.Post("/querybalance", srv.QueryBalance),
		rest.Post("/buyprop", srv.BuyProp),
	)
	if err != nil {
		log.Error(err)
	}
	api.SetApp(router)
	err = http.ListenAndServe("0.0.0.0:"+srv.port, api.MakeHandler())
	if err != nil {
		log.Error("listen error", err)
		return
	}
}
