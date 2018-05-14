package main

import (
	"sync"
	"time"
	//"time"
)

var spacesInst *Spaces

type Spaces struct {
	spaces *sync.Map
}

//获取spaces全局实例

func GetSpacesInst() *Spaces {
	if spacesInst == nil {
		spacesInst = &Spaces{}
	}
	return spacesInst
}

func (sps *Spaces) Init() {
	sps.spaces = &sync.Map{}
}

func (sps *Spaces) Run() {
	sps.Init()
	go sps.doloop()
}

func (sps *Spaces) doloop() {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sps.spaces.Range(
				func(k, spaceV interface{}) bool {
					space := spaceV.(*Space)
					space.doloop()
					return true
				})
		}

	}

}

func (sps *Spaces) getSpaces() *sync.Map {
	return sps.spaces
}

/*func (sps *Spaces) getSpace(cellID uint64) interface{} {

	if sp, ok := sps.spaces.Load(cellID); ok {
		return sp
	}
	return nil
}*/

//获得大地图
func (sps *Spaces) getSpace() *Space {
	var s *Space = nil
	sps.spaces.Range(
		func(K, spaceV interface{}) bool {
			space := spaceV.(*Space)
			if space.stype == 1 {
				s = space
				return false
			}
			return true
		})
	return s
}

//是否存在space, 后面如果有副本的话，需要给space加一个类型，判定是大地图还是副本。
func (sps *Spaces) isExistSpace() bool {

	var isExist bool = false
	sps.spaces.Range(
		func(K, spaceV interface{}) bool {
			space := spaceV.(*Space)
			if space.stype == 1 {
				isExist = true
				return false
			}
			return true
		})
	return isExist
}

func (sps *Spaces) newSpace(xmin float64, xmax float64, ymin float64, ymax float64, srvID uint64) (*Space, *Cell) {

	cellID := GetSrvInst().FetchTempID()

	space := &Space{
		spaceID_: cellID,
		cells:    &sync.Map{},
		tree:     nil,
	}
	space.Init(xmin, xmax, ymin, ymax)
	cell := space.newCell(srvID)

	//新建一棵bsp树
	space.tree = &bspTree{
		root:   nil,
		layer:  1,
		isInit: false,
	}
	//新建第一个bspNode
	firstNode := space.tree.newBspNode(cell)
	firstNode.Parent = nil
	cell.Init(xmin, xmax, ymin, ymax, firstNode)

	space.tree.Init(xmin, xmax, ymin, ymax, space, cell)
	sps.spaces.Store(cellID, space)

	//return cellID, cell.getID()
	return space, cell
}

//这个函数不应该被调到，后面处理掉
func (sps *Spaces) newSpaceByLoad(cellID uint64, cellload map[uint64]uint32) {

	space := &Space{
		spaceID_: cellID,
	}
	//space.Init()
	space.newCellByLoad(cellload)
	sps.spaces.Store(cellID, space)

}

//根据地图名和坐标拿到对应的cell
func (sps *Spaces) getSpaceByMapName(mapname string) *Space {
	var s *Space = nil
	sps.spaces.Range(
		func(K, spaceV interface{}) bool {
			space := spaceV.(*Space)
			s = space
			return false
		})

	return s

}

//获取还没有分配给cellapp的cell
func (sps *Spaces) getNoOwnerCell() *Cell {

	sp := sps.getSpace()
	var cell *Cell = nil
	sp.cells.Range(
		func(k, cellV interface{}) bool {
			if cellV.(*Cell).getSrvID() == 0 {
				cell = cellV.(*Cell)
				return false
			} else {
				return true
			}
		})

	return cell
}
