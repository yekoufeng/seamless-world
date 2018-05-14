package main

import (
	"io/ioutil"
	"time"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/server"
	"zeus/sess"
	"zeus/tlog"

	"runtime/debug"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var srvInst *GatewaySrv

// GetSrvInst 获取服务器全局实例
func GetSrvInst() *GatewaySrv {
	if srvInst == nil {
		srvInst = &GatewaySrv{}
		msgdef.Init()

		srvID := uint64(viper.GetInt("Gateway.FlagId"))
		fps := viper.GetInt("Gateway.FPS")
		pmin := viper.GetInt("Gateway.PortMin")
		pmax := viper.GetInt("Gateway.PortMax")

		/*
			配置中, OuterAddr+OuterPort写入到redis中, 作为服务器的对外地址
			OuterListen+随机出来的端口作为实际对外监听端口
			当配置中OuterPort为0时, OuterPort就是随机端口
		*/

		outpmin := viper.GetInt("Gateway.OuterPortMin")
		outpmax := viper.GetInt("Gateway.OuterPortMax")

		innerPort := server.GetValidSrvPort(pmin, pmax)
		for {
			srvInst.outerListenPort = server.GetValidSrvPort(int(outpmin), int(outpmax))
			if srvInst.outerListenPort != innerPort {
				break
			}
		}
		innerAddr := viper.GetString("Gateway.InnerAddr")
		outerAddr := viper.GetString("Gateway.OuterAddr")
		outerPort := srvInst.outerListenPort
		if viper.GetString("Gateway.OuterPort") != "0" {
			outerPort = viper.GetString("Gateway.OuterPort")
			srvInst.outerListenPort = outerPort
		}
		srvInst.IServer = server.NewServer(iserver.ServerTypeGateway, srvID, innerAddr+":"+innerPort, outerAddr+":"+outerPort, fps, srvInst)

		// consoleAddr := env.Get("Gateway", "ConsoleAddr")
		// consolePort, _ := strconv.ParseUint(env.Get("Gateway", "Console"), 10, 32)
		// admin := env.Get("Admin", "Domain")
		// srvInst.ConfigHTTPAdmin(consoleAddr, consolePort, admin)

		tlogAddr := viper.GetString("Config.TLogAddr")
		if tlogAddr != "" {
			if err := tlog.ConfigRemoteAddr(tlogAddr); err != nil {
				log.Error(err)
			}
		}

		log.Info("Gateway Init")
		log.Info("ServerID:", srvID)
		log.Info("InnerAddr:", innerAddr+":"+innerPort)
		log.Info("OuterAddr:", outerAddr+":"+outerPort)

	}

	return srvInst
}

// GatewaySrv 网关服务器
type GatewaySrv struct {
	iserver.IServer
	outerListenPort string
	protoSync       *msgdef.ProtoSync
	clientSrv       sess.IMsgServer

	stateLogger *tlog.StateLogger
}

// Init 初始化
func (srv *GatewaySrv) Init() error {

	log.Debug(string(debug.Stack()))

	srv.RegProtoType("Player", &GateUser{}, true)
	outlisten := viper.GetString("Gateway.OuterListen")
	maxConns := viper.GetInt("Gateway.MaxConns")
	srv.clientSrv = sess.NewMsgServer("tcp", outlisten+":"+srv.outerListenPort, maxConns)
	srv.clientSrv.RegMsgProc(&GatewaySrvMsgProc{})
	// srv.clientSrv.SetForwardMode()

	encryptEnabled := viper.GetBool("Config.EncryptEnabled")
	if encryptEnabled {
		srv.clientSrv.SetEncryptEnabled()
	}

	if err := srv.clientSrv.Start(); err != nil {
		log.Error("start gateway error", err)
		return err
	}

	// 读取Proto消息二进制数据
	if err := srv.genProtoSync(); err != nil {
		log.Error("start gateway error", err)
		return err
	}

	srv.stateLogger = tlog.NewStateLogger(srv.GetSrvAddr(), 0, 5*time.Minute)
	srv.stateLogger.Start()

	log.Info("OuterListen:", outlisten+":"+srv.outerListenPort)
	log.Info("Gateway Start")
	return nil
}

// MainLoop 逻辑帧每一帧都会调用
func (srv *GatewaySrv) MainLoop() {
	//log.Debug(string(debug.Stack()))

}

// Destroy 退出时调用
func (srv *GatewaySrv) Destroy() {
	GetUserMgr().forceLogout()
	srv.stateLogger.Stop()

	srv.RemoveListenerByObjInst(srv)
	log.Info("Gateway Shutdown")
}

func (srv *GatewaySrv) genProtoSync() error {
	srv.protoSync = &msgdef.ProtoSync{}
	protoFile := viper.GetString("Gateway.ProtoFile")
	data, err := ioutil.ReadFile(protoFile)
	if err != nil {
		return err
	}
	srv.protoSync.Data = data
	return nil
}

func (srv *GatewaySrv) OnServerConnect(srvID uint64, serverType uint8) {

	log.Debug("GatewaySrv OnServerConnect.... serverType is ", serverType)

}

func (srv *GatewaySrv) GetEntities(cellID uint64) iserver.IEntities {
	return nil
}
