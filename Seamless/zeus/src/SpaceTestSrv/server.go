package main

import (
	//"common"
	"flag"
	"strconv"
	//"zeus/env"
	"zeus/iserver"
	"zeus/server"
)

type Server struct {
	iserver.IServer
}

var srvInst *Server

func GetSrvInst() *Server {
	if srvInst == nil {
		srvInst = &Server{}

		// 主程序初始化时，从 flag中或配置文件读取到的配置信息
		srvID := flag.Uint64("id", 40000, "Server ID")
		flag.Parse()
		pmin, _ := strconv.ParseUint(env.Get("Simple", "PortMin"), 10, 32)
		pmax, _ := strconv.ParseUint(env.Get("Simple", "PortMax"), 10, 32)
		innerAddr := env.Get("Simple", "InnerAddr")
		innerPort := server.GetValidSrvPort(int(pmin), int(pmax))
		fps, _ := strconv.ParseUint(env.Get("Simple", "FPS"), 10, 32)
		srvInst.IServer = server.NewServer(4, *srvID, innerAddr+":"+innerPort, "", int(fps), srvInst)
	}

	return srvInst
}

// Init 初始化
func (srv *Server) Init() error {
	common.InitMsg()
	//注册可创建的实例
	srv.RegProtoType("Player", &SpaceUser{})
	srv.RegProtoType("TestSpace", &TestSpace{})

	return nil
}

// MainLoop 逻辑帧每一帧都会调用
func (srv *Server) MainLoop() {

}

// Destroy 退出时调用
func (srv *Server) Destroy() {

}
