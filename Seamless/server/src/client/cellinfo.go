package main

import (
	"sync"

	"client/sdl2"

	"fmt"

	log "github.com/cihub/seelog"
	"github.com/veandco/go-sdl2/sdl"
)

type CellInfo struct {
	SpaceID uint64
	CellID  uint64
	SrvID   uint64
	MinX    float64
	MaxX    float64
	MinY    float64
	MaxY    float64
}

func GetNewCellInfo(spaceId, cellid, srvId uint64, minX, maxX, minY, maxY float64) *CellInfo {
	obj := &CellInfo{
		SpaceID: spaceId,
		CellID:  cellid,
		SrvID:   srvId,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
	}

	return obj
}

type CellInfoMgr struct {
	cells []*CellInfo
}

var CIF *CellInfoMgr
var lock *sync.Mutex = &sync.Mutex{} // 同步锁

func GetCellInfoMgr() *CellInfoMgr {
	if CIF != nil {
		return CIF
	}

	CIF = &CellInfoMgr{
		cells: make([]*CellInfo, 0),
	}
	return CIF
}

func (cif *CellInfoMgr) addOne(obj *CellInfo) {
	for _, cell := range cif.cells {
		if cell.CellID == obj.CellID {
			cell.SrvID = obj.SrvID
			cell.SpaceID = obj.SpaceID
			cell.MinX = obj.MinX
			cell.MaxX = obj.MaxX
			cell.MinY = obj.MinY
			cell.MaxY = obj.MaxY

			log.Debugf("当前有重复消息 覆盖CellID:%d", cell.CellID)
			return
		}
	}

	cif.cells = append(cif.cells, obj)
}

func (cif *CellInfoMgr) AddCell(spaceId, cellId, srvId uint64, minX, maxX, minY, maxY float64) {
	lock.Lock()
	defer lock.Unlock()

	obj := GetNewCellInfo(spaceId, cellId, srvId, minX, maxX, minY, maxY)
	cif.addOne(obj)
}

//func collision(rect *sdl.Rect, calcRect *sdl.Rect, id uint64) bool {
//	/*
//		https://www.cnblogs.com/klobohyz/archive/2012/06/25/2562089.html
//	*/
//	x1 := rect.X
//	y1 := rect.Y
//	x2 := rect.X + rect.W
//	y2 := rect.Y + rect.H
//
//	x3 := calcRect.X
//	y3 := calcRect.Y
//	x4 := calcRect.X + calcRect.W
//	y4 := calcRect.Y + calcRect.H
//
//	//log.Debug("CellID:", id)
//	result := ((x1 >= x3 && x1 <= x4) || (x3 >= x1 && x3 <= x2)) && ((y1 >= y3 && y1 < y4) || (y3 >= y1 && y3 <= y2))
//	//log.Debug("Src", rect)
//	//log.Debug("Calc", calcRect)
//	//
//	//log.Debug("y1:", y1, "y2", y2, "y3", y3, "y4", y4)
//	//result := ((x1 >= x3 && x1 <= x4) || (x3 >= x1 && x3 <= x2)) && ((y1 >= y3 && y1 <= y4) || (y3 >= y1 && y3 <= y2))
//	//log.Debug("Result:", result)
//	//log.Flush()
//	return bool(result)
//}

func (cif *CellInfoMgr) GetDrawRect(minX, minY, maxX, maxY int32) {
	lock.Lock()
	defer lock.Unlock()

	winRect := &sdl.Rect{X: minX, Y: minY, W: MainViewWidth, H: MainViewHeight}
	var drawRect *sdl.Rect

	for _, cellInfo := range cif.cells {
		drawRect = &sdl.Rect{
			X: int32(cellInfo.MinX*10.0) / 10,
			Y: int32(cellInfo.MinY*10.0) / 10,
			W: int32(cellInfo.MaxX*10.0)/10 - int32(cellInfo.MinX*10.0)/10,
			H: int32(cellInfo.MaxY*10.0)/10 - int32(cellInfo.MinY*10.0)/10,
		}

		if ok := drawRect.HasIntersection(winRect); ok == false {
			//log.Debug("Current not found cell id", cellInfo.CellID)
			continue
		}

		drawRect.X -= winRect.X
		drawRect.Y -= winRect.Y
		sdl2.GetNewRectMgr().Add(drawRect)
	}
}

func (cif *CellInfoMgr) DrawCellID(minX, minY, maxX, maxY int32, mv *MainView) {
	lock.Lock()
	defer lock.Unlock()

	winRect := &sdl.Rect{X: minX, Y: minY, W: MainViewWidth, H: MainViewHeight}
	var drawRect *sdl.Rect

	for _, cellInfo := range cif.cells {
		drawRect = &sdl.Rect{
			X: int32(cellInfo.MinX*10.0) / 10,
			Y: int32(cellInfo.MinY*10.0) / 10,
			W: int32(cellInfo.MaxX*10.0)/10 - int32(cellInfo.MinX*10.0)/10,
			H: int32(cellInfo.MaxY*10.0)/10 - int32(cellInfo.MinY*10.0)/10,
		}

		if ok := drawRect.HasIntersection(winRect); ok == false {
			//log.Debug("Current not found cell id", cellInfo.CellID)
			continue
		}

		drawRect.X -= winRect.X
		drawRect.Y -= winRect.Y

		mv.DrawInputText(drawRect.X+10, drawRect.Y+10, fmt.Sprintf("SvrID:%d", cellInfo.SrvID))
		mv.DrawInputText(drawRect.X+10, drawRect.Y+drawRect.H-30, fmt.Sprintf("SvrID:%d", cellInfo.SrvID))
		mv.DrawInputText(drawRect.X+drawRect.W-200, drawRect.Y+10, fmt.Sprintf("SvrID:%d", cellInfo.SrvID))
		mv.DrawInputText(drawRect.X+drawRect.W-200, drawRect.Y+drawRect.H-30, fmt.Sprintf("SvrID:%d", cellInfo.SrvID))
	}
}

func (cif *CellInfoMgr) GetCellInfoByPos(x, y int32) *CellInfo {
	lock.Lock()
	defer lock.Unlock()

	calcPoint := &sdl.Point{X: x, Y: y}

	for _, cell := range cif.cells {
		calcRect := &sdl.Rect{X: int32(cell.MinX), Y: int32(cell.MinY),
			W: int32(cell.MaxX - cell.MinX), H: int32(cell.MaxY - cell.MinY)}
		if calcPoint.InRect(calcRect) == true {
			return cell
		}
	}
	return nil
}

func (cif *CellInfoMgr) GetRectByID(cellid uint64) *sdl.Rect {
	lock.Lock()
	defer lock.Unlock()

	var ret *sdl.Rect
	for _, cell := range cif.cells {
		if cell.CellID == cellid {
			ret = &sdl.Rect{
				X: int32(cell.MinX*10) / 10,
				Y: int32(cell.MinY*10) / 10,
				H: int32(cell.MaxY*10)/10 - int32(cell.MinY*10)/10,
				W: int32(cell.MaxX*10)/10 - int32(cell.MinX*10)/10,
			}
		}
	}
	return ret
}

func (cif *CellInfoMgr) GetAllCellInfo() []*CellInfo {
	lock.Lock()
	defer lock.Unlock()

	var lst []*CellInfo
	for _, v := range cif.cells {
		c := *v
		lst = append(lst, &c)
	}
	return lst
}
