package main

import (
	"sync"

	"zeus/linmath"

	log "github.com/cihub/seelog"
)

type Spaces struct {
	spaces  *sync.Map
	cellSrv *CellAppSrv
}

func (sps *Spaces) Init(cellSrv *CellAppSrv) {
	sps.cellSrv = cellSrv
	sps.spaces = &sync.Map{}
}

//创建一个space
func (sps *Spaces) createNewSpace(space_id uint64, srect linmath.Rect, cell_id uint64, crect linmath.Rect, mapName string) {

	space := Space{
		spaceID: space_id,
		flag:    false,
		mapName: mapName,
		cellSrv: sps.cellSrv,
	}

	space.Init(cell_id, srect, crect, mapName)
	space.flag = true

	sps.spaces.Store(space_id, &space)

	
	log.Info("spaces crreateNewSapce.... spaceid = ",
		space_id, " cell_id = ", cell_id, " flag = ", space.flag, ", mapName: ", mapName)

}

func (sps *Spaces) self() *sync.Map {
	return sps.spaces
}

func (sps *Spaces) destorySpace() {

}

func (sps *Spaces) getSpace(cellID uint64) *Space {
	space, ok := sps.spaces.Load(cellID)
	if ok {
		return space.(*Space)
	}
	return nil
}

func (sps *Spaces) update() {

	//todo: 遍历所有的space
}
