package entity

import (
	"reflect"
	"zeus/global"
	"zeus/iserver"

	log "github.com/cihub/seelog"
)

// ProtoType 保存实体映射列表
type ProtoType struct {
	entityProtoType map[string]reflect.Type
}

// NewProtoType 创建新的对象
func NewProtoType() *ProtoType {
	return &ProtoType{
		entityProtoType: make(map[string]reflect.Type),
	}
}

// RegProtoType 注册原型
func (es *ProtoType) RegProtoType(name string, protoType iserver.IEntity, autoCreate bool) {

	if protoType == nil {
		log.Error("ProtoType is nil", name)
		return
	}

	_, ok := es.entityProtoType[name]
	if ok {
		log.Warn("Registed type", name)
		return
	}

	t := reflect.TypeOf(protoType)
	es.entityProtoType[name] = t

	// only not space entity should be register to server type list
	//if _, ok := protoType.(iserver.ICellEntity); !ok {
	if autoCreate {
		// dbservice.EntityTypeUtil(name).RegSrvType(iserver.GetSrvInst().GetSrvType())
		types := global.GetGlobalInst().GetGlobalIntSlice("EntitySrvTypes:" + name)
		curType := iserver.GetSrvInst().GetSrvType()
		for _, typ := range types {
			if typ == int(curType) {
				return
			}
		}

		types = append(types, int(curType))
		global.GetGlobalInst().SetGlobalIntSlice("EntitySrvTypes:"+name, types)
	}
}

// NewEntityByProtoType 创建 一个对象
func (es *ProtoType) NewEntityByProtoType(entityType string) interface{} {
	proto, ok := es.entityProtoType[entityType]
	if !ok {
		log.Error("Entity type not exist", entityType)
		return nil
	}

	e := reflect.New(proto.Elem()).Interface()
	return e
}
