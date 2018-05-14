package main

import (
	"common"
	"db"
	"protoMsg"
	"zeus/entity"
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
	"zeus/serializer"
)

// IEntityCtrl 内部使用接口
type IEntityCtrl interface {
	onEnterCell()
	onLeaveCell()
}

// IEnterCell 上层回调
type IEnterCell interface {
	OnEnterCell()
}

// ILeaveCell 上层回调
type ILeaveCell interface {
	OnLeaveCell()
}

// ILeaveCell 上层回调
type IGetEntity interface {
	GetEntity() *Entity
}

// iGameStateSender 发送完整的游戏状态信息
type iGameStateSender interface {
	SendFullGameState()
}

// iLateLooper 内部接口, Entity后处理
type iLateLooper interface {
	onLateLoop()
}

// iAOIUpdater 刷新AOI接口
type iAOIUpdater interface {
	updateAOI()
}

// IWatcher 观察者
type IWatcher interface {
	iserver.IEntity
	iserver.IPos

	PostToClient(msgdef.IMsg) error

	/*
		getWatchAOIRange() float32
		markerLeaveAOI(IMarker)
		markerEnterAOI(IMarker)

		addMarker(IMarker)
		removeMarker(IMarker)
		isExistMarker(IMarker) bool
	*/
}

// // IAOIEnterTrigger 进入aoi区域
// type IAOIEnterTrigger interface {
// 	OnMarkerEnterAOI(m IMarker)
// }

// // IAOILeaveTrigger 离开aoi区域
// type IAOILeaveTrigger interface {
// 	OnMarkerLeaveAOI(m IMarker)
// }

// // IMarker 被观察者
// type IMarker interface {
// 	iserver.IEntity
// 	iserver.IPos
// 	IAOIProp
// 	addWatcher(IWatcher)
// 	removeWatcher(IWatcher)
// 	isExistWatcher(IWatcher) bool
// 	IsMarker() bool
// 	getMarkAOIRange() float32
// }

const (
	EnterAOILength     = 15.0 //AOI的进入范围
	LeaveAOILength     = 16.0 //AOI的退出范围
	TransferRealLength = 1.0  //转移real的长度
)

// Entity 空间中的实体
type Entity struct {
	entity.Entity

	cell *Cell

	pos  linmath.Vector3 //坐标
	rota linmath.Vector3 //旋转

	lastAOIPos    linmath.Vector3
	needUpdateAOI bool
	_isWatcher    bool //是否是观察者

	//临时保存离开或者进入aoi实体列表，更新后清空
	aoies []AOIInfo
	//被观察的数量
	beWatchedNums int

	extWatchList map[uint64]*extWatchEntity

	aoiSyncMsg *msgdef.AOISyncUserState

	// CastToAll相关的消息缓存
	delayedCastMsgs []*delayedCastMsg

	//如果是real则指针有效
	realEntity *RealEntity
}

// AOIInfo aoi信息
type AOIInfo struct {
	isEnter bool
	entity  iserver.ICoordEntity
}

// extWatchEntity 额外关注列表
type extWatchEntity struct {
	entity iserver.ICoordEntity

	isInAOI bool // 是否在AOI范围内
}

// OnInit 构造函数
func (e *Entity) OnInit() {
	e.Entity.OnInit()

	e.pos = linmath.Vector3_Invalid()
	e.lastAOIPos = linmath.Vector3_Invalid()
	e.needUpdateAOI = false

	e.aoies = make([]AOIInfo, 0, 5)

	e.delayedCastMsgs = make([]*delayedCastMsg, 0, 1)

	e._isWatcher = false

	if !e.IsGhost() {
		realEntity := CreateRealEntity()
		e.SetRealEntity(realEntity)
	}

	e.aoiSyncMsg = msgdef.NewAOISyncUserState()
	e.RegMsgProc(&EntityMsgProc{e: e})
	cell, ok := GetSrvInst().Cells.Load(e.Entity.GetCellID())
	if ok {
		e.cell = cell.(*Cell)
	} else {
		e.Error("cell is null , cellID:  ", e.Entity.GetCellID())
	}
}

// OnAfterInit 后代的初始化
func (e *Entity) OnAfterInit() {
	e.Entity.OnAfterInit()
	e.onEnterCell()
	//e.updatePosCoord(e.pos)
}

// OnDestroy 析构函数
func (e *Entity) OnDestroy() {
	//保存一些信息到数据中
	if !e.IsGhost() {
		e.SaveMapData()

		//删除所有ghost
		if e.realEntity != nil {
			for cellID, haunt := range e.realEntity.hauntMap {
				e.DeleteGhost(cellID, haunt.serverid)
			}
		}
	}

	e.onLeaveCell()

	e.Entity.OnDestroy()
}

// SaveMapData 保存玩家地图信息
func (e *Entity) SaveMapData() {
	data := &db.PlayerMapData{
		ID:      e.GetDBID(),
		MapName: e.GetCell().GetMapName(),
		Pos:     e.GetPos(),
		Rota:    e.GetRota(),
	}

	db.PlayerMapUtil(e.GetDBID()).SetPlayerMapData(data)
}

// GetSpace 获取所在的空间
func (e *Entity) GetSpace() iserver.ICell {
	return nil
}

// GetCell 获取所在的Cell
func (e *Entity) GetCell() *Cell {
	return e.cell
}

func (e *Entity) onEnterCell() {
	e.Debug("onEnterCell, isGhost: ", e.IsGhost())

	if e.IsGhost() {
		return
	}

	ic, ok := e.GetRealPtr().(IEnterCell)
	if ok {
		ic.OnEnterCell()
	}

	if e.IsWatcher() {
		msg := &msgdef.EnterCell{
			CellID:   e.GetCell().GetID(),
			MapName:  e.GetCell().GetMapName(),
			EntityID: e.GetID(),
			Addr:     iserver.GetSrvInst().GetCurSrvInfo().OuterAddress,
			//TimeStamp: e.GetSpace().GetTimeStamp(),
		}
		if err := e.Post(iserver.ServerTypeClient, msg); err != nil {
			e.Error("Send EnterCell failed ", err)
		}

		sendMsg := &protoMsg.EnterCellOk{
			//MapID:   e.GetCell().GetMapName(),
			MapName: e.GetCell().GetMapName(),
			Pos:     &protoMsg.Vector3{e.GetPos().X, e.GetPos().Y, e.GetPos().Z},
			Rota:    &protoMsg.Vector3{e.GetRota().X, e.GetRota().Y, e.GetRota().Z},
			CellID:  e.GetCellID(),
		}

		if err := e.Post(iserver.ServerTypeClient, sendMsg); err != nil {
			e.Error("Send EnterCellOk failed ", err)
		}

		e.aoies = append(e.aoies, AOIInfo{true, e})
	}
}

func (e *Entity) onLeaveCell() {
	e.Debug("onLeaveCell, isGhost: ", e.IsGhost())
	ic, ok := e.GetRealPtr().(ILeaveCell)
	if ok {
		ic.OnLeaveCell()
	}

	if e.GetCell() != nil {
		e.GetCell().RemoveFromCoord(e)
	}

	if e.IsWatcher() && !e.IsGhost() {
		e.aoies = append(e.aoies, AOIInfo{false, e})
		e.clearExtWatchs()
		e.updateAOI()

		msg := &msgdef.LeaveCell{}
		if err := e.Post(iserver.ServerTypeClient, msg); err != nil {
			e.Error("Send LeaveCell failed ", err)
		}
	}
}

// Entity帧处理顺序
// 处理消息和业务逻辑, 在业务逻辑中会有RPC和CastToAll
// 更新坐标系中的位置
// Space更新所有Entity的AOI
// Space调用所有Entity的LateUpdate, 发送属性消息, 延迟发送消息和缓存的CastToAll消息

// OnLoop 循环调用
func (e *Entity) OnLoop() {

	e.Entity.DoMsg()
	e.Entity.DoLooper()
	e.updatePosCoord(e.pos)
	// e.updateAOI()

	// e.updateAOI()
	// e.resetState()
	// e.Entity.OnLoop()
	// e.updatePosCoord(e.pos)
	// e.updateState()
}

// onLateLoop 后处理
func (e *Entity) onLateLoop() {
	e.Entity.ReflushDirtyProp()
	e.Entity.FlushDelayedMsgs()

	// 真正发送所有消息
	e.FlushDelayedCastMsgs()
}

func (e *Entity) syncClock() {
	e.PostToClient(&msgdef.SyncClock{
		TimeStamp: e.GetCell().GetTimeStamp(),
	})
}

// IsCellEntity 是否是个SpaceEntity
func (e *Entity) IsCellEntity() bool {
	return true
}

// GetRealEntity 获取realEntity
func (e *Entity) GetRealEntity() *RealEntity {
	return e.realEntity
}

// ClearRealEntity 设置real实体
func (e *Entity) SetRealEntity(real *RealEntity) {
	e.realEntity = real
}

//CheckReals 检查real
func (e *Entity) CheckReals() {
	if e.realEntity == nil {
		return
	}

	//先检查ghost，再判断迁移，顺序不能变
	e.CheckGhosts()

	if e.realEntity == nil {
		e.Info("CheckReals, realEntity is nil ,CheckGhosts中修改了")
		return
	}

	//检查是否要迁移real
	e.CheckTransferReal()
}

//CheckGhosts 检查ghost的创建和删除
func (e *Entity) CheckGhosts() {
	cell := e.GetCell()

	//判断是不是进入无ghost区域
	ret := cell.cellInfo.GetRect().IsInInnerRect(float64(e.GetPos().X), float64(e.GetPos().Z), LeaveAOILength)
	if ret {
		for cellID, haunt := range e.realEntity.hauntMap {
			e.DeleteGhost(cellID, haunt.serverid)
		}

		return
	}

	// ret = cell.cellInfo.GetRect().IsInInnerRect(float64(e.GetPos().X), float64(e.GetPos().Z), EnterAOILength)
	// if ret {
	// 	return
	// }

	var cellID uint64
	var cellInfo *common.CellInfo

	space := cell.getSpace()
	if space == nil {
		e.Info("CheckGhosts, space is nil ")
		return
	}

	space.cellinfos.Range(
		func(k, v interface{}) bool {
			cellID = k.(uint64)
			cellInfo = v.(*common.CellInfo)

			if cellInfo == nil {
				e.Error("CheckGhosts, cellInfo is nil , cellID: ", cellID)
				return true
			}

			//是自己则忽略
			if cellID == e.GetCellID() {
				return true
			}

			if e.realEntity == nil {
				e.Error("CheckGhosts, realEntity is nil , cellID: ", cellID)
				return false
			}

			if e.realEntity.hauntMap == nil {
				e.Error("CheckGhosts, e.realEntity.hauntMap is nil , cellID: ", cellID)
				return false
			}

			//获取是否已经创建ghost
			_, ok := e.realEntity.hauntMap[cellInfo.GetCellID()]

			if ok {
				//离开了AOI区域，需要删除ghost
				if !cellInfo.GetRect().IsInOuterRect(float64(e.GetPos().X), float64(e.GetPos().Z), LeaveAOILength) {
					e.DeleteGhost(cellInfo.GetCellID(), cellInfo.GetCellSrvID())
				}
			} else {
				//进入AOI范围但是还没有创建ghost，则创建新ghost
				if cellInfo.GetRect().IsInOuterRect(float64(e.GetPos().X), float64(e.GetPos().Z), EnterAOILength) {
					e.CreateGhost(cellInfo.GetCellID(), cellInfo.GetCellSrvID())
				}
			}

			return true
		})
}

//CheckTransferReal 检查是否迁移real到其他cell
func (e *Entity) CheckTransferReal() {
	cell := e.GetCell()
	ret := cell.cellInfo.GetRect().IsInOuterRect(float64(e.GetPos().X), float64(e.GetPos().Z), TransferRealLength)
	if !ret {
		var cellID uint64
		var cellInfo *common.CellInfo

		space := cell.getSpace()
		if space == nil {
			e.Info("CheckTransferReal, space is nil ")
			return
		}

		//real实体已经超出了，准备迁移
		//首先判断real落在哪个cell中，如果这个cell已经有ghost就直接迁移，如果没有则先创建ghost再迁移

		space.cellinfos.Range(
			func(k, v interface{}) bool {
				cellID = k.(uint64)
				cellInfo = v.(*common.CellInfo)

				if cellInfo == nil {
					e.Error("CheckTransferReal, cellInfo is nil , cellID: ", cellID)
					return true
				}

				//是自己则忽略
				if cellID == e.GetCellID() {
					return true
				}

				if e.realEntity == nil {
					e.Error("CheckTransferReal, realEntity is nil , cellID: ", cellID)
					return false
				}

				if e.realEntity.hauntMap == nil {
					e.Error("CheckTransferReal, e.realEntity.hauntMap is nil , cellID: ", cellID)
					return false
				}

				//如果没有ghost则创建
				_, ok := e.realEntity.hauntMap[cellInfo.GetCellID()]

				//判断是否进入其他cell的范围
				if cellInfo.GetRect().IsInInnerRect(float64(e.GetPos().X), float64(e.GetPos().Z), 0) {
					//如果没有ghost则创建
					if !ok {
						e.Info("CheckTransferReal, 没有创建ghost, cellID: ", cellID, ", severID: ", cellInfo.GetCellSrvID(), ", entityID:", e.GetID())
						e.CreateGhost(cellInfo.GetCellID(), cellInfo.GetCellSrvID())
					}

					//开始转移
					e.TransferReal(cellInfo.GetCellID(), cellInfo.GetCellSrvID())
					return false
				}

				return true
			})
	}
}

//DeleteGhost 删除ghost
func (e *Entity) DeleteGhost(cellID, serverID uint64) {

	e.Info("DeleteGhost, cellID: ", cellID, ", severID: ", serverID, ", entityID:", e.GetID())
	msg := &msgdef.DeleteGhostReq{
		EntityID: e.GetID(),
	}

	if err := iserver.GetSrvInst().PostMsgToCell(serverID, cellID, msg); err != nil {
		e.Error("DeleteGhostReq failed: ", err)
	}

	delete(e.realEntity.hauntMap, cellID)
}

//创建ghost
func (e *Entity) CreateGhost(cellID, serverID uint64) {
	e.Info("CreateGhost, cellID: ", cellID, ", severID: ", serverID, ", entityID:", e.GetID())

	num, data := e.GetAllProp()

	msg := &msgdef.CreateGhostReq{
		EntityType:   "Player",
		EntityID:     e.GetID(),
		DBID:         e.GetDBID(),
		InitParam:    serializer.Serialize(e.GetInitParam()),
		CellID:       cellID,
		Pos:          e.GetPos(),
		RealServerID: srvInst.GetCurSrvInfo().ServerID,
		RealCellID:   e.GetCellID(),
		PropNum:      uint32(num),
		Props:        data,
	}

	if err := iserver.GetSrvInst().PostMsgToCell(serverID, cellID, msg); err != nil {
		e.Error("CreateGhostReq failed: ", err)
	}

	haunt := &Haunt{serverid: serverID}
	e.realEntity.hauntMap[cellID] = haunt
}

//TransferReal 删除ghost
func (e *Entity) TransferReal(cellID, serverID uint64) {
	e.Info("TransferReal, cellID: ", cellID, ", severID: ", serverID, ", entityID:", e.GetID())

	//删除即将要变成real的ghost记录
	delete(e.realEntity.hauntMap, cellID)
	//加入当前的即将变成real的记录
	e.realEntity.hauntMap[e.GetCellID()] = &Haunt{serverid: srvInst.GetCurSrvInfo().ServerID}

	hauntNum, hauntData := e.GetRealEntity().packHaunt()

	num, data := e.GetAllProp()

	msg := &msgdef.TransferRealReq{
		EntityID:     e.GetID(),
		PropNum:      uint32(num),
		Props:        data,
		Pos:          e.GetPos(),
		RealServerID: srvInst.GetCurSrvInfo().ServerID,
		GhostNum:     uint32(hauntNum),
		GhostData:    hauntData,
	}

	if err := iserver.GetSrvInst().PostMsgToCell(serverID, cellID, msg); err != nil {
		e.Error("TransferRealReq failed: ", err)
		return
	}

	//不用注销，防止出现一段时间没有real
	//e.UnregSrvID()

	//转移了，需要把当前的real改成ghost
	e.SetRealServerID(serverID)
	e.SetRealCellID(cellID)
	e.SetRealEntity(nil)

	e.GetCell().ExchangeRealGhost(e.GetID())
}

func (e *Entity) SyncToGhosts(msg msgdef.IMsg) {
	if e.IsGhost() || e.realEntity == nil {
		return
	}

	srvType := iserver.GetSrvInst().GetSrvType()
	for cellID, info := range e.realEntity.hauntMap {
		if packMsg, err := e.ExportPackMsg(srvType, cellID, msg); err == nil {
			iserver.GetSrvInst().PostMsgToSrv(info.serverid, packMsg)
		}
	}
}

//
func (e *Entity) GetEntity() *Entity {
	return e
}

//SendMsgToReal 发送消息给real 实现ISendMsgToReal接口
func (e *Entity) SendMsgToReal(msg msgdef.IMsg) {
	if !e.IsGhost() || e.GetRealServerID() == 0 {
		return
	}

	if packMsg, err := e.ExportPackMsg(iserver.GetSrvInst().GetSrvType(), e.GetRealCellID(), msg); err == nil {
		iserver.GetSrvInst().PostMsgToSrv(e.GetRealServerID(), packMsg)
	}
}

// SetPos 设置位置
func (e *Entity) SetPos(pos linmath.Vector3) {
	e.SetCoordPos(pos)
}

// GetPos 获取位置
func (e *Entity) GetPos() linmath.Vector3 {
	return e.pos
}

// SetRota 设置旋转
func (e *Entity) SetRota(rota linmath.Vector3) {
	e.rota = rota
}

// GetRota 获取旋转
func (e *Entity) GetRota() linmath.Vector3 {
	return e.rota
}
