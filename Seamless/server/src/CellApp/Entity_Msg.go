package main

import (
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
	"zeus/serializer"
)

type delayedCastMsg struct {
	msg        msgdef.IMsg
	isCastToMe bool
}

// CastMsgToAllClient 发送消息给所有关注我的客户端
func (e *Entity) CastMsgToAllClient(msg msgdef.IMsg) {
	e.delayedCastMsgs = append(e.delayedCastMsgs, &delayedCastMsg{
		msg:        msg,
		isCastToMe: true,
	})
}

// CastMsgToMe 发送消息给自己
func (e *Entity) CastMsgToMe(msg msgdef.IMsg) {
	e.Post(iserver.ServerTypeClient, msg)
}

// CastMsgToAllClientExceptMe  发送给除了自己外的其它人
func (e *Entity) CastMsgToAllClientExceptMe(msg msgdef.IMsg) {
	e.delayedCastMsgs = append(e.delayedCastMsgs, &delayedCastMsg{
		msg:        msg,
		isCastToMe: false,
	})
}

func (e *Entity) CastMsgToRangeExceptMe(center *linmath.Vector3, radius int, msg msgdef.IMsg) {
	if e.GetCell() == nil {
		return
	}
	e.GetCell().TravsalRange(center, radius, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.Entity_State_Loop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok && ie.GetID() != e.GetID() {
				ie.Post(iserver.ServerTypeClient, msg)
			}
		}
	})
}

func (e *Entity) CastMsgToCenterExceptMe(center *linmath.Vector3, radius int, msg msgdef.IMsg) {
	if e.GetCell() == nil {
		return
	}

	e.GetCell().TravsalCenter(center, radius, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.Entity_State_Loop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok && ie.GetID() != e.GetID() {
				ie.Post(iserver.ServerTypeClient, msg)
			}
		}
	})
}

// CastRPCToAllClient 触发RPC消息给所有关注我的客户端
func (e *Entity) CastRPCToAllClient(methodName string, args ...interface{}) {
	data := serializer.Serialize(args...)
	msg := &msgdef.RPCMsg{}
	msg.ServerType = iserver.ServerTypeClient
	msg.MethodName = methodName
	msg.Data = data

	e.delayedCastMsgs = append(e.delayedCastMsgs, &delayedCastMsg{
		msg:        msg,
		isCastToMe: true,
	})
}

// CastRPCToMe 触发自己的RPC消息
func (e *Entity) CastRPCToMe(methodName string, args ...interface{}) {
	e.RPC(iserver.ServerTypeClient, methodName, args...)
}

// CastRPCToAllClientExceptMe 触发除了自己以外的其它客户端的RPC消息
func (e *Entity) CastRPCToAllClientExceptMe(methodName string, args ...interface{}) {
	data := serializer.Serialize(args...)
	msg := &msgdef.RPCMsg{}
	msg.ServerType = iserver.ServerTypeClient
	msg.MethodName = methodName
	msg.Data = data

	e.delayedCastMsgs = append(e.delayedCastMsgs, &delayedCastMsg{
		msg:        msg,
		isCastToMe: false,
	})
}

// PostToClient 投递消息给客户端
func (e *Entity) PostToClient(msg msgdef.IMsg) error {
	return e.Post(iserver.ServerTypeClient, msg)
}

// FlushDelayedCastMsgs 发送所有缓冲的Cast消息
func (e *Entity) FlushDelayedCastMsgs() {
	if len(e.delayedCastMsgs) == 0 {
		return
	}

	if e.GetCell() == nil {
		return
	}

	for _, dcm := range e.delayedCastMsgs {
		// 填充RPC消息中的SrcEntityID字段
		if rpcMsg, ok := dcm.msg.(*msgdef.RPCMsg); ok {
			rpcMsg.SrcEntityID = e.GetID()
		}

		e.GetCell().TravsalAOI(e, func(ia iserver.ICoordEntity) {
			if ise, ok := ia.(iserver.IEntityStateGetter); ok {
				if ise.GetEntityState() != iserver.Entity_State_Loop {
					return
				}

				if ie, ok := ia.(iserver.IEntity); ok && (e.GetID() != ie.GetID() || dcm.isCastToMe) {
					ie.Post(iserver.ServerTypeClient, dcm.msg)
				}
			}
		})

		e.TravsalExtWatchs(func(o *extWatchEntity) {
			if ise, ok := o.entity.(iserver.IEntityStateGetter); ok {
				if ise.GetEntityState() != iserver.Entity_State_Loop {
					return
				}

				if ie, ok := o.entity.(iserver.IEntity); ok {
					ie.Post(iserver.ServerTypeClient, dcm.msg)
				}
			}
		})
	}

	e.delayedCastMsgs = e.delayedCastMsgs[0:0]
}
