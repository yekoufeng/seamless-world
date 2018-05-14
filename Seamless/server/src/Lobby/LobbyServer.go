package main

import (
	"Lobby/online"
	"common"
	"errors"
	"fmt"
	"strings"
	"time"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/server"
	"zeus/tlog"
	"zeus/tsssdk"

	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var srvInst *LobbySrv

// GetSrvInst 获取服务器全局实例
func GetSrvInst() *LobbySrv {
	if srvInst == nil {
		common.InitMsg()

		srvInst = &LobbySrv{}

		srvID := uint64(viper.GetInt("Lobby.FlagId"))
		pmin := viper.GetInt("Lobby.PortMin")
		pmax := viper.GetInt("Lobby.PortMax")

		srvInst.sqlUser = viper.GetString("Lobby.MySQLUser")
		srvInst.sqlPwd = viper.GetString("Lobby.MySQLPwd")
		srvInst.sqlAddr = viper.GetString("Lobby.MySQLAddr")
		srvInst.sqlDB = viper.GetString("Lobby.MySQLDB")
		srvInst.sqlTB = viper.GetString("Lobby.MySQLTable")
		srvInst.msdkAddr = viper.GetString("Config.MSDKAddr")

		innerPort := server.GetValidSrvPort(pmin, pmax)
		innerAddr := viper.GetString("Lobby.InnerAddr")
		fps := viper.GetInt("Lobby.FPS")
		srvInst.IServer = server.NewServer(common.ServerTypeLobby, srvID, innerAddr+":"+innerPort, "", fps, srvInst)

		srvInst.onlineCnter = make(map[string]*online.Cnter)

		tlogAddr := viper.GetString("Config.TLogAddr")
		if tlogAddr != "" {
			if err := tlog.ConfigRemoteAddr(tlogAddr); err != nil {
				log.Error(err)
			}
		}

		log.Info("Lobby Init")
		log.Info("ServerID:", srvID)
		log.Info("InnerAddr:", innerAddr+":"+innerPort)
	}

	return srvInst
}

// LobbySrv 大厅服务器
type LobbySrv struct {
	iserver.IServer

	sqlUser  string
	sqlPwd   string
	sqlAddr  string
	sqlDB    string
	sqlTB    string
	msdkAddr string

	// tlog 相关
	onlineCnter map[string]*online.Cnter
	stateLogger *tlog.StateLogger
}

// Init 初始化
func (srv *LobbySrv) Init() error {
	srv.RegProtoType("Player", &LobbyUser{}, true)

	srv.RegMsgProc(&LobbySrvMsgProc{srv: srv})

	srv.stateLogger = tlog.NewStateLogger(srv.GetSrvAddr(), 0, 5*time.Minute)
	srv.stateLogger.Start()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if strings.Contains(e.Name, "server.json") {
			if err := viper.ReadInConfig(); err != nil {
				log.Error(err)
			} else {
				log.Info("Reload server.json success")
			}
		}
	})

	tsssdk.Init(GetSrvInst().GetSrvID())
	srv.regBroadcaster()
	log.Info("Lobby Start")

	return nil
}

// MainLoop 逻辑帧每一帧都会调用
func (srv *LobbySrv) MainLoop() {
}

// Destroy 退出时调用
func (srv *LobbySrv) Destroy() {
	srv.stopOnlineUpdate()
	srv.stateLogger.Stop()
	tsssdk.Destroy()
	log.Info("Lobby Shutdown")
}

func (srv *LobbySrv) stopOnlineUpdate() {
	for _, cnt := range srv.onlineCnter {
		if err := cnt.Stop(); err != nil {
			log.Error(err)
		}
	}
}

var errGameAppIDEmpty = errors.New("GameAppID empty")

//login更新在线人数统计Cnter
func (srv *LobbySrv) loginCnt(gameApp string, platID int) error {
	if platID != 0 && platID != 1 {
		return fmt.Errorf("Login PlatID Error %d", platID)
	}
	if gameApp == "" {
		return errGameAppIDEmpty
	}

	cnt, ok := srv.onlineCnter[gameApp]
	if !ok {
		var err error
		srv.onlineCnter[gameApp], err = online.NewCnter(srv.sqlUser, srv.sqlPwd, srv.sqlAddr,
			srv.sqlDB, srv.sqlTB, gameApp, 0, srv.GetSrvID())
		if err != nil {
			return err
		}

		cnt = srv.onlineCnter[gameApp]
		cnt.Start()
	}

	cnt.ReportOnline(platID, 1)
	return nil
}

var errLogoutWithoutLogin = errors.New("Logout without login")

//logout更新在线人数统计Cnter
func (srv *LobbySrv) logoutCnt(gameApp string, platID int) error {
	if platID != 0 && platID != 1 {
		return fmt.Errorf("Logout PlatID Error %d", platID)
	}
	if gameApp == "" {
		return errGameAppIDEmpty
	}

	cnt, ok := srv.onlineCnter[gameApp]
	if !ok {
		return errLogoutWithoutLogin
	}

	cnt.ReportOnline(platID, -1)
	return nil
}

func (srv *LobbySrv) regBroadcaster() {
	if 0 == srv.AddListener(iserver.BroadcastChannel, srv, "BroadcastMsg") {
		panic("Register broadcast channel failed")
	}
	if 0 == srv.AddListener(iserver.RPCChannel, srv, "RPCClients") {
		panic("Register rpc channel failed")
	}
}

// BroadcastMsg 广播消息到所有客户端
func (srv *LobbySrv) BroadcastMsg(msg msgdef.IMsg) {
	log.Info("BroadcastMsg ", msg)

	srv.TravsalEntity("Player", func(o iserver.IEntity) {
		o.Post(iserver.ServerTypeClient, msg)
	})
}

// RPCClients 调用所有客户端的RPC消息
func (srv *LobbySrv) RPCClients(method string, args ...interface{}) {
	log.Info("RPCClients ", method, args)

	srv.TravsalEntity("Player", func(o iserver.IEntity) {
		o.RPC(iserver.ServerTypeClient, method, args...)
	})
}

// InvalidHandle 服务器不可用时处理回调
func (srv *LobbySrv) InvalidHandle(entityID uint64) {
	e := srv.GetEntity(entityID)
	if e.GetType() == "Player" {

		log.Warn("服务器把你炸下线了 ", e)

		srv.DestroyEntityAll(entityID)
	}
}

func (srv *LobbySrv) OnServerConnect(srvID uint64, serverType uint8) {

}

func (srv *LobbySrv) GetEntities(cellID uint64) iserver.IEntities {
	e := srv.GetEntity(cellID)

	if e == nil {
		return nil
	}

	return e.(iserver.IEntities)
}
