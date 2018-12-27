package main

import (
	"common"
	"sync"
	"time"
	"zeus/iserver"
	"zeus/linmath"

	"protoMsg"

	log "github.com/cihub/seelog"
)

//区别于zeus下的空间实体
type Space struct {
	//SpaceID
	spaceID uint64
	cellSrv *CellAppSrv

	//本进程所负责space上的cell
	cells *sync.Map

	//全局space上的cellinfo
	//map <cellid, * common.CellInfo>
	cellinfos *sync.Map

	//space是否创建完标记
	flag bool

	mapName string
	rect    linmath.Rect
}

func (sp *Space) Init(cellid uint64, srect linmath.Rect, crect linmath.Rect, mapName string) {

	sp.cells = &sync.Map{}
	sp.cellinfos = &sync.Map{}
	sp.rect = srect
	sp.mapName = mapName
	// 第一个space 初始为一个cell
	if cellid != 0 {

		cellinfo := &common.CellInfo{}
		cellinfo.SetCellSrvID(sp.cellSrv.GetSrvID())
		cellinfo.Init(cellid, crect.Xmin, crect.Xmax, crect.Ymin, crect.Ymax)
		sp.newCell(cellid, cellinfo)
	}
}

func (sp *Space) newCell(cellID uint64, cellInfo *common.CellInfo) *Cell {
	cell := &Cell{
		cellID:   cellID,
		cellInfo: cellInfo,
		space:    sp,
	}

	cell.Init()
	cell.space = sp
	sp.cells.Store(cellID, cell)
	sp.cellSrv.Cells.Store(cellID, cell)

	//把cellinfo广播给所有cellapp
	sp.syncCellinfoToSpace(sp.cellSrv.GetCellMgrID(), cell.GetID(), cell.getProtoMsgRect(), GetSrvInst().GetSrvID(), 1)

	log.Info("newCell， cellID: ", cellID)

	deltaTime := time.Millisecond * time.Duration(1000/gFPS)
	go func() {
		ticker := time.NewTicker(deltaTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if cell.IsDestroyed() {
					return
				}
				cell.MainLoop()
			}
		}
	}()
	return cell
}

func (sp *Space) delCell(cellid uint64) {

	sp.cells.Delete(cellid)

	//把cellinfo广播给所有cellapp
	sp.syncCellinfoToSpace(sp.cellSrv.GetCellMgrID(), cellid, &protoMsg.RectInfo{}, GetSrvInst().GetSrvID(), 3)

}

func (sp *Space) getID() uint64 {
	return sp.spaceID
}

func (sp *Space) GetMapName() string {
	return sp.mapName
}

func (sp *Space) SetMapName(mapName string) {
	sp.mapName = mapName
}

// GetTimeStamp 获取当前的时间戳
func (s *Space) GetTimeStamp() uint32 {
	//return uint32(time.Now().Sub(s.startTime) / iserver.GetSrvInst().GetFrameDeltaTime())
	return 0
}

func (sp *Space) getCell(cellid uint64) *Cell {
	cell, ok := sp.cells.Load(cellid)
	if ok {
		return cell.(*Cell)
	} else {
		return nil
	}
}

func (sp *Space) isCellExist(cellid uint64) bool {
	_, ok := sp.cells.Load(cellid)
	return ok
}

func (sp *Space) update() {

	//todo: 判断所管辖cell上实体的数量，如果实体数量为0,返回false,让Spaces去删除

	//todo: 清理ghost等

}

//同步cellinfo到全局Space, 通过cellappmgr转发
func (sp *Space) syncCellinfoToSpace(cellappmgrid uint64, cellID uint64, rectinfo *protoMsg.RectInfo, srvID uint64, op uint32) {

	msg := &protoMsg.CellInfoNotify{
		Operate:  op,
		SpaceID:  sp.getID(),
		CellID:   cellID,
		RectInfo: rectinfo,
		SrvID:    srvID,
	}

	iserver.GetSrvInst().PostMsgToCell(cellappmgrid, 0, msg)
}

func (sp *Space) setCellInfo(cellid uint64, operate uint32, srvID uint64, xmin float64, xmax float64, ymin float64, ymax float64) {

	var cellinfo *common.CellInfo
	cellinfoI, ok := sp.cellinfos.Load(cellid)

	if !ok {
		cellinfo = &common.CellInfo{}
		cellinfo.SetCellSrvID(srvID)
		cellinfo.Init(cellid, xmin, xmax, ymin, ymax)
		sp.cellinfos.Store(cellid, cellinfo)
	} else {
		if operate == 2 {
			cellinfo = cellinfoI.(*common.CellInfo)
			cellinfo.Init(cellid, xmin, xmax, ymin, ymax)
		} else if operate == 3 {
			sp.cellinfos.Delete(cellid)
		}
	}

}

func (sp *Space) getCellInfos() *sync.Map {

	return sp.cellinfos
}
