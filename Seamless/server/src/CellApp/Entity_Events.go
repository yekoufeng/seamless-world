package main

import (
	"zeus/msgdef"
	"zeus/serializer"
)

// BroadcastEvent 广播事件
func (e *Entity) BroadcastEvent(event string, args ...interface{}) {
	msg := &msgdef.EntityEvent{}
	msg.SrcEntityID = e.GetID()
	msg.EventName = event
	msg.Data = serializer.Serialize(args...)

	e.CastMsgToAllClient(msg)
}

// BroadcastEventExceptMe 广播事件
func (e *Entity) BroadcastEventExceptMe(event string, args ...interface{}) {
	msg := &msgdef.EntityEvent{}
	msg.SrcEntityID = e.GetID()
	msg.EventName = event
	msg.Data = serializer.Serialize(args...)

	e.CastMsgToAllClientExceptMe(msg)
}
