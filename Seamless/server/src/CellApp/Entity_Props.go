package main

import (
	"errors"
	"zeus/msgdef"
)

type iPropsSender interface {
	SendFullProps() error
}

var errPropsNil = errors.New("Props num is zero")

// SendFullProps 发送完整AOI属性信息
func (e *Entity) SendFullProps() error {
	num, data := e.GetAOIProp()
	if num == 0 {
		return errPropsNil
	}

	msg := &msgdef.PropsSyncClient{
		EntityID: e.GetID(),
		Num:      uint32(num),
		Data:     data,
	}
	e.CastMsgToMe(msg)
	return nil
}

// GetAOIProp 获得进入其它人AOI范围内需要收到的属性数据
func (e *Entity) GetAOIProp() (int, []byte) {
	return e.PackProps(true)
}

// GetAllProp 获得所有属性数据
func (e *Entity) GetAllProp() (int, []byte) {
	return e.PackProps(false)
}
