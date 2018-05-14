package main

import (
	"reflect"
	"zeus/common"
	"zeus/iserver"
	"zeus/linmath"

	"github.com/cihub/seelog"
)

/*
	TinyEntity 轻量级的Entity 在场景中可以大量放置的Entity
	只有最简单的标志自己位置相关的能力
	所有的属性必须在后代的 Init函数中完成
	没有独立的Loop相关方法，属性也无法改变
*/

// ITinyEntity TinyEntity接口
type ITinyEntity interface {
	onInit()
	onDestroy()

	onEntityCreated(entityID uint64, entityType string, cell iserver.ICell, initParam interface{}, realPtr interface{})
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

// TinyEntity 轻量级的Entity
type TinyEntity struct {
	id      uint64
	typ     string
	realPtr interface{}
	//cell     iserver.ICell
	cell      *Cell
	initParam interface{}

	pos       linmath.Vector3
	rota      linmath.Vector3
	packState []byte

	props     map[string]*tinyPropInfo
	packProps []byte
}

type tinyPropInfo struct {
	typ   reflect.Kind
	value interface{}
}

func (e *TinyEntity) onEntityCreated(entityID uint64, entityType string, cell iserver.ICell, initParam interface{}, realPtr interface{}) {
	e.id = entityID
	e.typ = entityType
	e.realPtr = realPtr
	//e.space = space
	e.initParam = initParam

	ps, ok := e.realPtr.(iserver.IEntityPropsSetter)
	if ok {
		ps.SetPropsSetter(e.realPtr.(iserver.IEntityProps))
	}

	e.props = make(map[string]*tinyPropInfo)
}

// GetID 获取ID号
func (e *TinyEntity) GetID() uint64 {
	return e.id
}

// GetType 获取类型
func (e *TinyEntity) GetType() string {
	return e.typ
}

// GetRealPtr 获取后代类
func (e *TinyEntity) GetRealPtr() interface{} {
	return e.realPtr
}

// GetCell 获取Cell
func (e *TinyEntity) GetCell() *Cell {
	return e.cell
}

// GetInitParam 获取初始化参数
func (e *TinyEntity) GetInitParam() interface{} {
	return e.initParam
}

func (e *TinyEntity) onInit() {
	ii, ok := e.GetRealPtr().(iEntityInit)
	if ok {
		ii.Init(e.GetInitParam())
	}

	e.initProps()
	e.initState()

	if e.GetCell() != nil {
		e.GetCell().UpdateCoord(e)
	}
}

func (e *TinyEntity) onDestroy() {
	ii, ok := e.GetRealPtr().(iEntityDestroy)
	if ok {
		ii.Destroy()
	}

	if e.GetCell() != nil {
		e.GetCell().RemoveFromCoord(e)
	}
}

// SetProp 设置属生
func (e *TinyEntity) SetProp(name string, v interface{}) {
	e.props[name] = &tinyPropInfo{reflect.TypeOf(v).Kind(), v}
}

// PropDirty 设置属生为空
func (e *TinyEntity) PropDirty(name string) {

}

// GetProp 获取属性
func (e *TinyEntity) GetProp(name string) interface{} {
	return e.props[name]
}

// SetPos 设置位置
func (e *TinyEntity) SetPos(pos linmath.Vector3) {
	e.pos = pos
}

// GetPos 获取位置
func (e *TinyEntity) GetPos() linmath.Vector3 {
	return e.pos
}

// SetRota 设置旋转
func (e *TinyEntity) SetRota(rota linmath.Vector3) {
	e.rota = rota
}

// GetRota 获取旋转
func (e *TinyEntity) GetRota() linmath.Vector3 {
	return e.rota
}

// GetAOIProp 获取属性的打包
func (e *TinyEntity) GetAOIProp() (int, []byte) {
	return len(e.props), e.packProps
}

//IsNearAOILayer 是否近的层级
func (e *TinyEntity) IsNearAOILayer() bool {
	return true
}

// IsAOITrigger 是否AOI触发
func (e *TinyEntity) IsAOITrigger() bool {
	return false
}

// IsWatcher 不是观察者
func (e *TinyEntity) IsWatcher() bool {
	return false
}

func (e *TinyEntity) initProps() {
	size := 0
	for name, value := range e.props {
		size = size + len(name) + 2
		size = size + e.getValueStreamSize(value)
	}

	if size == 0 {
		return
	}

	e.packProps = make([]byte, size)
	bs := common.NewByteStream(e.packProps)

	for name, value := range e.props {
		if err := bs.WriteStr(name); err != nil {
			seelog.Error(err, name, e)
		}
		if err := e.writeValueToStream(bs, value); err != nil {
			seelog.Error(err, name, e)
		}
	}
}

func (e *TinyEntity) getValueStreamSize(v *tinyPropInfo) int {

	s := 0
	switch v.typ {
	case reflect.Int8, reflect.Uint8:
		s = 1
	case reflect.Int16, reflect.Uint16:
		s = 2
	case reflect.Int32, reflect.Uint32:
		s = 4
	case reflect.Int64, reflect.Uint64:
		s = 8
	case reflect.Float32:
		s = 4
	case reflect.Float64:
		s = 8
	case reflect.String:
		s = len(v.value.(string)) + 2
	default:
		panic("no support type")
	}

	typeName := v.typ.String()

	return s + len(typeName) + 2
}

func (e *TinyEntity) writeValueToStream(bs *common.ByteStream, v *tinyPropInfo) error {

	var err error

	typeName := v.typ.String()
	err = bs.WriteStr(typeName)
	if err != nil {
		return err
	}

	switch v.typ {
	case reflect.Int8:
		err = bs.WriteInt8(v.value.(int8))
	case reflect.Uint8:
		err = bs.WriteByte(v.value.(uint8))
	case reflect.Int16:
		err = bs.WriteInt16(v.value.(int16))
	case reflect.Uint16:
		err = bs.WriteUInt16(v.value.(uint16))
	case reflect.Int32:
		err = bs.WriteInt32(v.value.(int32))
	case reflect.Uint32:
		err = bs.WriteUInt32(v.value.(uint32))
	case reflect.Int64:
		err = bs.WriteInt64(v.value.(int64))
	case reflect.Uint64:
		err = bs.WriteUInt64(v.value.(uint64))
	case reflect.String:
		err = bs.WriteStr(v.value.(string))
	case reflect.Float32:
		err = bs.WriteFloat32(v.value.(float32))
	case reflect.Float64:
		err = bs.WriteFloat64(v.value.(float64))
	default:
		panic("no support type")
	}

	return err
}

func (e *TinyEntity) initState() {
	// e.packState = make([]byte, 4+4+6*4)
	// bs := common.NewByteStream(e.packState)

	// var mask uint32
	// mask = EntityStateMask_Pos_X | EntityStateMask_Pos_Y | EntityStateMask_Pos_Z | EntityStateMask_Rota_X | EntityStateMask_Rota_Y | EntityStateMask_Rota_Z

	// bs.WriteUInt32(e.GetCell().GetTimeStamp())
	// bs.WriteUInt32(mask)
	// bs.WriteFloat32(e.pos.X)
	// bs.WriteFloat32(e.pos.Y)
	// bs.WriteFloat32(e.pos.Z)

	// bs.WriteFloat32(e.rota.X)
	// bs.WriteFloat32(e.rota.Y)
	// bs.WriteFloat32(e.rota.Z)
}

func (e *TinyEntity) IsGhost() bool {
	return false
}
