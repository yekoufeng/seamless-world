package iserver

// IEntityProps Entity属性相关的操作
type IEntityProps interface {
	SetProp(name string, value interface{})
	PropDirty(name string)
	GetProp(name string) interface{}
}

// IEntityPropsSetter Entity属性相关的操作
type IEntityPropsSetter interface {
	SetPropsSetter(IEntityProps)
}
