package main

import (
	"common"
	_ "net/http/pprof"
	"zeus/iserver"
	"zeus/server"
	"zeus/tlog"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	//	"zeus/msgdef"

	"protoMsg"
	"sync"
	"time"
)

var srvInst *CellAppSrv
var gFPS int

//cellSrv     场景服务器
type CellAppSrv struct {
	iserver.IServer
	spaces       Spaces
	Cells        *sync.Map
	ReportTicker *time.Ticker
	//cellappmgr id, 后续多个cellappmgr时需要扩展
	cellappmgr_id uint64
}

// GetSrvInst 获取服务器全局实例
func GetSrvInst() *CellAppSrv {
	if srvInst == nil {
		common.InitMsg()

		srvInst = &CellAppSrv{}

		srvID := uint64(viper.GetInt("CellApp.FlagId"))
		pmin := viper.GetInt("CellApp.PortMin")
		pmax := viper.GetInt("CellApp.PortMax")

		innerPort := server.GetValidSrvPort(pmin, pmax)
		innerAddr := viper.GetString("CellApp.InnerAddr")

		opmin := viper.GetInt("CellApp.OuterPortMin")
		opmax := viper.GetInt("CellApp.OuterPortMax")

		fps := viper.GetInt("CellApp.FPS")
		gFPS = fps

		/*
			配置中, OuterAddr+OuterPort写入到redis中, 作为服务器的对外地址
			OuterListen+随机出来的端口作为实际对外监听端口
			当配置中OuterPort为0时, OuterPort就是随机端口
		*/
		//outerListen := viper.GetString("Room.OuterListen")
		listenPort := server.GetValidSrvPort(opmin, opmax)
		outerAddr := viper.GetString("CellApp.OuterAddr")
		outerPort := listenPort
		if viper.GetString("CellApp.OuterPort") != "0" {
			outerPort = viper.GetString("CellApp.OuterPort")
			listenPort = outerPort
		}

		srvInst.IServer = server.NewServer(common.ServerTypeCellApp, srvID, innerAddr+":"+innerPort, outerAddr+":"+outerPort, fps, srvInst)

		//protocal := viper.GetString("CellApp.SpaceProtocal")
		//maxConns := viper.GetInt("CellApp.MaxConns")

		tlogAddr := viper.GetString("Config.TLogAddr")
		if tlogAddr != "" {
			if err := tlog.ConfigRemoteAddr(tlogAddr); err != nil {
				log.Error(err)
			}
		}

		log.Info("CellApp Init")
		log.Info("ServerID:", srvID)
		log.Info("InnerAddr:", innerAddr+":"+innerPort)
	}

	return srvInst
}

// Init 初始化
func (srv *CellAppSrv) Init() error {
	srv.RegProtoType("Player", &CellUser{}, false)

	//向cellappmgrSrv发送创建 newspace请求消息

	srv.RegMsgProc(&CellAppSrvMsgProc{srv: srv})
	srv.Cells = &sync.Map{}
	srv.spaces.Init(srv)

	srv.ReportTicker = time.NewTicker(100 * time.Second)
	log.Info("CellAppSrv Start")

	return nil
}

func (srv *CellAppSrv) MainLoop() {

	//定时向cellappmgr汇报，cell负载

	//暂时每个space做一个汇报
	select {
	case <-srv.ReportTicker.C:
		//log.Debug("cell to cellappmgr report...")
		cellload := make(map[uint64]uint32)
		srv.spaces.self().Range(
			func(k, v interface{}) bool {
				s := v.(*Space)
				//log.Debug("cell to cellappmgr report key = ", k.(uint64), " flag= ", s.flag)
				if s.flag {
					s.cells.Range(
						func(cellK, cellV interface{}) bool {
							cell := cellV.(*Cell)
							cellload[cellK.(uint64)] = cell.getRealEntityNum()
							//temp for test
							//cell.getRealEntityNum()
							//cellload[cellK.(uint64)] = uint32(rand.Intn(100))
							return true
						})

					msg := &protoMsg.ReportCellLoad{
						SrvID:    srv.GetSrvID(),
						SpaceID:  k.(uint64),
						Cellload: cellload,
						SrvLoad:  uint32(srv.GetLoad()),
					}

					if err := iserver.GetSrvInst().PostMsgToCell(srv.cellappmgr_id, 0, msg); err != nil {
						log.Error(err)
					}
				}
				return true
			})

	default:
	}

}

func (srv *CellAppSrv) Destroy() {

}

func (srv *CellAppSrv) OnServerConnect(srvID uint64, serverType uint8) {

	log.Debug("CellAppSrv OnServerConnect.... serverType is ", serverType, " serverID is ", srvID)

	if serverType == common.ServerTypeCellAppMgr {

		//请求创建space
		/*msg := &protoMsg.CreateSpaceReq{
			CellappIndex: 1,
			SrvID:        srv.GetSrvID(),
		}

		log.Debug("CellAppSrv send req create space msg, srvID: ", srv.GetSrvID())
		if err := iserver.GetSrvInst().PostMsgToCell(srvID, 0, msg); err != nil {
			log.Error(err)
		}
		*/
		srv.cellappmgr_id = srvID

	}
}

func (p *CellAppSrv) GetEntities(cellID uint64) iserver.IEntities {

	cell, ok := p.Cells.Load(cellID)

	if ok {
		entity, isOK := cell.(iserver.IEntities)
		if isOK {
			//log.Info("GetEntities, 找到了 ")
			return entity
		} else {
			log.Info("GetEntities, 不能转化为iserver.IEntities ")
		}
	}

	log.Info("GetEntities，找不到 cellID: ", cellID)
	return nil
}

func (srv *CellAppSrv) GetCell(spaceID uint64, cellID uint64) *Cell {

	space := srv.spaces.getSpace(spaceID)
	if space != nil {
		return space.getCell(cellID)
	}
	return nil
}

func (srv *CellAppSrv) GetCellMgrID() uint64 {
	return srv.cellappmgr_id
}
