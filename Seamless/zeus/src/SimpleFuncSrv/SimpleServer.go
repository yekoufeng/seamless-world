package main

import (
	//"common"
	"flag"
	"fmt"
	"strconv"
	"time"
	//"zeus/env"
	"zeus/iserver"
	"zeus/server"
)

type SimpleServer struct {
	iserver.IServer
}

var srvInst *SimpleServer

func GetSimpleSrvInst() *SimpleServer {
	if srvInst == nil {
		srvInst = &SimpleServer{}

		// 主程序初始化时，从 flag中或配置文件读取到的配置信息
		srvID := flag.Uint64("id", 30000, "Server ID")
		flag.Parse()
		pmin, _ := strconv.ParseUint(env.Get("Simple", "PortMin"), 10, 32)
		pmax, _ := strconv.ParseUint(env.Get("Simple", "PortMax"), 10, 32)
		innerAddr := env.Get("Simple", "InnerAddr")
		innerPort := server.GetValidSrvPort(int(pmin), int(pmax))
		fps, _ := strconv.ParseUint(env.Get("Simple", "FPS"), 10, 32)
		srvInst.IServer = server.NewServer(3, *srvID, innerAddr+":"+innerPort, "", int(fps), srvInst)
	}

	return srvInst
}

// Init 初始化
func (srv *SimpleServer) Init() error {

	common.InitMsg()

	srv.RegProtoType("Player", &SimpleUser{})
	srv.RegProtoType("TestSpace", &TestSpace{})

	srv.AddDelayCall(func() {
		if _, err := srv.CreateEntityAll("TestSpace", 50000, 0); err != nil {
			fmt.Println(err)
		}
	}, 1*time.Second)

	srv.AddDelayCall(func() {
		if _, err := srv.CreateEntityAll("Player", server.GetEntityTempID(), 1000); err != nil {
			fmt.Println(err)
		}
	}, 3*time.Second)

	return nil
}

// MainLoop 逻辑帧每一帧都会调用
func (srv *SimpleServer) MainLoop() {

}

// Destroy 退出时调用
func (srv *SimpleServer) Destroy() {

}
