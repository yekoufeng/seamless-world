package main

//"sync"

import (
	log "github.com/cihub/seelog"
	"sync"
)

type CellApp struct {
	id     uint64

	//是否已分配cell标记
	flag       bool
	//负责哪些场景区域
	cells               *sync.Map

	//每个cellapp管理的大地图区域面积总和
	area   uint32
	//每个cellapp管理的大地图区域的负载(人数)
	load   uint32
	//cellapp的cpu或者内存负载
	overload uint32
}

func (app *CellApp) setValid(b bool) {
	app.flag = b
}

func (app *CellApp) isValid() bool {
	return app.flag
}

func (app *CellApp) start(srvID uint64) {
	app.id = srvID
	app.flag  = false
	app.cells = &sync.Map{}
	log.Debug("Cellapp srvID = ", srvID)

}

func (app *CellApp) getSpaces() *Spaces {
	return GetSpacesInst()
}

func (app *CellApp) getID() uint64 {
	return app.id
}

func (app *CellApp) setOverLoad(load uint32) {
	app.overload = load
}

func (app *CellApp) getOverLoad() uint32 {
	return app.overload
}

func (app *CellApp) addCell(cell * Cell) {
	app.cells.Store(cell.getID(), cell)
}