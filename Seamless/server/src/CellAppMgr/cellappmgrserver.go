package main

import (
	"common"
	_ "net/http/pprof"
	"protoMsg"
	"time"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/server"
	"zeus/tlog"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	//"math/rand"
)

var srvInst *CellAppMgrSrv
var gFPS int

//var spaceID uint64 = 0
//var cellID  uint64 = 0

type CellAppMgrSrv struct {
	iserver.IServer
	//所有的cellapp管理, cellID和cellapp的对应表
	cellapps map[uint64]*CellApp

	ReportTicker *time.Ticker

	//大地图ID，所有的cellapp共用
	spaceID uint64
}

// GetSrvInst 获取服务器全局实例
func GetSrvInst() *CellAppMgrSrv {
	if srvInst == nil {
		common.InitMsg()

		srvInst = &CellAppMgrSrv{}

		srvID := uint64(viper.GetInt("CellAppMgr.FlagId"))
		pmin := viper.GetInt("CellAppMgr.PortMin")
		pmax := viper.GetInt("CellAppMgr.PortMax")

		innerPort := server.GetValidSrvPort(pmin, pmax)
		innerAddr := viper.GetString("CellAppMgr.InnerAddr")

		opmin := viper.GetInt("CellAppMgr.OuterPortMin")
		opmax := viper.GetInt("CellAppMgr.OuterPortMax")

		fps := viper.GetInt("CellAppMgr.FPS")
		gFPS = fps

		/*
			配置中, OuterAddr+OuterPort写入到redis中, 作为服务器的对外地址
			OuterListen+随机出来的端口作为实际对外监听端口
			当配置中OuterPort为0时, OuterPort就是随机端口
		*/
		//outerListen := viper.GetString("Room.OuterListen")
		listenPort := server.GetValidSrvPort(opmin, opmax)
		outerAddr := viper.GetString("CellAppMgr.OuterAddr")
		outerPort := listenPort
		if viper.GetString("CellAppMgr.OuterPort") != "0" {
			outerPort = viper.GetString("CellAppMgr.OuterPort")
			listenPort = outerPort
		}

		srvInst.IServer = server.NewServer(common.ServerTypeCellAppMgr, srvID, innerAddr+":"+innerPort, outerAddr+":"+outerPort, fps, srvInst)

		//protocal := viper.GetString("CellAppMgr.SpaceProtocal")
		//maxConns := viper.GetInt("CellAppMgr.MaxConns")

		tlogAddr := viper.GetString("Config.TLogAddr")
		if tlogAddr != "" {
			if err := tlog.ConfigRemoteAddr(tlogAddr); err != nil {
				log.Error(err)
			}
		}

		log.Info("CellAppMgr Init")
		log.Info("ServerID:", srvID)
		log.Info("InnerAddr:", innerAddr+":"+innerPort)

	}
	return srvInst
}

// Init 初始化
func (srv *CellAppMgrSrv) Init() error {

	srv.RegMsgProc(&CellAppMgrSrvMsgProc{srv: srv})
	srv.cellapps = make(map[uint64]*CellApp)

	srv.spaceID = 0
	srv.allocSpaceCell()

	log.Info("CellAppMgrSrv Start")

	return nil
}

func (srv *CellAppMgrSrv) MainLoop() {

	//定时检测每个cellapp的压力
	//判定这个cellapp上的压力，如果压力达到交出cell的条件，就交给另外一个合适的cellapp去处理

	//如果别的cellapp里有负载比较轻的

	var maxloadServer []uint64
	var minloadServer []uint64

	//遍历一遍, 找出超过压力警戒线的cellapp, 找出压力比较低的几个cellapp
	for _, cellapp := range srv.cellapps {
		if cellapp.getOverLoad() > 80 {
			maxloadServer = append(maxloadServer, cellapp.getID())
		}
		if cellapp.getOverLoad() < 20 {
			minloadServer = append(minloadServer, cellapp.getID())
		}
	}

	/*for _, srvID := range maxloadServer {

		//暂时随机挑一个，后面再优化
		if len(minloadServer) > 0 {
			var randresult = 0
			//log.Debug(" srvID =", srvID, " number=", number )
			randresult = rand.Intn(len(minloadServer))
			msgreq := &protoMsg.MigrateCellReq{
				SrcCellApp: srvID,
				ToCellApp:   minloadServer[randresult],
				CellID:  0,
			}
			srv.PostMsgToCell(srvID, 0, msgreq)

		}
	}*/

}

func (srv *CellAppMgrSrv) Destroy() {

}

func (srv *CellAppMgrSrv) OnServerConnect(srvID uint64, serverType uint8) {

	log.Debug("CellAppMgrSrv OnServerConnect.... serverType is ", serverType, ", srvID is ", srvID)

	//把cellapp加进来
	if serverType == common.ServerTypeCellApp {
		cellapp := &CellApp{}

		cellapp.start(srvID)
		srv.cellapps[srvID] = cellapp

		cell := GetSpacesInst().getNoOwnerCell()

		if cell != nil {
			space := GetSpacesInst().getSpace()
			if space != nil {
				space.allocNewCellNotify(cell, cell.getProtoMsgRect())
			}
		}
	}

}

func (srv *CellAppMgrSrv) GetEntities(spaceID uint64) iserver.IEntities {
	e := srv.GetEntity(spaceID)

	if e == nil {
		return nil
	}

	return e.(iserver.IEntities)
}

func (srv *CellAppMgrSrv) getSpaceID() {

}

//分配space, cell
func (srv *CellAppMgrSrv) allocSpace(srvID uint64) (*Space, *Cell) {

	//创建 space cell
	cellapp := srv.cellapps[srvID]
	if cellapp != nil {
		//先写死，0,100,0,100
		//return cellapp.spaces.newSpace(0, 100, 0, 100, srvID)

		space := GetSpacesInst().getSpace()
		if space != nil {
			return space, nil
		} else {
			return GetSpacesInst().newSpace(0, 100, 0, 100, srvID)
		}
	}

	return nil, nil
}

//分配space, cell
func (srv *CellAppMgrSrv) allocSpaceCell() (*Space, *Cell) {

	//创建 space cell
	//先写死，0,100,0,100

	space := GetSpacesInst().getSpace()
	if space != nil {
		return space, nil
	} else {
		//srvID为0,现在的cell还不属于任何srvID
		return GetSpacesInst().newSpace(0, 100, 0, 100, 0)
	}

	return nil, nil
}

//把拆分后的cell发给合适的cellapp
//todo: 后面可以封装成统一的转发，而不需要每个转发函数都要定义一个函数
func (srv *CellAppMgrSrv) broadcastCellInfo(content msgdef.IMsg) {
	msg := content.(*protoMsg.CellInfoNotify)

	//判定这个cellapp上的压力，如果压力达到交出cell的条件，就交给另外一个合适的cellapp去处理
	//app := srv.getCellApp(msg.SrvID)
	//load := app.getOverLoad()
	//if load > 90 {
	for _, cellapp := range srv.cellapps {
		//如果别的cellapp里有负载比较轻的
		//if (cellapp.getOverLoad <50){
		log.Debug("broadcastCellInfo to cellapp id = ", cellapp.getID(), " cellid = ", msg.CellID, " xmin = ", msg.RectInfo.Xmin, " xmax = ", msg.RectInfo.Xmax, " ymin = ", msg.RectInfo.Ymin, " ymax = ", msg.RectInfo.Ymax, " operate = ", msg.Operate)
		srv.PostMsgToCell(cellapp.getID(), 0, msg)

		//}
	}
	//}
}

//根据地图名和坐标拿到对应的cell
func (srv *CellAppMgrSrv) getCell(mapname string, pos *protoMsg.Vector3) (*Space, *Cell) {
	for _, app := range srv.cellapps {
		space := app.getSpaces().getSpaceByMapName(mapname)
		cell := space.getCellByPos(pos)
		if (space != nil) && (cell != nil) {
			return space, cell
		}
	}
	return nil, nil
}

func (srv *CellAppMgrSrv) getCellApp(srvID uint64) *CellApp {
	return srv.cellapps[srvID]
}

//1：优先获取还没有分配cell的空闲Cellapp
//2：如果所有的cellapp都分配cell了, 拿压力最小的Cellapp

func (srv *CellAppMgrSrv) getFreeCellApp() *CellApp {

	var app *CellApp = nil
	for _, cellapp := range srv.cellapps {
		if !cellapp.isValid() {
			app = cellapp
			return app
		} else {
			if app == nil {
				app = cellapp
			} else {
				if cellapp.getOverLoad() < app.getOverLoad() {
					app = cellapp
				}
			}
		}
	}
	return app
}
