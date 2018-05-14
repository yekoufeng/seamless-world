package server

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"zeus/admin"
	"zeus/dbservice"
	"zeus/entity"
	"zeus/iserver"
	"zeus/serverMgr"

	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// IServerCtrl 后代服务器类需要继承的接口
type IServerCtrl interface {
	Init() error
	MainLoop()
	Destroy()

	GetLoad() int

	OnServerConnect(srvID uint64, serverType uint8)

	GetEntities(cellID uint64) iserver.IEntities
}

// Server 服务器
type Server struct {
	*SrvNet
	*serverMgr.LoadUpdater
	*entity.ProtoType
	*entity.Entities
	*IDFetcher
	srvCtrl IServerCtrl

	ticker         *time.Ticker
	frameDeltaTime time.Duration
	startupTime    time.Time
	loopStopC      chan bool
	srvStopC       chan bool
}

// NewServer 创建一个新的服务器
func NewServer(srvType uint8, srvID uint64, addr string, outerAddr string, fps int, srvCtrl IServerCtrl) iserver.IServer {

	log.Debug("srvType = ", srvType)
	if iserver.GetSrvInst() != nil {
		log.Error("服务器已经存在，一个应用只能创建一个服务器实例")
		return nil
	}

	if srvID > iserver.MaxServerID {
		log.Error("serverID 超过了最大ID号,服务器最大ID为 ", iserver.MaxServerID)
		return nil
	}

	deltaTime := time.Millisecond * time.Duration(1000/fps)

	srv := &Server{
		SrvNet: NewSrvNet(srvType, srvID, addr, outerAddr),

		ProtoType: entity.NewProtoType(),

		Entities:  entity.NewEntities(true),
		IDFetcher: newIDFetcher(srvID),
		srvCtrl:   srvCtrl,

		frameDeltaTime: deltaTime,
		startupTime:    time.Now(),
		ticker:         time.NewTicker(deltaTime),
		loopStopC:      make(chan bool),
		srvStopC:       make(chan bool),
	}

	//srv.Admin = admin.NewAdminApp(srv)
	srv.LoadUpdater = serverMgr.NewLoadUpdater(srvCtrl, srvID, 5*time.Second)

	iserver.SetSrvInst(srv)

	return srv
}

func (srv *Server) init() error {

	if err := srv.SrvNet.init(); err != nil {
		return err
	}

	srv.SrvNet.regMsgProc(&serverMsgProc{srv: srv})
	srv.SrvNet.regMsgProc(srv.srvCtrl)

	if err := srv.IDFetcher.init(); err != nil {
		return err
	}

	srv.Entities.Init()
	return srv.srvCtrl.Init()
}

func (srv *Server) destroy() {
	srv.LoadUpdater.Stop()

	srv.srvCtrl.Destroy()
	// srv.Entities.Destroy()
	srv.Entities.SyncDestroy()
	srv.SrvNet.destroy()

	srv.ticker.Stop()
}

//Close 关闭服务器
func (srv *Server) Close() {
	srv.loopStopC <- true
	close(srv.loopStopC)
}

// Run 逻辑入口
func (srv *Server) Run() {

	if err := srv.init(); err != nil {
		panic(err)
	}

	go srv.doLoop()

	//go srv.StartConsole()
	// if srv.Admin.GetConsolePort() != 0 {
	// 	go srv.StartHTTP()
	// }
	srv.LoadUpdater.Start() //启动负载更新器

	<-srv.srvStopC

	srv.destroy()
}

func (srv *Server) doLoop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	for {
		select {
		case <-srv.loopStopC:
			srv.srvStopC <- true
			close(srv.srvStopC)
			return
		case sig := <-c:
			log.Info(sig.String())
			srv.srvStopC <- true
			close(srv.srvStopC)
			return
		case <-srv.ticker.C:
			srv.MainLoop()
		}
	}
}

// MainLoop 主循环
func (srv *Server) MainLoop() {
	srv.SrvNet.MainLoop() //网关
	srv.srvCtrl.MainLoop()
	srv.Entities.MainLoop()
	//srv.Admin.ProcCmds() //命令行处理
}

// ConfigHTTPAdmin 配置控制台相关
func (srv *Server) ConfigHTTPAdmin(addr string, port uint64, admin string) {
	//srv.Admin.ConfigHTTPAdmin(addr, port, admin)
	// srv.SrvNet.SetConsole(port)
}

// HandleCommand 控制台接口
func (srv *Server) HandleCommand(c []string) *admin.CmdResp {
	switch c[0] {
	case "status":
		return srv.doStatusCmd(c)
	default:
		if iCmd, ok := srv.srvCtrl.(admin.IServerCommand); ok {
			return iCmd.HandleCommand(c)
		}
	}

	return nil
}

func (srv *Server) doStatusCmd(cmd []string) *admin.CmdResp {
	resp := &admin.CmdResp{}

	if len(cmd) == 2 {
		targetStatus, _ := strconv.ParseInt(cmd[1], 10, 32)
		// srv.SetStatus(int(targetStatus))
		str := fmt.Sprintf("切换服务器状态至%d", targetStatus)
		log.Warn(str)
		resp.Result = 0
		resp.ResultStr = str
	} else {
		resp.Result = -1
		resp.ResultStr = "参数错误"
	}

	return resp
}

// GetLoad 获取服务器负载信息, 取CPU和内存的大值
func (srv *Server) GetLoad() int {
	var c int
	var vm int

	if loads, err := cpu.Percent(0, false); err == nil {
		if len(loads) > 0 {
			c = int(loads[0])
		}
	} else {
		log.Error(err)
		c = 100
	}

	if memorys, err := mem.VirtualMemory(); err == nil {
		vm = int(memorys.UsedPercent)
	} else {
		log.Error(err)
		vm = 100
	}

	if c > vm {
		return c
	}

	return vm
}

// GetFrameDeltaTime 获取每帧间的间隔
func (srv *Server) GetFrameDeltaTime() time.Duration {
	return srv.frameDeltaTime
}

// GetStartupTime 获取服务器启动时间
func (srv *Server) GetStartupTime() time.Time {
	return srv.startupTime
}

// GetCurSrvInfo 获取当前服务器信息
// func (srv *Server) GetCurSrvInfo() *iserver.ServerInfo {
// 	return srv.SrvNet.GetSrvInfo(srv.GetSrvID())
// }

// IsSrvValid 服务是否可用
func (srv *Server) IsSrvValid() bool {
	return dbservice.DBValid && dbservice.SrvRedisValid && dbservice.SingletonRedisValid
}

type iInvalidHandler interface {
	InvalidHandle(entityID uint64)
}

// HandlerSrvInvalid 处理服务不可用
func (srv *Server) HandlerSrvInvalid(entityID uint64) {
	if iih, ok := srv.srvCtrl.(iInvalidHandler); ok {
		iih.InvalidHandle(entityID)
	}
}

func (srv *Server) OnServerConnect(srvID uint64, serverType uint8) {

	srv.srvCtrl.OnServerConnect(srvID, serverType)

}

func (srv *Server) GetEntities(cellID uint64) iserver.IEntities {
	return srv.srvCtrl.GetEntities(cellID)
}
