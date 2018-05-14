package main

import (
	"zeus/dbservice"
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
)

func (e *TinyEntity) GetDBID() uint64 {
	return 0
}

func (e *TinyEntity) GetCellID() uint64 {
	if e.cell == nil {
		return 0
	}

	return e.GetCell().GetID()
}

func (e *TinyEntity) GetSrvIDS() map[uint8]*dbservice.EntitySrvInfo {
	return nil
}

func (e *TinyEntity) IsOwnerCellEntity() bool {
	return true
}

func (e *TinyEntity) IsCellEntity() bool {
	return true
}

func (e *TinyEntity) Post(srvType uint8, msg msgdef.IMsg) error {
	return nil
}

func (e *TinyEntity) RPC(srvType uint8, methodName string, args ...interface{}) error {
	return nil
}

/*func (e *TinyEntity) RPCOther(srvType uint8, srcEntityID uint64, methodName string, args ...interface{}) error {
	return nil
}*/

func (e *TinyEntity) EnterCell(cellID uint64) {

}

func (e *TinyEntity) LeaveCell() {

}

func (e *TinyEntity) GetProxy() iserver.IEntityProxy {
	return nil
}

func (e *TinyEntity) SetCoordPos(pos linmath.Vector3) {}
