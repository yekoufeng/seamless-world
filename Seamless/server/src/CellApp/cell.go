package main

import (
	"common"
	"errors"
	"excel"
	"time"

	"zeus/entity"
	"zeus/iserver"
	"zeus/msgdef"
	//"zeus/linmath"

	"sync"
	"zeus/linmath"

	"protoMsg"

	log "github.com/cihub/seelog"
)

type Cell struct {
	iserver.ICoord
	//cell上所有实体
	*entity.Entities

	space *Space
	//cellID
	cellID   uint64
	cellInfo *common.CellInfo
	mapData  excel.MapsData

	//cell上的ghost实体
	ghostEntityMap map[uint64]*Entity
	//cell上的real实体
	realEntityMap map[uint64]*Entity

	mapInfo *Map
	mapName string

	tinyEntities map[uint64]ITinyEntity

	startTime time.Time

	isMapLoaded bool
}

//Init 初始化
func (c *Cell) Init() {
	mapData, ok := excel.GetMaps(common.StringToUint64(c.space.GetMapName()))
	if !ok {
		return
	}
	c.mapData = mapData

	c.ghostEntityMap = make(map[uint64]*Entity)
	c.realEntityMap = make(map[uint64]*Entity)

	// 暂时先写一个最大的尺寸
	c.ICoord = NewTileCoord(9000, 9000)
	c.Entities = entity.NewEntities(false)
	c.Entities.Init()

	c.tinyEntities = make(map[uint64]ITinyEntity)
	c.startTime = time.Now()

	c.ghostEntityMap = make(map[uint64]*Entity)
	c.realEntityMap = make(map[uint64]*Entity)

	c.RegMsgProc(&CellMsgProc{cell: c})
	c.loadMap()
}

func (c *Cell) setRect(xmin float64, xmax float64, ymin float64, ymax float64) {

	c.cellInfo.Init(c.GetID(), xmin, xmax, ymin, ymax)
}

func (c *Cell) getRect() *linmath.Rect {
	return c.cellInfo.GetRect()
}

func (c *Cell) getProtoMsgRect() *protoMsg.RectInfo {

	var rectinfo protoMsg.RectInfo
	rectinfo.Xmin = c.getRect().Xmin
	rectinfo.Xmax = c.getRect().Xmax
	rectinfo.Ymin = c.getRect().Ymin
	rectinfo.Ymax = c.getRect().Ymax
	return &rectinfo

}

// 获取cell上的玩家数
func (c *Cell) getRealEntityNum() uint32 {
	if c.getRect().GetArea() > 2001 {
		return 9999
	} else {
		return uint32(len(c.realEntityMap))
	}
}

//OnMapLoadSucceed 地图加载成功, 框架层回调
func (c *Cell) OnMapLoadSucceed() {

}

// 获取space
func (c *Cell) getSpace() *Space {
	return c.space
}

func (c *Cell) GetID() uint64 {
	return c.cellID
}

func (c *Cell) GetMapName() string {
	if c.space != nil {
		return c.space.GetMapName()
	}
	return ""
}

func (c *Cell) GetTimeStamp() uint32 {
	return 0
}

//OnMapLoadFailed 地图加载失败
func (c *Cell) OnMapLoadFailed() {
	log.Error("地图加载失败 ", c)
}

func (c *Cell) IsDestroyed() bool {

	return false
}

func (c *Cell) MainLoop() {
	c.Entities.MainLoop()

	c.Entities.Range(func(k, v interface{}) bool {
		if iA, ok := v.(iAOIUpdater); ok {
			iA.updateAOI()
		}

		return true
	})

	c.Entities.Range(func(k, v interface{}) bool {
		if iL, ok := v.(iLateLooper); ok {
			iL.onLateLoop()
		}

		return true
	})

	c.CheckReals()
}

//CheckReals 检查ghost的创建与销毁
func (c *Cell) CheckReals() {
	for _, r := range c.realEntityMap {
		r.CheckReals()
	}
}

// AddEntity 在空间中添加entity
func (c *Cell) AddEntity(entityType string, entityID uint64, dbid uint64, initParam interface{}, syncInit bool, realServerID uint64, realCellID uint64) (*Entity, error) {
	log.Info("AddEntity entityType:", entityType, ", entityID:", entityID, ", dbid: ", dbid)

	e, err := c.CreateEntity(entityType, entityID, dbid, c.cellID, initParam, syncInit, realServerID)
	if err != nil {
		return nil, err
	}

	_, ok := e.(iserver.ICellEntity)
	if !ok {
		c.DestroyEntity(e.GetID())
		return nil, errors.New("the entity which add to space must be ICellEntity ")
	}

	entityGet, ok := e.GetRealPtr().(IGetEntity)
	entity := entityGet.GetEntity()
	if entity == nil {
		log.Info("AddEntity 获取实体为nil, entityType:", entityType, ", entityID:", entityID, "dbid: ", dbid)
		return nil, errors.New("the entity is nil ")
	}

	if realServerID != 0 {
		c.ghostEntityMap[entityID] = entityGet.GetEntity()
	} else {
		c.realEntityMap[entityID] = entityGet.GetEntity()
	}

	entity.SetRealCellID(realCellID)

	return entity, nil
}

// OnEntityDestory 删除entity
func (c *Cell) OnEntityDestory(entity *Entity) {
	if entity == nil {
		return
	}

	log.Info("OnEntityDestory entityType:", entity.GetType(), ", entityID:", entity.GetID(), ", dbid: ", entity.GetDBID())

	if entity.IsGhost() {
		delete(c.ghostEntityMap, entity.GetID())
	} else {
		delete(c.realEntityMap, entity.GetID())
	}
}

func (c *Cell) getRectEntities(rect *linmath.Rect) *sync.Map {

	var entities sync.Map
	c.Entities.Range(
		func(k, v interface{}) bool {
			entity := v.(Entity)
			if (float64(entity.pos.X) >= rect.Xmin) && (float64(entity.pos.X) <= rect.Xmax) && (float64(entity.pos.Y) >= rect.Ymin) && (float64(entity.pos.Y) <= rect.Ymax) {
				entities.Store(k, v)
			}
			return true
		})

	return &entities
}

func (c *Cell) GhostToReal(msg *msgdef.TransferRealReq) {
	e, ok := c.ghostEntityMap[msg.EntityID]
	if !ok {
		log.Error("GhostToReal 不在ghostMap中EntityID:", msg.EntityID)
		return
	}

	realEntity := CreateRealEntity()
	realEntity.umpackHaunt((int)(msg.GhostNum), msg.GhostData)
	e.SetRealEntity(realEntity)

	e.ReflushFromMsg((int)(msg.PropNum), msg.Props)
	e.SetRealServerID(0)
	e.SetRealCellID(0)

	c.ExchangeRealGhost(msg.EntityID)

	//向数据库中注册
	e.RegSrvID()

	//设置最新坐标
	e.SetPos(msg.Pos)

	//通知所有ghost更新realcellID
	sendMsg := &msgdef.NewRealNotify{
		RealServerID: srvInst.GetCurSrvInfo().ServerID,
		RealCellID:   e.GetCellID(),
	}

	e.SyncToGhosts(sendMsg)
}

//ExchangeRealGhost 虚实交换
func (c *Cell) ExchangeRealGhost(entityID uint64) {
	if entity, ok := c.realEntityMap[entityID]; ok {
		c.ghostEntityMap[entityID] = entity
		delete(c.realEntityMap, entityID)

		return
	}

	if entity, ok := c.ghostEntityMap[entityID]; ok {
		c.realEntityMap[entityID] = entity
		delete(c.ghostEntityMap, entityID)
	}

}

func (c *Cell) GetSpaceCellInfos() *sync.Map {
	return c.space.getCellInfos()
}
