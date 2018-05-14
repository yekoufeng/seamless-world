package entity

import (
	"fmt"
	"sync"
	"zeus/dbservice"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/msghandler"

	"github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// iEntityState Entity的状态相关函数
type iEntityState interface {
	OnInit()
	OnAfterInit()
	OnLoop()
	OnDestroy()
	IsDestroyed() bool
}

// iEntityState 留给后代作一些回调方法
type iEntityInit interface {
	Init(interface{})
}

type iEntityLoop interface {
	Loop()
}

type iEntityDestroy interface {
	Destroy()
}

type iEntityInitFinish interface {
	InitFinish()
}

// delayedSendMsg 需要延迟发送的消息结构
type delayedSendMsg struct {
	srvType uint8
	msg     msgdef.IMsg
}

// Entity 代表服务器端一个可通讯对象
type Entity struct {
	msghandler.IMsgHandlers

	entityType  string
	entityID    uint64
	dbid        uint64
	cellID      uint64
	cellSrvType uint8

	// real实体所在的服务器id，仅对ghost实体有效
	realServerID uint64
	//real entity所在的cellID
	realCellID uint64

	initParam interface{}

	state byte

	srvIDS    map[uint8]*dbservice.EntitySrvInfo
	srvIDSMux *sync.RWMutex

	clientSess iserver.ISess

	isDebug bool

	realPtr     interface{}
	entitiesPtr iserver.IEntities
	ieState     iEntityState

	// props relate
	props            map[string]*PropInfo
	def              *Def
	dirtyPropList    []*PropInfo
	dirtySaveProps   map[string]*PropInfo
	dirtyProps       map[uint8][]*PropInfo
	ghostProps       []*PropInfo
	dirtyClientProps []*PropInfo
	dirtyMRoleProps  []*PropInfo
	// props end

	// 消息缓存, 在帧末发送的消息
	delayedMsgs []*delayedSendMsg
}

// SetClient 设置客户端的连接
func (e *Entity) SetClient(s iserver.ISess) {
	e.clientSess = s
	if s == nil {
		return
	}

	e.clientSess.SetMsgHandler(e.IMsgHandlers)
}

// GetClientSess 获取客户端连接
func (e *Entity) GetClientSess() iserver.ISess {

	if e.clientSess != nil && !e.clientSess.IsClosed() {
		return e.clientSess
	}

	return nil
}

// GetDBID 获取实体的DBID
func (e *Entity) GetDBID() uint64 {
	return e.dbid
}

// GetID 获取实体ID
func (e *Entity) GetID() uint64 {
	return e.entityID
}

// GetType 获取实体类型
func (e *Entity) GetType() string {
	return e.entityType
}

// GetRealPtr 获取真实的后代对象的指针
func (e *Entity) GetRealPtr() interface{} {
	return e.realPtr
}

// GetInitParam 获取初始化参数
func (e *Entity) GetInitParam() interface{} {
	return e.initParam
}

// GetEntityState 获取当前状态
func (e *Entity) GetEntityState() uint8 {
	return e.state
}

// IsGhost 是否是ghost实体
func (e *Entity) IsGhost() bool {
	return e.realServerID != 0
}

func (e *Entity) SetRealServerID(srvID uint64) {
	e.realServerID = srvID
}

// GetRealServerID 获取realServerID
func (e *Entity) GetRealServerID() uint64 {
	return e.realServerID
}

func (e *Entity) SetRealCellID(cellID uint64) {
	e.realCellID = cellID
}

// GetRealCellID
func (e *Entity) GetRealCellID() uint64 {
	return e.realCellID
}

// GetEntities 获取包含自己的Entities指针
func (e *Entity) GetEntities() iserver.IEntities {
	return e.entitiesPtr
}

// GetProxy 获取实体的代理对象
func (e *Entity) GetProxy() iserver.IEntityProxy {
	return NewEntityProxy(iserver.GetSrvInst().GetSrvID(), e.GetCellID(), e.entityID)
}

// OnEntityCreated 初始化
func (e *Entity) OnEntityCreated(entityID uint64, entityType string, dbid uint64, cellID uint64, protoType interface{}, entities iserver.IEntities, initParam interface{}, syncInit bool, realServerID uint64) {

	e.IMsgHandlers = msghandler.NewMsgHandlers()

	e.entityType = entityType
	e.entityID = entityID
	e.dbid = dbid
	e.cellID = cellID
	e.initParam = initParam

	e.realServerID = realServerID

	e.realPtr = protoType
	e.entitiesPtr = entities

	e.state = iserver.Entity_State_Init
	e.ieState = protoType.(iEntityState)

	e.srvIDS = make(map[uint8]*dbservice.EntitySrvInfo)
	e.srvIDSMux = &sync.RWMutex{}

	e.props = make(map[string]*PropInfo)
	e.dirtyPropList = make([]*PropInfo, 0, 1)
	e.dirtySaveProps = make(map[string]*PropInfo)
	e.dirtyProps = make(map[uint8][]*PropInfo)
	e.ghostProps = make([]*PropInfo, 0, 1)
	e.dirtyClientProps = make([]*PropInfo, 0, 1)
	e.dirtyMRoleProps = make([]*PropInfo, 0, 1)

	e.delayedMsgs = make([]*delayedSendMsg, 0, 1)

	ps, ok := e.realPtr.(iserver.IEntityPropsSetter)
	if ok {
		ps.SetPropsSetter(e.realPtr.(iserver.IEntityProps))
	}

	if syncInit {
		ies := e.ieState
		ies.OnInit()
		ies.OnAfterInit()
	}
}

// OnEntityDestroyed 当Entity销毁时调用
func (e *Entity) OnEntityDestroyed() {
	if e.state == iserver.Entity_State_Loop || e.state == iserver.Entity_State_Init {
		e.state = iserver.Entity_State_Destroy
		//e.reflushToDB()
	}
}

// IsDestroyed 是否删除
func (e *Entity) IsDestroyed() bool {
	return e.state == iserver.Entity_State_InValid
}

// MainLoop 主循环
func (e *Entity) MainLoop() {

	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err, e)
			if viper.GetString("Config.Recover") == "0" {
				panic(fmt.Sprintln(err, e))
			}
		}
	}()

	ies := e.ieState

	switch e.state {
	case iserver.Entity_State_Init:
		{
			ies.OnInit()
			ies.OnAfterInit()
		}
	case iserver.Entity_State_Loop:
		{
			ies.OnLoop()
		}
	case iserver.Entity_State_Destroy:
		{
			ies.OnDestroy()
			e.state = iserver.Entity_State_InValid
		}
	default:
		{
			// do nothing
		}
	}
}

// OnInit 初始化
func (e *Entity) OnInit() {
	e.state = iserver.Entity_State_Loop

	e.InitProp(GetDefs().GetDef(e.entityType))
	e.RegMsgProc(e.GetRealPtr())

	e.RegSrvID()
}

// OnAfterInit 后代的初始化
func (e *Entity) OnAfterInit() {
	ii, ok := e.GetRealPtr().(iEntityInit)
	if ok {
		ii.Init(e.GetInitParam())
	} else {
		seelog.Error("the entity ", e.GetType(), " no init method")
	}
}

// OnDestroy 销毁
func (e *Entity) OnDestroy() {

	seelog.Info("OnDestroy the entity ", e.GetType(), ", ID: ", e.GetID())

	if id, ok := e.GetRealPtr().(iEntityDestroy); ok {
		id.Destroy()
	}

	if !e.IsGhost() {
		e.UnregSrvID()
		e.reflushToDB()
	}

	if e.clientSess != nil {
		e.clientSess.Close()
	}

	e.state = iserver.Entity_State_InValid
}

// OnLoop 运行时
func (e *Entity) OnLoop() {
	// 每帧处理顺序:
	// 处理消息和业务逻辑, 在业务逻辑中更改属性, 发送RPC消息, 此时仅缓存需要发送的RPC消息
	// 同步属性变化消息
	// 真正发送RPC消息

	e.DoMsg()
	e.DoLooper()
	e.ReflushDirtyProp()
	e.FlushDelayedMsgs()
}

// DoLooper 调用后代类的loop函数
func (e *Entity) DoLooper() {
	if ie, ok := e.GetRealPtr().(iEntityLoop); ok {
		ie.Loop()
	}
}
