package iserver

import (
	"zeus/linmath"
)

// ICell 代表一个空间
type ICell interface {
	IEntity
	ICoord
	IEntities

	GetTimeStamp() uint32

	AddEntity(entityType string, entityID uint64, dbid uint64, initParam interface{}, syncInit bool, isGhost bool) error
	RemoveEntity(entityID uint64) error

	AddTinyEntity(entityType string, entityID uint64, initParam interface{}) error
	RemoveTinyEntity(entityID uint64) error

	IsMapLoaded() bool
	FindPath(srcPos, destPos linmath.Vector3) ([]linmath.Vector3, error)
	Raycast(origin, direction linmath.Vector3, length float32, mask int32) (float32, linmath.Vector3, int32, bool, error)
	CapsuleRaycast(head, foot linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool, error)
	SphereRaycast(center linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool, error)
	GetHeight(x, z float32) (float32, error)
	IsWater(x, z float32, waterlevel float32) (bool, error)

	//TravsalEntityInRange(ICellEntity, float32, func(ICellEntity))
}
