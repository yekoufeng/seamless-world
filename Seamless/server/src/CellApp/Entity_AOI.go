package main

import (
	"common"
	"math"
	"protoMsg"
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
)

const (
	// aoiRange     = 50.0
	aoiTolerance = 1.0
)

type iAOIPacker interface {
	GetID() uint64
	GetType() string
	GetAOIProp() (int, []byte)
}

// iAOISender 断线重连之后打包发送完整AOI信息
type iAOISender interface {
	SendFullAOIs() error
}

func (e *Entity) SendCellInfos() {

	cellinfos := e.cell.GetSpaceCellInfos()
	cellinfos.Range(
		func(id, cellinfoI interface{}) bool {
			cellinfo := cellinfoI.(*common.CellInfo)
			if cellinfo != nil {
				var rect protoMsg.RectInfo
				rect.Xmin = cellinfo.GetRect().Xmin
				rect.Xmax = cellinfo.GetRect().Xmax
				rect.Ymin = cellinfo.GetRect().Ymin
				rect.Ymax = cellinfo.GetRect().Ymax
				msg := &protoMsg.CellInfoNotify{
					Operate:  1,
					SpaceID:  e.cell.getSpace().getID(),
					CellID:   id.(uint64),
					RectInfo: &rect,
					SrvID:    cellinfo.GetCellSrvID(),
				}
				e.Post(iserver.ServerTypeClient, msg)
			}
			return true
		})

}

// SendFullAOIs 发送完整的AOI信息
func (e *Entity) SendFullAOIs() error {
	msg := msgdef.NewEntityAOISMsg()

	if e.GetCell() == nil {
		return nil
	}

	e.GetCell().TravsalAOI(e, func(o iserver.ICoordEntity) {
		ip, ok := o.(iAOIPacker)
		if !ok {
			e.Error("Get AOIPacker failed")
			return
		}

		num, propBytes := ip.GetAOIProp()
		m := &msgdef.EnterAOI{
			EntityID:   ip.GetID(),
			EntityType: ip.GetType(),
			PropNum:    uint16(num),
			Properties: propBytes,
			Pos:        o.GetPos(),
			Rota:       o.GetRota(),
		}

		data := make([]byte, m.Size()+1)
		data[0] = 1
		m.MarshalTo(data[1:])

		msg.AddData(data)
	})

	return e.Post(iserver.ServerTypeClient, msg)
}

// SetWatcher 设置当前entity 为watcher
func (e *Entity) SetWatcher(b bool) {
	e._isWatcher = b
}

// IsWatcher 是否观察者
func (e *Entity) IsWatcher() bool {
	return e._isWatcher
}

func (e *Entity) isBeWatch() bool {
	return e.beWatchedNums > 0
}

func (e *Entity) determineAOIFlag() {
	updataDist := 10.0 //aoiRange * 0.01
	if math.Abs(float64(e.pos.X-e.lastAOIPos.X)) > updataDist ||
		math.Abs(float64(e.pos.Y-e.lastAOIPos.Y)) > updataDist ||
		math.Abs(float64(e.pos.Z-e.lastAOIPos.Z)) > updataDist {
		e.needUpdateAOI = true
	}
}

func (e *Entity) updatePosCoord(pos linmath.Vector3) {

	if e.needUpdateAOI {
		c := e.GetCell()
		if c != nil {
			c.UpdateCoord(e)
		}

		e.lastAOIPos = pos
		e.needUpdateAOI = false
	}
}

// AddExtWatchEntity 增加额外关注对象
func (e *Entity) AddExtWatchEntity(o iserver.ICoordEntity) {
	if e.extWatchList == nil {
		e.extWatchList = make(map[uint64]*extWatchEntity)
	}

	if _, ok := e.extWatchList[o.GetID()]; ok {
		return
	}

	inMyAOI := false
	if e.GetCell() != nil {
		e.GetCell().TravsalAOI(e, func(n iserver.ICoordEntity) {
			// 已经在AOI范围内
			if n.GetID() == o.GetID() {
				inMyAOI = true
			}
		})
	}

	if !inMyAOI {
		e.OnEntityEnterAOI(o)
	}

	e.extWatchList[o.GetID()] = &extWatchEntity{
		entity:  o,
		isInAOI: inMyAOI,
	}
}

// RemoveExtWatchEntity 删除额外关注对象
func (e *Entity) RemoveExtWatchEntity(o iserver.ICoordEntity) {
	if e.extWatchList == nil {
		return
	}

	if _, ok := e.extWatchList[o.GetID()]; !ok {
		return
	}

	inMyAOI := false

	if e.GetCell() != nil {
		e.GetCell().TravsalAOI(e, func(n iserver.ICoordEntity) {
			// 已经在AOI范围内
			if n.GetID() == o.GetID() {
				inMyAOI = true
			}
		})
	}

	delete(e.extWatchList, o.GetID())

	if !inMyAOI {
		e.OnEntityLeaveAOI(o)
	}
}

func (e *Entity) clearExtWatchs() {

	for id, we := range e.extWatchList {
		delete(e.extWatchList, id)
		if !we.isInAOI {
			e.OnEntityLeaveAOI(we.entity)
		}
	}
}

// TravsalExtWatchs 遍历额外观察者列表
func (e *Entity) TravsalExtWatchs(f func(*extWatchEntity)) {
	if len(e.extWatchList) == 0 {
		return
	}

	for _, extWatch := range e.extWatchList {
		if !extWatch.isInAOI {
			f(extWatch)
		}
	}
}

//OnEntityEnterAOI 实体进入AOI范围
func (e *Entity) OnEntityEnterAOI(o iserver.ICoordEntity) {
	// 当o在我的额外关注列表中时, 不触发真正的EnterAOI, 只是打个标记
	if extWatch, ok := e.extWatchList[o.GetID()]; ok {
		extWatch.isInAOI = true
		return
	}

	if e._isWatcher {
		e.aoies = append(e.aoies, AOIInfo{true, o})

		e.Info("OnEntityEnterAOI  append, isGhost: ", e.IsGhost(), " o.ID: ", o.GetID(), ", o.isGhost: ", o.IsGhost())
	}

	if o.IsWatcher() {
		e.beWatchedNums++
	}

	e.Info("OnEntityEnterAOI  isGhost: ", e.IsGhost(), " o.ID: ", o.GetID(), ", o.isGhost: ", o.IsGhost())
}

//OnEntityLeaveAOI 实体离开AOI范围
func (e *Entity) OnEntityLeaveAOI(o iserver.ICoordEntity) {
	// 当o在我的额外关注列表中时, 不触发真正的LeaveAOI
	if extWatch, ok := e.extWatchList[o.GetID()]; ok {
		extWatch.isInAOI = false
		return
	}

	if e._isWatcher {
		e.aoies = append(e.aoies, AOIInfo{false, o})

		e.Info("OnEntityLeaveAOI append,  isGhost: ", e.IsGhost(), " o.ID: ", o.GetID(), ", o.isGhost: ", o.IsGhost())
	}

	if o.IsWatcher() {
		e.beWatchedNums--
	}

	e.Info("OnEntityLeaveAOI  isGhost: ", e.IsGhost(), " o.ID: ", o.GetID(), ", o.isGhost: ", o.IsGhost())
}

func (e *Entity) updateAOI() {

	if len(e.aoies) != 0 && e._isWatcher {
		msg := msgdef.NewEntityAOISMsg()
		for i := 0; i < len(e.aoies); i++ {

			if msg.Num >= 20 {
				e.PostToClient(msg)
				msg = msgdef.NewEntityAOISMsg()
			}

			info := e.aoies[i]

			ip := info.entity.(iAOIPacker)

			var data []byte

			if info.isEnter {
				num, propBytes := ip.GetAOIProp()
				m := &msgdef.EnterAOI{
					EntityID:   ip.GetID(),
					EntityType: ip.GetType(),
					PropNum:    uint16(num),
					Properties: propBytes,
					Pos:        info.entity.GetPos(),
					Rota:       info.entity.GetRota(),
				}

				data = make([]byte, m.Size()+1)
				data[0] = 1
				m.MarshalTo(data[1:])

			} else {
				m := &msgdef.LeaveAOI{
					EntityID: ip.GetID(),
				}

				data = make([]byte, m.Size()+1)
				data[0] = 0
				m.MarshalTo(data[1:])
			}

			msg.AddData(data)
		}

		e.PostToClient(msg)
		e.aoies = e.aoies[0:0]
	}
}

//IsNearAOILayer 是否视野近的层
func (e *Entity) IsNearAOILayer() bool {
	return false
}

//IsAOITrigger 是否要解发AOITrigger事件
func (e *Entity) IsAOITrigger() bool {
	return true //e.IsWatcher()
}
