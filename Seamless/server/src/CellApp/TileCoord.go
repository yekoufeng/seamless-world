package main

import (
	"container/list"
	"errors"
	"zeus/iserver"
	"zeus/linmath"

	"github.com/cihub/seelog"
)

//目前先写死了，支持两种类型的可视距离
const (
	farDist  = 300
	nearDist = 40
)

// TileCoord 基于九宫格的坐标系统
type TileCoord struct {
	farTiles  *Tiles
	nearTiles *Tiles
}

// NewTileCoord  创建新的坐标系统
func NewTileCoord(width, height int) *TileCoord {
	return &TileCoord{
		farTiles:  newTiles(width, height, false),
		nearTiles: newTiles(width, height, true),
	}
}

//UpdateCoord 更新坐标位置
func (c *TileCoord) UpdateCoord(n iserver.ICoordEntity) {

	if n.IsAOITrigger() || !n.IsNearAOILayer() {
		c.farTiles.update(n)
	}

	// if n.IsAOITrigger() || n.IsNearAOILayer() {
	// 	c.nearTiles.update(n)
	// }
}

//RemoveFromCoord 从坐标系统中删除
func (c *TileCoord) RemoveFromCoord(n iserver.ICoordEntity) {

	seelog.Debug("RemoveFromCoord, ID: ", n.GetID())

	if n.IsAOITrigger() || !n.IsNearAOILayer() {
		c.farTiles.remove(n)
	}

	// if n.IsAOITrigger() || n.IsNearAOILayer() {
	// 	c.nearTiles.remove(n)
	// }
}

// TravsalAOI 遍AOI范围内的对象
func (c *TileCoord) TravsalAOI(n iserver.ICoordEntity, cb func(iserver.ICoordEntity)) {

	if n.IsAOITrigger() || !n.IsNearAOILayer() {
		c.farTiles.TravsalAOI(n, cb)
	}

	// if n.IsAOITrigger() || n.IsNearAOILayer() {
	// 	c.nearTiles.TravsalAOI(n, cb)
	// }
}

// 遍历center为中心，半径范围radius内的所有实体执行cb
func (c *TileCoord) TravsalRange(center *linmath.Vector3, radius int, cb func(iserver.ICoordEntity)) {
	c.farTiles.TravsalRange(center, radius, cb)
}

// 遍历center所在的Tower，在该Tower内的center为中心，半径范围radius内的所有实体执行cb
func (c *TileCoord) TravsalCenter(center *linmath.Vector3, radius int, cb func(iserver.ICoordEntity)) {
	c.farTiles.TravsalCenter(center, radius, cb)
}

////////////////////////////////////////////////////////////////

// CoordPos Coord系统使用的位置
type CoordPos struct {
	X int
	Z int
}

func newCoordPos(pos linmath.Vector3) CoordPos {
	return CoordPos{
		X: int(pos.X),
		Z: int(pos.Z),
	}
}

//////////////////////////////////////////////////////////////

const (
	towerDir_All       = 0
	towerDir_Left      = 1
	towerDir_LeftDown  = 2
	towerDir_Down      = 3
	towerDir_RightDown = 4
	towerDir_Right     = 5
	towerDir_RightUp   = 6
	towerDir_Up        = 7
	towerDir_LeftUp    = 8
)

// Tower 灯塔
type Tower struct {
	tiles *Tiles

	gridX int
	gridZ int

	neighbours []int

	aoiEntities *list.List //IsAOITrigger true
	entities    *list.List //IsAOITrigger false
}

func newTower(tiles *Tiles, gridX, gridZ int) *Tower {

	t := new(Tower)
	t.tiles = tiles
	t.gridX = gridX
	t.gridZ = gridZ
	t.entities = list.New()
	t.aoiEntities = list.New()

	t.init()

	return t
}

func (t *Tower) init() {
	t.neighbours = make([]int, 9)

	t.neighbours[towerDir_All] = t.getNeighbourID(0, 0)
	t.neighbours[towerDir_Left] = t.getNeighbourID(-1, 0)
	t.neighbours[towerDir_LeftDown] = t.getNeighbourID(-1, -1)
	t.neighbours[towerDir_Down] = t.getNeighbourID(0, -1)
	t.neighbours[towerDir_RightDown] = t.getNeighbourID(1, -1)
	t.neighbours[towerDir_Right] = t.getNeighbourID(1, 0)
	t.neighbours[towerDir_RightUp] = t.getNeighbourID(1, 1)
	t.neighbours[towerDir_Up] = t.getNeighbourID(0, 1)
	t.neighbours[towerDir_LeftUp] = t.getNeighbourID(-1, 1)
}

func (t *Tower) getNeighbourID(deltaX, deltaZ int) int {
	return t.tiles.getTowerID(t.gridX+deltaX, t.gridZ+deltaZ)
}

func (t *Tower) add(n iserver.ICoordEntity) {
	t.notifyTowerAdd(n)
	t.travsalNeighour(towerDir_All, func(tt *Tower) { tt.notifyTowerAdd(n) })

	t.addToList(n)
}

func (t *Tower) addToList(n iserver.ICoordEntity) {
	if n.IsAOITrigger() {
		t.aoiEntities.PushFront(n)
	} else {
		t.entities.PushFront(n)
	}
}

func (t *Tower) remove(n iserver.ICoordEntity) {
	t.removeFromList(n)

	t.notifyTowerRemove(n)
	t.travsalNeighour(towerDir_All, func(tt *Tower) { tt.notifyTowerRemove(n) })
}

func (t *Tower) removeFromList(n iserver.ICoordEntity) {

	if n.IsAOITrigger() {
		for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
			if e.Value == n {
				t.aoiEntities.Remove(e)
				break
			}
		}
	} else {
		for e := t.entities.Front(); e != nil; e = e.Next() {
			if e.Value == n {
				t.entities.Remove(e)
				break
			}
		}
	}
}

func (t *Tower) moveTo(n iserver.ICoordEntity, nt *Tower) {

	deltaX := nt.gridX - t.gridX
	deltaZ := nt.gridZ - t.gridZ

	if deltaX == 0 && deltaZ == 0 {
		return
	}

	if deltaX > 1 || deltaX < -1 || deltaZ > 1 || deltaZ < -1 {
		t.remove(n)
		nt.add(n)
		return
	}

	t.removeFromList(n)
	t.travsalNeighour(t.getInvertDir(t.getDir(deltaX, deltaZ)), func(tt *Tower) { tt.notifyTowerRemove(n) })

	nt.travsalNeighour(t.getDir(deltaX, deltaZ), func(tt *Tower) { tt.notifyTowerAdd(n) })
	nt.addToList(n)
}

func (t *Tower) getDir(deltaX, deltaZ int) int {

	var dir int

	if deltaX == 1 && deltaZ == 0 {
		dir = towerDir_Right
	} else if deltaX == 1 && deltaZ == -1 {
		dir = towerDir_RightDown
	} else if deltaX == 0 && deltaZ == -1 {
		dir = towerDir_Down
	} else if deltaX == -1 && deltaZ == -1 {
		dir = towerDir_LeftDown
	} else if deltaX == -1 && deltaZ == 0 {
		dir = towerDir_Left
	} else if deltaX == -1 && deltaZ == 1 {
		dir = towerDir_LeftUp
	} else if deltaX == 0 && deltaZ == 1 {
		dir = towerDir_Up
	} else if deltaX == 1 && deltaZ == 1 {
		dir = towerDir_RightUp
	}

	return dir
}

func (t *Tower) getInvertDir(dir int) int {

	var invDir int

	switch dir {
	case towerDir_Left:
		invDir = towerDir_Right
	case towerDir_LeftDown:
		invDir = towerDir_RightUp
	case towerDir_Down:
		invDir = towerDir_Up
	case towerDir_RightDown:
		invDir = towerDir_LeftUp
	case towerDir_Right:
		invDir = towerDir_Left
	case towerDir_RightUp:
		invDir = towerDir_LeftDown
	case towerDir_Up:
		invDir = towerDir_Down
	case towerDir_LeftUp:
		invDir = towerDir_RightDown
	}

	return invDir
}

func (t *Tower) travsalNeighour(dir int, cb func(*Tower)) {
	t.tiles.travsalNeighour(t, dir, cb)
}

func (t *Tower) notifyTowerAdd(n iserver.ICoordEntity) {
	// if n.IsGhost() {
	// 	return
	// }

	if !n.IsAOITrigger() {

		for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
			ii := e.Value.(iserver.IAOITriggerEntity)
			ii.OnEntityEnterAOI(n)
		}

	} else {

		in := n.(iserver.IAOITriggerEntity)

		if n.IsNearAOILayer() == t.tiles.isNearLayer {

			for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
				ii := e.Value.(iserver.IAOITriggerEntity)
				ii.OnEntityEnterAOI(n)
			}
		}

		for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
			ii := e.Value.(iserver.ICoordEntity)
			if ii.IsNearAOILayer() == t.tiles.isNearLayer {
				in.OnEntityEnterAOI(ii)
			}
		}

		for e := t.entities.Front(); e != nil; e = e.Next() {
			ii := e.Value.(iserver.ICoordEntity)
			in.OnEntityEnterAOI(ii)
		}
	}
}

func (t *Tower) notifyTowerRemove(n iserver.ICoordEntity) {
	if n.IsGhost() {
		return
	}

	if !n.IsAOITrigger() {

		for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
			ii := e.Value.(iserver.IAOITriggerEntity)
			ii.OnEntityLeaveAOI(n)
		}

	} else {

		in := n.(iserver.IAOITriggerEntity)

		if n.IsNearAOILayer() == t.tiles.isNearLayer {
			for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
				ii := e.Value.(iserver.IAOITriggerEntity)
				ii.OnEntityLeaveAOI(n)
			}
		}

		for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
			ii := e.Value.(iserver.ICoordEntity)
			if ii.IsNearAOILayer() == t.tiles.isNearLayer {
				in.OnEntityLeaveAOI(ii)
			}
		}

		for e := t.entities.Front(); e != nil; e = e.Next() {
			ii := e.Value.(iserver.ICoordEntity)
			in.OnEntityLeaveAOI(ii)
		}
	}
}

// TravsalAOI 遍历AOI对象
func (t *Tower) TravsalAOI(cb func(iserver.ICoordEntity)) {

	for e := t.aoiEntities.Front(); e != nil; e = e.Next() {
		ii := e.Value.(iserver.ICoordEntity)

		if ii.IsNearAOILayer() == t.tiles.isNearLayer && ii.IsAOITrigger() {
			cb(ii)
		}
	}
}

////////////////////////////////////////////////////////////////

// CoordInfo 坐标结节
type CoordInfo struct {
	entity iserver.ICoordEntity
	tower  *Tower
}

////////////////////////////////////////////////////////////////

// Tiles 九宫格系统
type Tiles struct {
	startX        int
	startZ        int
	width         int
	height        int
	gridSize      int
	gridWidthNum  int
	gridHeightNum int
	dist          int
	isNearLayer   bool

	towers  []*Tower
	entites map[uint64]*CoordInfo

	towerNeighbours [][]int
}

// dist是可视距离
func newTiles(width, height int, isNearLayer bool) *Tiles {

	ts := new(Tiles)
	ts.width = width
	ts.height = height
	ts.isNearLayer = isNearLayer

	if isNearLayer {
		ts.dist = nearDist
	} else {
		ts.dist = farDist
	}

	ts.init()
	return ts
}

func (t *Tiles) init() {
	t.gridSize = t.dist
	t.gridWidthNum = t.width/t.gridSize + 1
	t.gridHeightNum = t.height/t.gridSize + 1
	t.towers = make([]*Tower, t.gridWidthNum*t.gridHeightNum)

	t.entites = make(map[uint64]*CoordInfo)

	t.initNeighbour()
}

func (t *Tiles) initNeighbour() {

	t.towerNeighbours = make([][]int, 9)
	t.towerNeighbours[towerDir_All] = []int{towerDir_Left, towerDir_LeftDown, towerDir_Down, towerDir_RightDown, towerDir_Right, towerDir_RightUp, towerDir_Up, towerDir_LeftUp}
	t.towerNeighbours[towerDir_Left] = []int{towerDir_LeftUp, towerDir_Left, towerDir_LeftDown}
	t.towerNeighbours[towerDir_LeftDown] = []int{towerDir_LeftUp, towerDir_Left, towerDir_LeftDown, towerDir_Down, towerDir_RightDown}
	t.towerNeighbours[towerDir_Down] = []int{towerDir_LeftDown, towerDir_Down, towerDir_RightDown}
	t.towerNeighbours[towerDir_RightDown] = []int{towerDir_LeftDown, towerDir_Down, towerDir_RightDown, towerDir_Right, towerDir_RightUp}
	t.towerNeighbours[towerDir_Right] = []int{towerDir_RightDown, towerDir_Right, towerDir_RightUp}
	t.towerNeighbours[towerDir_RightUp] = []int{towerDir_RightDown, towerDir_Right, towerDir_RightUp, towerDir_Up, towerDir_LeftUp}
	t.towerNeighbours[towerDir_Up] = []int{towerDir_RightUp, towerDir_Up, towerDir_LeftUp}
	t.towerNeighbours[towerDir_LeftUp] = []int{towerDir_RightUp, towerDir_Up, towerDir_LeftUp, towerDir_Left, towerDir_LeftDown}

}

// Update 更新位置
func (t *Tiles) update(n iserver.ICoordEntity) {
	pos := newCoordPos(n.GetPos())

	if !t.isValidPos(pos) {
		t.remove(n)
		return
	}

	info, ok := t.entites[n.GetID()]
	if !ok {
		t.add(n, pos)
	} else {
		t.move(info, pos)
	}
}

// Remove 从坐标系统中删除
func (t *Tiles) remove(n iserver.ICoordEntity) {

	info, ok := t.entites[n.GetID()]
	if !ok {
		return
	}

	seelog.Debug("Tiles remove, num: ", len(t.entites)-1, ", entityID: ", n.GetID())

	info.tower.remove(n)
	delete(t.entites, n.GetID())

}

func (t *Tiles) add(n iserver.ICoordEntity, pos CoordPos) {

	tower := t.getTower(pos)
	if tower == nil {
		seelog.Error("add failed ", pos, n.GetID())
		return
		// panic("inner wrong ")
	}

	t.entites[n.GetID()] = &CoordInfo{n, tower}
	tower.add(n)

	seelog.Debug("Tiles add, num: ", len(t.entites), ", entityID: ", n.GetID())
}

func (t *Tiles) move(info *CoordInfo, pos CoordPos) {

	nt := t.getTower(pos)
	if info.tower == nt || nt == nil {
		return
	}

	info.tower.moveTo(info.entity, nt)
	info.tower = nt
}

func (t *Tiles) getTower(pos CoordPos) *Tower {

	gridX := pos.X / t.gridSize
	gridZ := pos.Z / t.gridSize

	id := t.getTowerID(gridX, gridZ)
	tower, err := t.getTowerByID(id)
	if err != nil {
		seelog.Error("invalid pos ", pos)
		return nil
		// panic("inner wrong")
	}

	if tower == nil {
		tower = newTower(t, gridX, gridZ)
		t.towers[id] = tower
	}

	return tower
}

func (t *Tiles) getTowerID(gridX, gridZ int) int {

	if gridX < 0 || gridX > t.gridWidthNum || gridZ < 0 || gridZ > t.gridHeightNum {
		return -1
	}

	return gridZ*t.gridWidthNum + gridX
}

var errInvalidIndex = errors.New("invalid index")

func (t *Tiles) getTowerByID(id int) (*Tower, error) {
	if id >= len(t.towers) || id < 0 {
		return nil, errInvalidIndex
	}
	return t.towers[id], nil
}

func (t *Tiles) getTowerByPos(gridX, gridZ int) (*Tower, error) {
	id := gridZ*t.gridWidthNum + gridX
	return t.getTowerByID(id)
}

func (t *Tiles) isValidPos(pos CoordPos) bool {
	if pos.X < 0 || pos.X > t.width {
		return false
	}

	if pos.Z < 0 || pos.Z > t.height {
		return false
	}

	return true
}

func (t *Tiles) travsalNeighour(tc *Tower, dir int, cb func(*Tower)) {
	for _, v := range t.towerNeighbours[dir] {
		tn, err := t.getTowerByID(tc.neighbours[v])
		if err == nil && tn != nil {
			cb(tn)
		}
	}
}

// TravsalAOI 遍历AOI范围内的entity
func (t *Tiles) TravsalAOI(n iserver.ICoordEntity, cb func(iserver.ICoordEntity)) {
	pos := newCoordPos(n.GetPos())
	tt := t.getTower(pos)
	if tt == nil {
		return
	}

	tt.TravsalAOI(cb)
	t.travsalNeighour(tt, towerDir_All, func(nt *Tower) { nt.TravsalAOI(cb) })
}

// 遍历center为中心，半径范围radius内的所有实体执行cb
func (t *Tiles) TravsalRange(center *linmath.Vector3, radius int, cb func(iserver.ICoordEntity)) {
	minX := int(center.X) - radius
	minZ := int(center.Z) - radius
	maxX := int(center.X) + radius
	maxZ := int(center.Z) + radius

	if minX < 0 {
		minX = 0
	}
	if minZ < 0 {
		minZ = 0
	}
	if maxX > t.width {
		maxX = t.width
	}
	if maxZ > t.height {
		maxZ = t.height
	}

	minGridX := minX / t.gridSize
	mingridZ := minZ / t.gridSize
	maxGridX := maxX / t.gridSize
	maxgridZ := maxZ / t.gridSize

	// 获取center为中心，半径范围radius内的所有格子
	towers := make([]*Tower, 0)
	for gridX := minGridX; gridX <= maxGridX; gridX++ {
		for gridZ := mingridZ; gridZ <= maxgridZ; gridZ++ {
			tt, err := t.getTowerByPos(gridX, gridZ)
			if err == nil && tt != nil {
				towers = append(towers, tt)
			}
		}
	}

	// 遍历
	for _, tt := range towers {
		tt.TravsalAOI(func(ii iserver.ICoordEntity) {
			pos := ii.GetPos()
			if int(pos.X) >= minX && int(pos.X) <= maxX && int(pos.Z) >= minZ && int(pos.Z) <= maxZ {
				cb(ii)
			}
		})
	}
}

// 遍历center所在的Tower，在该Tower内的center为中心，半径范围radius内的所有实体执行cb
func (t *Tiles) TravsalCenter(center *linmath.Vector3, radius int, cb func(iserver.ICoordEntity)) {
	x := int(center.X)
	z := int(center.Z)
	gridX := x / t.gridSize
	gridZ := z / t.gridSize

	// 获取center为中心，半径范围radius内的所有格子
	tt, err := t.getTowerByPos(gridX, gridZ)
	if err != nil || tt == nil {
		return
	}

	minX := int(center.X) - radius
	minZ := int(center.Z) - radius
	maxX := int(center.X) + radius
	maxZ := int(center.Z) + radius

	if minX < 0 {
		minX = 0
	}
	if minZ < 0 {
		minZ = 0
	}
	if maxX > t.width {
		maxX = t.width
	}
	if maxZ > t.height {
		maxZ = t.height
	}

	// 遍历
	tt.TravsalAOI(func(ii iserver.ICoordEntity) {
		pos := ii.GetPos()
		if int(pos.X) >= minX && int(pos.X) <= maxX && int(pos.Z) >= minZ && int(pos.Z) <= maxZ {
			cb(ii)
		}
	})
}
