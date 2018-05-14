package main

import (
	"common"
	"time"
	"zeus/iserver"
	"zeus/server"
	"zeus/tlog"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var srvInst *Server

// GetSrvInst 获取服务器全局实例
func GetSrvInst() *Server {
	if srvInst == nil {
		common.InitMsg()

		srvInst = &Server{}

		srvID := uint64(viper.GetInt("Center.FlagId"))
		pmin := viper.GetInt("Center.PortMin")
		pmax := viper.GetInt("Center.PortMax")

		srvInst.idipAddr = viper.GetString("Center.IDIPAddr")
		srvInst.idipPort = viper.GetString("Center.IDIPPort")

		innerPort := server.GetValidSrvPort(pmin, pmax)
		innerAddr := viper.GetString("Center.InnerAddr")
		fps := viper.GetInt("Center.FPS")
		srvInst.IServer = server.NewServer(common.ServerTypeCenter, srvID, innerAddr+":"+innerPort, "", fps, srvInst)

		tlogAddr := viper.GetString("Config.TLogAddr")
		if tlogAddr != "" {
			if err := tlog.ConfigRemoteAddr(tlogAddr); err != nil {
				log.Error(err)
			}
		}

		log.Info("Center Init")
		log.Info("ServerID:", srvID)
		log.Info("InnerAddr:", innerAddr+":"+innerPort)
	}

	return srvInst
}

// Server 中心服务器
type Server struct {
	iserver.IServer

	idipAddr string
	idipPort string

	annuonceMgr *AnnuonceMgr
	ht          *HttpService

	// tlog相关
	stateLogger *tlog.StateLogger
}

// Init 初始化
func (srv *Server) Init() error {
	srv.annuonceMgr = NewAnnuonceMgr(srv)
	srv.ht = NewHttpService(srv)
	srv.stateLogger = tlog.NewStateLogger(srv.GetSrvAddr(), 0, 5*time.Minute)
	srv.stateLogger.Start()

	srv.RegMsgProc(&CenterSrvMsgProc{srv: srv})

	go srv.startIDIP()

	log.Info("Center Start")
	return nil
}

// MainLoop 逻辑帧每一帧都会调用
func (srv *Server) MainLoop() {
}

// Destroy 退出时调用
func (srv *Server) Destroy() {
	srv.stateLogger.Stop()

	log.Info("Center Shutdown")
}

func (srv *Server) GetEntities(cellID uint64) iserver.IEntities {
	return nil
}