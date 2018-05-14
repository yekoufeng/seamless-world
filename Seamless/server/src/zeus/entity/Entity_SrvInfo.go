package entity

import (
	"fmt"
	"zeus/dbservice"
	"zeus/iserver"
	"zeus/msgdef"
)

func (e *Entity) isEntityExisted(srvType uint8) bool {

	e.srvIDSMux.RLock()
	defer e.srvIDSMux.RUnlock()

	_, ok := e.srvIDS[srvType]

	return ok
}

func (e *Entity) getEntitySrvID(srvType uint8) (uint64, uint64, error) {

	e.srvIDSMux.RLock()
	srvID, ok := e.srvIDS[srvType]
	e.srvIDSMux.RUnlock()

	if ok {
		return srvID.SrvID, srvID.CellID, nil
	}

	// 第一次尝试不成功, 则先刷新一次信息
	e.RefreshSrvIDS()

	e.srvIDSMux.RLock()
	srvID, ok = e.srvIDS[srvType]
	e.srvIDSMux.RUnlock()

	if !ok {
		return 0, 0, fmt.Errorf("Entity srvType [%d] not existed", srvType)
	}

	return srvID.SrvID, srvID.CellID, nil
}

// GetSrvIDS 获取玩家的分布式实体所在的服务器列表
func (e *Entity) GetSrvIDS() map[uint8]*dbservice.EntitySrvInfo {
	return e.srvIDS
}

// RefreshSrvIDS 从redis上更新 当前 entity所有的分布式信息
func (e *Entity) RefreshSrvIDS() {

	e.srvIDSMux.Lock()
	defer e.srvIDSMux.Unlock()

	srvIDs, err := dbservice.EntitySrvUtil(e.entityID).GetSrvIDs()
	if err != nil {
		e.Error("Get entity srv info failed")
		return
	}

	e.srvIDS = srvIDs

	//ghost不更新cellID
	if e.IsGhost() {
		return
	}

	e.cellID = 0
	e.cellSrvType = 0

	for srvType, info := range e.srvIDS {

		//e.Debug("RefreshSrvIDS,  CellID: ", info.CellID, "srvType: ", srvType)

		if info.CellID != 0 {
			e.cellID = info.CellID
			e.cellSrvType = srvType
		}
	}

}

// RegSrvID 将当前部分的Entity注册到Redis上
func (e *Entity) RegSrvID() {
	if !e.IsGhost() {
		if err := dbservice.EntitySrvUtil(e.entityID).RegSrvID(
			iserver.GetSrvInst().GetSrvType(),
			iserver.GetSrvInst().GetSrvID(),
			e.GetCellID(),
			e.entityType,
			e.dbid); err != nil {
			e.Error("Reg SrvID failed ", err)
			return
		}

		if iserver.GetSrvInst().GetSrvType() == iserver.ServerTypeGateway {
			if err := dbservice.EntitySrvUtil(e.entityID).RegSrvID(
				iserver.ServerTypeClient,
				iserver.GetSrvInst().GetSrvID(),
				e.GetCellID(),
				e.entityType,
				e.dbid); err != nil {
				e.Error("Reg SrvID failed ", err)
				return
			}
		}

		e.RefreshSrvIDS()
		e.broadcastSrvInfo()
	} else {
		e.RefreshSrvIDS()
	}
}

// unregSrvID 将当前部分的Entity从Redis上删除
func (e *Entity) UnregSrvID() {
	if err := dbservice.EntitySrvUtil(e.entityID).UnRegSrvID(
		iserver.GetSrvInst().GetSrvType(),
		iserver.GetSrvInst().GetSrvID(),
		e.GetCellID()); err != nil {
		e.Error("Unreg SrvID failed ", err)
	}

	if iserver.GetSrvInst().GetSrvType() == iserver.ServerTypeGateway {
		if err := dbservice.EntitySrvUtil(e.entityID).UnRegSrvID(
			iserver.ServerTypeClient,
			iserver.GetSrvInst().GetSrvID(),
			e.GetCellID()); err != nil {
			e.Error("Unreg SrvID failed ", err)
		}
	}

	e.broadcastSrvInfo()
}

func (e *Entity) broadcastSrvInfo() {

	e.srvIDSMux.RLock()
	defer e.srvIDSMux.RUnlock()

	for srvType := range e.srvIDS {
		if srvType != iserver.ServerTypeClient && srvType != iserver.GetSrvInst().GetSrvType() {
			e.Post(srvType, &msgdef.EntitySrvInfoNotify{})
		}
	}
}

func (e *Entity) MsgProc_EntitySrvInfoNotify(msgdef.IMsg) {
	e.RefreshSrvIDS()
}
