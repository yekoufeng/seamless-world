package main

import (
	"db"
	"protoMsg"
	"zeus/iserver"

	log "github.com/cihub/seelog"
)

// AnnuonceMgr 跑马灯公告管理
type AnnuonceMgr struct {
	srv *Server
}

// NewAnnuonceMgr 获取公告管理器
func NewAnnuonceMgr(srv *Server) *AnnuonceMgr {
	mgr := &AnnuonceMgr{
		srv: srv,
	}
	mgr.init()

	return mgr
}

// init 初始化管理器
func (mgr *AnnuonceMgr) init() {
	//log.Info("AnnuonceMgr init")
}

// AddAnnuouce 添加公告
func (mgr *AnnuonceMgr) AddAnnuouce(id uint64) {

	data := db.GetAnnuonceData(id)
	if data == nil {
		log.Warn("AddAnnuouce data is nill id = ", id)
		return
	}

	//log.Info("AddAnnuouce ", data, id)

	retMsg := &protoMsg.AnnuonceInfo{
		Id:           data.ID,
		StartTime:    int64(data.StartTime),
		EndTime:      int64(data.EndTime),
		InternalTime: int64(data.InternalTime),
		Content:      data.Content,
	}

	iserver.GetSrvInst().FireEvent(iserver.RPCChannel, "AddAnnuouceData", retMsg)

	return
}

// DelAnnuoucing 删除进行中公告
func (mgr *AnnuonceMgr) DelAnnuoucing(id uint64) bool {
	//log.Info("DelAnnuoucing ", id)

	iserver.GetSrvInst().FireEvent(iserver.RPCChannel, "DelAnnuoucingData", id)

	return true
}
