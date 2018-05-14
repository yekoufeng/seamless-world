package iserver

import (
	"zeus/events"
	"zeus/msghandler"
	"zeus/timer"
)

// IEntities 用于实体的管理
type IEntities interface {
	msghandler.IMsgHandlers
	timer.ITimer
	events.IEvents

	CreateEntityAll(entityType string, dbid uint64, initParam interface{}, syncInit bool) (IEntity, error)
	DestroyEntityAll(entityID uint64) error

	CreateEntity(entityType string, entityID uint64, dbid uint64, cellID uint64, initParam interface{}, syncInit bool, realServerID uint64) (IEntity, error)

	DestroyEntity(entityID uint64) error
	DestroyEntityByDBID(entityType string, dbID uint64) error

	GetEntity(entityID uint64) IEntity
	GetEntityByDBID(entityType string, dbid uint64) IEntity

	TravsalEntity(entityType string, f func(IEntity))

	EntityCount() uint32
}
