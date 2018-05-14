package iserver

import (
	"zeus/linmath"
)

// ICoordEntity 坐标结点
type ICoordEntity interface {
	GetID() uint64
	SetPos(pos linmath.Vector3)
	SetCoordPos(pos linmath.Vector3)
	GetPos() linmath.Vector3
	SetRota(rota linmath.Vector3)
	GetRota() linmath.Vector3

	GetRealPtr() interface{}

	IsNearAOILayer() bool
	IsAOITrigger() bool
	IsWatcher() bool
	IsGhost() bool
}

// IAOITriggerEntity 拥有AOITrigger能力的Entity
type IAOITriggerEntity interface {
	OnEntityEnterAOI(ICoordEntity)
	OnEntityLeaveAOI(ICoordEntity)
}

// ICoord 座标系统
type ICoord interface {
	UpdateCoord(ICoordEntity)
	RemoveFromCoord(ICoordEntity)
	TravsalAOI(ICoordEntity, func(ICoordEntity))

	// 遍历center为中心，半径范围radius内的所有实体执行cb
	TravsalRange(center *linmath.Vector3, radius int, cb func(ICoordEntity))

	// 遍历center所在的Tower，在该Tower内的center为中心，半径范围radius内的所有实体执行cb
	TravsalCenter(center *linmath.Vector3, radius int, cb func(ICoordEntity))
}
