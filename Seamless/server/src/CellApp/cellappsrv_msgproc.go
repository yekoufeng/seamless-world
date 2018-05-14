package main

import (
	"common"
	"protoMsg"
	"zeus/linmath"
	"zeus/msgdef"

	"zeus/iserver"

	log "github.com/cihub/seelog"
)

type CellAppSrvMsgProc struct {
	srv *CellAppSrv
}

func (msgproc *CellAppSrvMsgProc) MsgProc_CellBorderChangeNotify(content msgdef.IMsg) {
	msg := content.(*protoMsg.CellBorderChangeNotify)

	space := msgproc.srv.spaces.getSpace(msg.SpaceID)
	if space != nil {
		cell, isExist := space.cells.Load(msg.CellID)
		if isExist {
			log.Debug("MsgProc_CellBorderChangeNotify  exist  spaceID= ", msg.SpaceID, " cellid =", msg.CellID, " xmin = ", msg.Rectinfo.Xmin, " xmax = ", msg.Rectinfo.Xmax, " ymin = ", msg.Rectinfo.Ymin, " ymax = ", msg.Rectinfo.Ymax)
			cell.(*Cell).setRect(msg.Rectinfo.Xmin, msg.Rectinfo.Xmax, msg.Rectinfo.Ymin, msg.Rectinfo.Ymax)
			//把cellinfo广播给所有cellapp
			space.syncCellinfoToSpace(msgproc.srv.GetCellMgrID(), cell.(*Cell).GetID(), msg.Rectinfo, msgproc.srv.GetSrvID(), 2)
		} else {
			log.Debug("MsgProc_CellChangeNotify not exist  spaceID= ", msg.SpaceID, " cellid = ", msg.CellID)
		}

	}
}

//收到广播同步的cellinfo信息
func (msgproc *CellAppSrvMsgProc) MsgProc_CellInfoNotify(content msgdef.IMsg) {

	msg := content.(*protoMsg.CellInfoNotify)
	space := msgproc.srv.spaces.getSpace(msg.SpaceID)

	log.Debug("MsgProc_CellInfoNotify cellid =", msg.CellID, "  SrvID=", msg.SrvID, " operate=", msg.Operate, msg.RectInfo.Xmin, msg.RectInfo.Xmax, msg.RectInfo.Ymin, msg.RectInfo.Ymax)
	if space != nil {
		space.setCellInfo(msg.CellID, msg.Operate, msg.SrvID, msg.RectInfo.Xmin, msg.RectInfo.Xmax, msg.RectInfo.Ymin, msg.RectInfo.Ymax)

		// 遍历当前所有的Cell
		msgproc.srv.Cells.Range(func(_, v interface{}) bool {
			cellPtr := v.(*Cell)
			if cellPtr == nil {
				log.Debug("Error msgproc.srv.Cells is nil")
				return true
			}

			// 遍历Cell 中的 Entity 发送消息
			cellPtr.Range(func(_, v1 interface{}) bool {
				e := v1.(iserver.IEntity)
				if entityGet, ok := e.GetRealPtr().(IGetEntity); ok {
					entityGet.GetEntity().SendCellInfos()
				}

				return true
			})

			return true
		})
	}

}

//new cell notity
func (msgproc *CellAppSrvMsgProc) MsgProc_CreateCellNotify(content msgdef.IMsg) {

	msg := content.(*protoMsg.CreateCellNotify)

	//创建space
	var srect, crect linmath.Rect
	if msg.CellID != 0 {
		log.Debug("MsgProc_CreateCellNotify create space cell notity   spaceID =  ", msg.SpaceID,
			//" xmin = ", msg.Srect.Xmin, " xmax = ",msg.Srect.Xmax, " ymin = ", msg.Srect.Ymin, " ymax=", msg.Srect.Ymax,
			" cellID = ", msg.CellID, " xmin = ", msg.Crect.Xmin, " xmax=", msg.Crect.Xmax, " ymin=", msg.Crect.Ymin, " ymax=", msg.Crect.Ymax)

		//srect.Xmin = msg.Srect.Xmin
		//srect.Xmax = msg.Srect.Xmax
		//srect.Ymin = msg.Srect.Ymin
		//srect.Ymax = msg.Srect.Ymax
		crect.Xmin = msg.Crect.Xmin
		crect.Xmax = msg.Crect.Xmax
		crect.Ymin = msg.Crect.Ymin
		crect.Ymax = msg.Crect.Ymax
	}

	//判定有没有space,没有就创建
	space := msgproc.srv.spaces.getSpace(msg.SpaceID)
	if space != nil {
		cell := space.getCell(msg.CellID)
		if cell == nil {
			cellinfo := &common.CellInfo{}
			cellinfo.SetCellSrvID(space.cellSrv.GetSrvID())
			cellinfo.Init(msg.CellID, msg.Crect.Xmin, msg.Crect.Xmax, msg.Crect.Ymin, msg.Crect.Ymax)
			space.newCell(msg.CellID, cellinfo)
		}
	} else {
		msgproc.srv.spaces.createNewSpace(msg.SpaceID, srect, msg.CellID, crect, msg.GetMapName())
	}

}

//delete cell
func (msgproc *CellAppSrvMsgProc) MsgProc_DeleteCellNotify(content msgdef.IMsg) {

	msg := content.(*protoMsg.DeleteCellNotify)

	log.Debug("MsgProc_DeleteCellNotify delete space cell notity   spaceID =  ", msg.SpaceID, " cellID=", msg.CellID)

	space := msgproc.srv.spaces.getSpace(msg.SpaceID)
	if space != nil {
		cell := space.getCell(msg.CellID)
		if cell != nil {
			space.delCell(msg.CellID)
		}
	}
}
