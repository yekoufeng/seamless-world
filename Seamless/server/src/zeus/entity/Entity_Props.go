package entity

import (
	"fmt"
	"reflect"
	"zeus/common"
	"zeus/dbservice"
	"zeus/iserver"
	"zeus/msgdef"
)

// InitProp 初始化属性列表
func (e *Entity) InitProp(def *Def) {
	if def == nil {
		return
	}

	e.def = def

	for _, p := range e.def.Props {

		isAdd := false

		for _, s := range p.InterestSrvs {
			if s == iserver.GetSrvInst().GetSrvType() {
				isAdd = true
				break
			}
		}

		if isAdd {
			e.addProp(p)
		}
	}

	e.reflushFromDB()
}

func (e *Entity) addProp(prop *PropDef) {
	e.props[prop.Name] = newPropInfo(prop)
}

// SetProp 设置一个属性的值
func (e *Entity) SetProp(name string, v interface{}) {
	p := e.props[name]

	if e.def == nil {
		panic(fmt.Errorf("no def file exist %s", name))
	}

	if !p.def.IsValidValue(v) {
		panic(fmt.Errorf("prop type error %s", name))
	}

	if reflect.DeepEqual(p.value, v) {
		return
	}

	p.value = v

	e.PropDirty(name)
}

// GetProp 获取一个属生的值
func (e *Entity) GetProp(name string) interface{} {
	return e.props[name].value
}

// PropDirty 设置某属性为Dirty
func (e *Entity) PropDirty(name string) {
	//如果是ghost则不用通知客户端或这保存到数据库
	if e.IsGhost() {
		return
	}

	p := e.props[name]

	if p.def.Sync {
		if p.syncFlag == false {
			e.dirtyPropList = append(e.dirtyPropList, p)
		}
		p.syncFlag = true
	}

	if p.def.Persistence {
		if p.dbFlag == false {
			e.dirtySaveProps[name] = p
		}
		p.dbFlag = true
	}
}

// ReflushDirtyProp 每一帧刷新属性通知
func (e *Entity) ReflushDirtyProp() {
	e.reflushSync()
	e.reflushToDB()
}

func (e *Entity) getDirtyProps(srvType uint8) []*PropInfo {
	m, ok := e.dirtyProps[srvType]
	if !ok {
		m = make([]*PropInfo, 0, 10)
		e.dirtyProps[srvType] = m
	}
	return m
}

func (e *Entity) sendPropsSyncMsg() {

	for s, m := range e.dirtyProps {
		if len(m) != 0 {
			if err := e.Post(s, &msgdef.PropsSync{Num: uint32(len(m)), Data: e.PackPropsToBytes(m)}); err != nil {
				e.Error("Send PropsSync failed ", err)
			}
			e.dirtyProps[s] = e.dirtyProps[s][0:0]
		}
	}

	if len(e.ghostProps) != 0 {
		if iSyncToGhost, ok := e.GetRealPtr().(iserver.ISyncToGhosts); ok {
			iSyncToGhost.SyncToGhosts(&msgdef.PropsSync{Num: uint32(len(e.ghostProps)), Data: e.PackPropsToBytes(e.ghostProps)})
		}
		e.ghostProps = e.ghostProps[0:0]
	}

	bc, ok := e.GetRealPtr().(iserver.IClientBroadcaster)
	if len(e.dirtyClientProps) != 0 {
		msg := &msgdef.PropsSyncClient{
			EntityID: e.GetID(),
			Num:      uint32(len(e.dirtyClientProps)),
			Data:     e.PackPropsToBytes(e.dirtyClientProps),
		}
		if e.GetCellID() != 0 {
			if ok {
				bc.CastMsgToAllClientExceptMe(msg)
			}
		} else {
			e.PostToCell(msg)
		}
		e.dirtyClientProps = e.dirtyClientProps[0:0]
	}

	if len(e.dirtyMRoleProps) != 0 {
		msg := &msgdef.PropsSyncClient{
			EntityID: e.GetID(),
			Num:      uint32(len(e.dirtyMRoleProps)),
			Data:     e.PackPropsToBytes(e.dirtyMRoleProps),
		}
		if ok {
			bc.CastMsgToMe(msg)
		} else {
			e.Post(iserver.ServerTypeClient, msg)
		}
		e.dirtyMRoleProps = e.dirtyMRoleProps[0:0]
	}
}

func (e *Entity) reflushSync() {
	if len(e.dirtyPropList) == 0 {
		return
	}

	for _, p := range e.dirtyPropList {
		for _, s := range p.def.InterestSrvs {
			if s != iserver.GetSrvInst().GetSrvType() && e.isEntityExisted(s) {
				m := e.getDirtyProps(s)
				e.dirtyProps[s] = append(m, p)
			}
		}

		if p.def.IsClientInterest {
			e.dirtyClientProps = append(e.dirtyClientProps, p)
			e.dirtyMRoleProps = append(e.dirtyMRoleProps, p)
		} else if p.def.IsMRoleInterest {
			e.dirtyMRoleProps = append(e.dirtyMRoleProps, p)
		}

		e.ghostProps = append(e.ghostProps, p)

		p.syncFlag = false
	}

	e.dirtyPropList = e.dirtyPropList[0:0]

	e.sendPropsSyncMsg()
}

// ReflushFromMsg 从消息中更新属性
func (e *Entity) ReflushFromMsg(num int, data []byte) {
	bs := common.NewByteStream(data)
	for i := 0; i < num; i++ {
		name, err := bs.ReadStr()
		if err != nil {
			e.Error("read prop name fail ", err)
			return
		}

		prop, ok := e.props[name]
		if !ok {
			e.Error("target entity not own prop ", name)
			return
		}

		err = prop.ReadValueFromStream(bs)
		if err != nil {
			e.Error("read prop from stream failed ", name, err)
			return
		}
	}
}

// 从数据库中恢复
func (e *Entity) reflushFromDB() {
	if e.dbid == 0 {
		return
	}

	var hashArgs []interface{}
	for n := range e.props {
		hashArgs = append(hashArgs, n)
	}

	if len(hashArgs) > 0 {
		values, err := dbservice.EntityUtil(e.entityType, e.dbid).GetValues(hashArgs)
		if err != nil {
			e.Error("Reflush from db failed ", err)
			return
		}

		for index, n := range hashArgs {
			info := e.props[n.(string)]
			info.UnPackValue(values[index])
		}
	}
}

func (e *Entity) reflushToDB() {
	if e.dbid == 0 {
		return
	}

	var hashArgs []interface{}
	for n, p := range e.dirtySaveProps {
		hashArgs = append(hashArgs, n)
		hashArgs = append(hashArgs, p.PackValue())

		p.dbFlag = false
	}

	if len(hashArgs) > 0 {
		if err := dbservice.EntityUtil(e.entityType, e.dbid).SetValues(hashArgs); err != nil {
			e.Error("Reflush to db failed ", err)
		}
	}

	e.dirtySaveProps = make(map[string]*PropInfo)
}

// PackProps 打包属性
func (e *Entity) PackProps(isClientInterest bool) (int, []byte) {

	props := make([]*PropInfo, 0, 10)

	for _, prop := range e.props {
		if isClientInterest && !prop.def.IsClientInterest {
			continue
		}
		props = append(props, prop)
	}

	return len(props), e.PackPropsToBytes(props)
}

// PackMRoleProps 打包主角关心的属性
func (e *Entity) PackMRoleProps() (int, []byte) {

	props := make([]*PropInfo, 0, 10)

	for _, prop := range e.props {
		if !prop.def.IsMRoleInterest {
			continue
		}
		props = append(props, prop)
	}

	return len(props), e.PackPropsToBytes(props)
}

// PackPropsToBytes 把属列表打包成bytes
func (e *Entity) PackPropsToBytes(props []*PropInfo) []byte {
	size := 0

	for _, prop := range props {
		size = size + len(prop.def.Name) + 2
		size = size + prop.GetValueStreamSize()
	}

	if size == 0 {
		return nil
	}

	buf := make([]byte, size)
	bs := common.NewByteStream(buf)

	for _, prop := range props {
		if err := bs.WriteStr(prop.def.Name); err != nil {
			e.Error("Pack props failed ", err, prop.def.Name)
		}
		if err := prop.WriteValueToStream(bs); err != nil {
			e.Error("Pack props failed ", err, prop.def.Name)
		}
	}

	return buf
}
