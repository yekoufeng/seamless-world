package main

import (
	"protoMsg"
	"zeus/iserver"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
	//	"zeus/linmath"
)

type CellAppMgrSrvMsgProc struct {
	srv *CellAppMgrSrv
}

//请求创建一个space
//暂时不用
func (mgrProc *CellAppMgrSrvMsgProc) MsgProc_CreateSpaceReq(content msgdef.IMsg) {

	msg := content.(*protoMsg.CreateSpaceReq)

	//把cellapp加进来
	cellapp := &CellApp{}

	cellapp.start(msg.SrvID)
	mgrProc.srv.cellapps[msg.SrvID] = cellapp

	space, cell := mgrProc.srv.allocSpace(msg.SrvID)

	var crect protoMsg.RectInfo
	var srect protoMsg.RectInfo

	srect.Xmin = space.getRect().Xmin
	srect.Xmax = space.getRect().Xmax
	srect.Ymin = space.getRect().Xmin
	srect.Ymax = space.getRect().Xmax

	var msgret *protoMsg.CreateSpaceRet
	if cell != nil {
		log.Debug("CellAppMgrSrvMsgProc receive create space req msg   Index =  ", msg.CellappIndex, " srvID = ", msg.SrvID, " spaceID= ", space.getID(), " cellID= ", cell.getID())
		crect.Xmin = cell.getRect().Xmin
		crect.Xmax = cell.getRect().Xmax
		crect.Ymin = cell.getRect().Xmin
		crect.Ymax = cell.getRect().Xmax

		msgret = &protoMsg.CreateSpaceRet{
			SpaceID: space.getID(),
			Srect:   &srect,
			CellID:  cell.getID(),
			Crect:   &crect,
			MapName: "1",
		}
	} else {
		msgret = &protoMsg.CreateSpaceRet{
			SpaceID: space.getID(),
			Srect:   &srect,
			CellID:  0,
			Crect:   nil,
			MapName: "1",
		}
	}

	if err := iserver.GetSrvInst().PostMsgToCell(msg.SrvID, 0, msgret); err != nil {
		log.Error(err)
	}

}

//查询所有相关进程space信息
func (mgr *CellAppMgrSrvMsgProc) RPC_querySpaces() {

}

//更新space数据
func (mgr *CellAppMgrSrvMsgProc) RPC_updateSpaceData() {

}

//汇报cell负载
func (mgr *CellAppMgrSrvMsgProc) MsgProc_ReportCellLoad(content msgdef.IMsg) {

	msg := content.(*protoMsg.ReportCellLoad)
	log.Debug("CellAppMgrSrvMsgProc report cell load   ", msg.SrvID, " spaceID = ", msg.SpaceID, " srvload = ", msg.SrvLoad)

	cellapp := mgr.srv.cellapps[msg.SrvID]
	if cellapp != nil {
		//判断有没有space,没有就创建，有就修改
		//sp := cellapp.spaces.getSpace(msg.SpaceID)
		sp := GetSpacesInst().getSpace()
		cellapp.setOverLoad(msg.SrvLoad)
		if sp == nil {
			log.Debug("new space..... SpaceID = ", msg.SpaceID)
			GetSpacesInst().newSpaceByLoad(msg.SpaceID, msg.Cellload)
		} else {
			log.Debug("update space cellload ..... ")
			sp.updateCellData(msg.Cellload)
		}

	} else {

	}
}

//MsgProc_CellInfoReq 根据地图坐标获取cell信息
func (mgr *CellAppMgrSrvMsgProc) MsgProc_CellInfoReq(content msgdef.IMsg) {
	msgRecv := content.(*protoMsg.CellInfoReq)

	space, cell := mgr.srv.getCell(msgRecv.MapName, msgRecv.Pos)

	msg := &protoMsg.CellInfoRet{
		EntityID:  msgRecv.GetEntityID(),
		SpaceID:   space.getID(),
		CellID:    cell.getID(),
		CellSrvID: cell.getSrvID(),
		Pos:       msgRecv.GetPos(),
	}

	if err := iserver.GetSrvInst().PostMsgToCell(msgRecv.GetSrvID(), 0, msg); err != nil {
		log.Error(err)
	}
}

func (mgrProc *CellAppMgrSrvMsgProc) MsgProc_CellInfoNotify(content msgdef.IMsg) {

	//log.Debug("CellAppMgrSrvMsgProc MsgProc_CellInfoNotify  ")
	mgrProc.srv.broadcastCellInfo(content)

}
