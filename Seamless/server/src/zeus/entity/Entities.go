package entity

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"zeus/dbservice"
	"zeus/events"
	"zeus/global"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/msghandler"
	"zeus/serializer"
	"zeus/timer"

	log "github.com/cihub/seelog"
)

type iEntityCtrl interface {
	FireMsg(name string, content interface{})
	GetClientSess() iserver.ISess
	OnEntityCreated(entityID uint64, entityType string, UID uint64, cellID uint64, protoType interface{}, entities iserver.IEntities, initParam interface{}, syncInit bool, realServerID uint64)
	OnEntityDestroyed()
	IsDestroyed() bool
	MainLoop()
}

// Entities entity的集合
type Entities struct {
	msghandler.IMsgHandlers
	*timer.Timer
	*events.GlobalEvents

	isMutilThread bool

	entityCnt int32

	entities       *sync.Map
	entitiesByType *sync.Map
	entitiesByDBID *sync.Map
}

// NewEntities 创建一个新的Entities
func NewEntities(isMutilThread bool) *Entities {
	return &Entities{
		IMsgHandlers: msghandler.NewMsgHandlers(),
		Timer:        timer.NewTimer(),
		GlobalEvents: events.NewGlobalEventsInst(),

		isMutilThread: isMutilThread,
		entityCnt:     0,

		entities:       &sync.Map{},
		entitiesByDBID: &sync.Map{},
		entitiesByType: &sync.Map{},
	}
}

// Init 初始化
func (es *Entities) Init() {
	es.RegMsgProc(es)
}

// Destroy 删除所有的实体
func (es *Entities) Destroy() {

	es.entities.Range(
		func(k, v interface{}) bool {
			if err := es.DestroyEntity(k.(uint64)); err != nil {
				log.Error(err)
			}
			return true
		})

	es.GlobalEvents.Destroy()

	es.entitiesByType = nil
	es.entitiesByDBID = nil
	es.entities = nil
}

// SyncDestroy 删除所有Entity, 并等待所有Entity删除结束
func (es *Entities) SyncDestroy() {
	es.entities.Range(
		func(k, v interface{}) bool {
			e := es.GetEntity(k.(uint64))
			if err := es.DestroyEntity(k.(uint64)); err != nil {
				log.Error(err)
			}
			for {
				if e.(iEntityState).IsDestroyed() {
					break
				}

				time.Sleep(1 * time.Millisecond)
			}
			return true
		})

	es.GlobalEvents.Destroy()
	es.entitiesByType = nil
	es.entitiesByDBID = nil
	es.entities = nil
}

// MainLoop 自己的逻辑线程
func (es *Entities) MainLoop() {

	es.DoMsg()

	es.Timer.Loop()
	es.GlobalEvents.HandleEvent()

	if !es.isMutilThread {
		es.entities.Range(func(k, v interface{}) bool {
			v.(iEntityCtrl).MainLoop()
			return true
		})
	}
}

// Range 遍历所有Entity
func (es *Entities) Range(f func(k, v interface{}) bool) {
	es.entities.Range(f)
}

// CreateEntity 创建实体
func (es *Entities) CreateEntity(entityType string, entityID uint64, dbid uint64, cellID uint64, initParam interface{}, syncInit bool, realServerID uint64) (iserver.IEntity, error) {

	_, ok := es.entities.Load(entityID)

	if ok {
		return nil, fmt.Errorf("EntityID existed")
	}

	ie := iserver.GetSrvInst().NewEntityByProtoType(entityType).(iEntityCtrl)

	if cellID != 0 {
		_, ok := ie.(iserver.ICellEntity)
		if !ok {
			return nil, fmt.Errorf("The entity must is a space entity")
		}
	}

	ie.OnEntityCreated(entityID, entityType, dbid, cellID, ie, es, initParam, syncInit, realServerID)
	es.addEntity(entityID, ie.(iserver.IEntity))

	if es.isMutilThread {
		go func() {
			ticker := time.NewTicker(iserver.GetSrvInst().GetFrameDeltaTime())
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if ie.IsDestroyed() {
						return
					}
					ie.MainLoop()
				}
			}
		}()
	}

	atomic.AddInt32(&es.entityCnt, 1)
	return ie.(iserver.IEntity), nil
}

// DestroyEntityByDBID 删除 Entity
func (es *Entities) DestroyEntityByDBID(entityType string, dbID uint64) error {
	e := es.GetEntityByDBID(entityType, dbID)
	if e == nil {
		return fmt.Errorf("Entity not existed")
	}

	return es.DestroyEntity(e.GetID())
}

// DestroyEntity 删除Entity
func (es *Entities) DestroyEntity(entityID uint64) error {
	e, ok := es.entities.Load(entityID)
	if !ok {
		return fmt.Errorf("Entity not existed")
	}

	es.delEntity(e.(iserver.IEntity))
	e.(iEntityCtrl).OnEntityDestroyed()

	if !es.isMutilThread {
		e.(iEntityState).OnDestroy()
	}

	atomic.AddInt32(&es.entityCnt, -1)
	return nil
}

// GetEntity 获取Entity
func (es *Entities) GetEntity(entityID uint64) iserver.IEntity {
	if ie, ok := es.entities.Load(entityID); ok {
		return ie.(iserver.IEntity)
	}
	return nil
}

// GetEntityByDBID 获取Entity
func (es *Entities) GetEntityByDBID(entityType string, dbid uint64) iserver.IEntity {
	if it, ok := es.entitiesByDBID.Load(entityType); ok {

		if id, ok := it.(*sync.Map).Load(dbid); ok {
			return id.(iserver.IEntity)
		}
	}
	return nil
}

// TravsalEntity 遍历某一类型的entity
func (es *Entities) TravsalEntity(entityType string, f func(iserver.IEntity)) {
	if it, ok := es.entitiesByType.Load(entityType); ok {
		it.(*sync.Map).Range(func(k, v interface{}) bool {
			ise := v.(iserver.IEntityStateGetter)
			if ise.GetEntityState() != iserver.Entity_State_Loop {
				return true
			}

			f(v.(iserver.IEntity))
			return true
		})
	}
}

// addEntity  增加entity
func (es *Entities) addEntity(entityID uint64, e iserver.IEntity) {

	es.entities.Store(entityID, e)

	if e.GetDBID() != 0 {
		var t *sync.Map
		if it, ok := es.entitiesByDBID.Load(e.GetType()); ok {
			t = it.(*sync.Map)
		} else {
			t = &sync.Map{}
			es.entitiesByDBID.Store(e.GetType(), t)
		}

		t.Store(e.GetDBID(), e)
	}

	var t *sync.Map
	if it, ok := es.entitiesByType.Load(e.GetType()); ok {
		t = it.(*sync.Map)
	} else {
		t = &sync.Map{}
		es.entitiesByType.Store(e.GetType(), t)
	}

	t.Store(e.GetID(), e)
}

// delEntity 删除entity
func (es *Entities) delEntity(e iserver.IEntity) {

	es.entities.Delete(e.GetID())

	if e.GetDBID() != 0 {
		if it, ok := es.entitiesByDBID.Load(e.GetType()); ok {
			it.(*sync.Map).Delete(e.GetDBID())
		}
	}

	if it, ok := es.entitiesByType.Load(e.GetType()); ok {
		it.(*sync.Map).Delete(e.GetID())
	}
}

// EntityCount 返回实体数
func (es *Entities) EntityCount() uint32 {
	return uint32(es.entityCnt)
}

// CreateEntityAll 创建实体的所有部分
func (es *Entities) CreateEntityAll(entityType string, dbid uint64, initParam interface{}, syncInit bool) (iserver.IEntity, error) {

	var entityID uint64
	entityID = iserver.GetSrvInst().FetchTempID()

	// 理论上tempID不可能重复
	/*
		if dbservice.EntitySrvUtil(entityID).IsExist() {
			return nil, fmt.Errorf("Entity existed")
		}
	*/

	e, err := es.CreateEntity(entityType, entityID, dbid, 0, initParam, syncInit, 0)
	if err != nil {
		return nil, err
	}

	srvList := global.GetGlobalInst().GetGlobalIntSlice("EntitySrvTypes:" + entityType)
	// srvList, err := dbservice.EntityTypeUtil(entityType).GetSrvType()
	// if err != nil {
	// 	es.DestroyEntity(e.GetID())
	// 	return nil, err
	// }

	for i := 0; i < len(srvList); i++ {
		srvType := uint8(srvList[i])
		if srvType == iserver.GetSrvInst().GetSrvType() {
			continue
		}

		srvID, err := iserver.GetSrvInst().GetSrvIDBySrvType(srvType)
		if err != nil {
			log.Error(err)
			continue
		}

		//提前注册，这样就可以提前发消息了
		// dbservice.EntitySrvUtil(entityID).RegSrvID(srvType, srvID, 0, entityType, dbid)

		// data, err := common.Marshal(initParam)
		// if err != nil {
		// 	log.Error("marshal init error ", err)
		// 	return nil, err
		// }

		msg := &msgdef.CreateEntityReq{
			EntityType: entityType,
			EntityID:   entityID,
			CellID:     0,
			InitParam:  serializer.Serialize(initParam),
			DBID:       dbid,
			SrcSrvType: iserver.GetSrvInst().GetSrvType(),
			SrcSrvID:   iserver.GetSrvInst().GetSrvID(),
			CallbackID: 0,
		}

		if err := iserver.GetSrvInst().PostMsgToCell(srvID, 0, msg); err != nil {
			log.Error(err)
		}
	}

	return e, nil
}

// DestroyEntityAll 销毁实体的所有部分
func (es *Entities) DestroyEntityAll(entityID uint64) error {

	var srvInfos map[uint8]*dbservice.EntitySrvInfo
	var err error

	e := es.GetEntity(entityID)
	if e == nil {
		srvInfos, err = dbservice.EntitySrvUtil(entityID).GetSrvIDs()
		if err != nil {
			log.Error("Get entity srv info failed ", err)
			return err
		}
	} else {
		srvInfos = e.GetSrvIDS()
	}

	// srvInfos, err := dbservice.EntitySrvUtil(entityID).GetSrvIDs()
	// if err != nil {
	// 	log.Error("Get entity srv info failed ", err)
	// 	return err
	// }

	for _, srvInfo := range srvInfos {
		if srvInfo.SrvID == iserver.GetSrvInst().GetSrvID() {
			continue
		}

		msg := &msgdef.DestroyEntityReq{
			EntityID:   entityID,
			SrcSrvType: iserver.GetSrvInst().GetSrvType(),
			SrcSrvID:   iserver.GetSrvInst().GetSrvID(),
			CallbackID: 0,
			CellID:     srvInfo.CellID,
		}

		if err := iserver.GetSrvInst().PostMsgToCell(srvInfo.SrvID, srvInfo.CellID, msg); err != nil {
			log.Error(err)
		}
	}

	return es.DestroyEntity(entityID)
}
