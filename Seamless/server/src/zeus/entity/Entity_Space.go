package entity

// LeaveCell 离开场景
func (e *Entity) LeaveCell() {
	// if e.IsCell() {
	// 	e.Warn("Space entity couldn't move into space")
	// 	return
	// }

	// if e.GetCellID() == 0 {
	// 	e.Warn("Entity not in space")
	// 	return
	// }

	// srvID, err := dbservice.CellUtil(e.GetCellID()).GetSrvID()
	// if err != nil {
	// 	e.Error("Get srvID error ", err)
	// 	return
	// }

	// msg := &msgdef.LeaveCellReq{
	// 	EntityID: e.entityID,
	// }

	// if err := iserver.GetSrvInst().PostMsgToCell(srvID, e.GetCellID(), msg); err != nil {
	// 	e.Error("Leave space failed ", err)
	// }
}

// GetCellID
func (e *Entity) GetCellID() uint64 {
	return e.cellID
}

// IsOwnerCellEntity 是否拥有SpaceEntity的部分
func (e *Entity) IsOwnerCellEntity() bool {
	return e.cellID != 0
}

// IsCell 是否是空间
func (e *Entity) IsCell() bool {
	return false
}

// IsCellEntity 是否是个空间实体
func (e *Entity) IsCellEntity() bool {
	return false
}
