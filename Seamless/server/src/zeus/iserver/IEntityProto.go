package iserver

// IEntityProto 实体原型
type IEntityProto interface {
	RegProtoType(name string, protoType IEntity, autoCreate bool)
	NewEntityByProtoType(entityType string) interface{}
}
