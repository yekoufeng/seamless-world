package main

import (
	"protoMsg"
	"time"
	"zeus/linmath"
	"zeus/msgdef"

	"fmt"

	log "github.com/cihub/seelog"
)

type ClientMsgProc struct {
}

func (proc *ClientMsgProc) MsgProc_ClientVertifySucceedRet(content interface{}) {
	var msg *msgdef.ClientVertifySucceedRet = content.(*msgdef.ClientVertifySucceedRet)
	log.Debug("MsgProc_ClientVertifySucceedRet:", *msg)

	time.Sleep(1e9)
	GetClient().TryEnterCell()

	// go GetMainView().Start()
	// viewStart <- struct{}{}
}

func (proc *ClientMsgProc) MsgProc_EnterCellOk(content interface{}) {
	GetMainView().EnterCellOk()
}

func (proc *ClientMsgProc) MsgProc_RPCMsg(content msgdef.IMsg) {
	msg, ok := content.(*msgdef.RPCMsg)
	if !ok {
		return
	}
	// 从客户端收到的消息要判断要调用的服务器类型
	GetClient().GetSess().FireRPC(msg.MethodName, msg.Data)
	// log.Debug("MsgProc_RPCMsg")
}

func (proc *ClientMsgProc) MsgProc_ProtoSync(content interface{}) {
	// msg, ok := content.(*msgdef.ProtoSync)
	// if !ok {
	// 	log.Error("MsgProc_ProtoSync failed")
	// 	return
	// }

	// log.Debug("MsgProc_ProtoSync")
}

func (proc *ClientMsgProc) MsgProc_MRolePropsSyncClient(content interface{}) {
	msg, ok := content.(*msgdef.MRolePropsSyncClient)
	if !ok {
		return
	}

	log.Debug("MsgProc_MRolePropsSyncClient, props'num:", msg.Num)
	GetClient().EntityID = msg.EntityID
	GetClient().ParseProps(uint16(msg.Num), msg.Data)
	GetClient().AddUser(GetClient().User)
}

func (proc *ClientMsgProc) MsgProc_PropsSync(content interface{}) {
	msg, ok := content.(*msgdef.PropsSync)
	if !ok {
		return
	}

	log.Debug("MsgProc_PropsSync, props'num:", msg.Num)
	GetClient().ParseProps(uint16(msg.Num), msg.Data)
}

func (proc *ClientMsgProc) MsgProc_PropsSyncClient(content interface{}) {
	msg, ok := content.(*msgdef.PropsSyncClient)
	if !ok {
		return
	}

	log.Debug("MsgProc_PropsSyncClient, props'num:", msg.Num)
	u, ok := GetClient().GetUser(msg.EntityID)
	if !ok {
		return
	}
	u.ParseProps(uint16(msg.Num), msg.Data)
}

func (proc *ClientMsgProc) MsgProc_EntityAOIS(content interface{}) {
	msg, ok := content.(*msgdef.EntityAOIS)
	if !ok {
		return
	}

	log.Debug("MsgProc_EntityAOIS, ", msg.Num)

	datas := msg.GetData()
	for _, data := range datas {
		aoi := &msgdef.EnterAOI{}
		aoi.Unmarshal(data[1:])
		GetClient().AddAOI(data[0], aoi)
	}
}

func (proc *ClientMsgProc) MsgProc_AOISyncUserState(content interface{}) {
	log.Debug("MsgProc_AOISyncUserState")
	// msg, ok := content.(*msgdef.AOISyncUserState)
	// if !ok {
	// 	return
	// }
	// for i := uint32(0); i < msg.Num; i++ {
	// 	eid := msg.EIDS[i]
	// 	data := msg.EDS[i]
	// 	if u, ok := GetClient().GetUser(eid); ok {
	// 		u.ParseEntityState(data)
	// 		u.StartMove()
	// 		u.lastMoveTime = time.Now().UnixNano() / 1e6
	// 	}
	// }
}

func (proc *ClientMsgProc) MsgProc_MoveUpdate(content interface{}) {
	msg, ok := content.(*protoMsg.MoveUpdate)
	if !ok {
		return
	}
	// log.Debug("MsgProc_MoveUpdate", msg)

	if u, ok := GetClient().GetUser(msg.EntityID); ok {
		if msg.Pos != nil {
			u.SetPos(linmath.Vector3{msg.Pos.X * MapRate, msg.Pos.Y * MapRate, msg.Pos.Z * MapRate})
		}
		if msg.Rota != nil {
			u.SetRota(linmath.Vector3{msg.Rota.X, msg.Rota.Y, msg.Rota.Z})
		}
		if msg.Stoped {
			u.um.model.StopAnimate()
		} else {
			u.um.model.StartAnimate()
		}
	}
}

func (proc *ClientMsgProc) MsgProc_DetectCell(content interface{}) {
	// log.Debug("MsgProc_DetectCell")
	GetClient().SetCellOk()
}

func (proc *ClientMsgProc) MsgProc_CellInfoNotify(content interface{}) {
	msg, ok := content.(*protoMsg.CellInfoNotify)
	if !ok {
		return
	}

	ri := msg.GetRectInfo()
	// 将CellInfo存入CellInfoMgr
	outStr := fmt.Sprintf("CellInfoNotify:SrvID:%d CellID:%d SpaceID:%d Xmin:%f Xmax:%f, Ymin:%f, Ymax:%f",
		msg.GetSrvID(), msg.GetCellID(), msg.GetSpaceID(), ri.Xmin, ri.Xmax, ri.Ymin, ri.Ymax)
	log.Debug(outStr)
	GetCellInfoMgr().AddCell(msg.GetSpaceID(), msg.GetCellID(), msg.GetSrvID(),
		ri.GetXmin()*MapRate, ri.GetXmax()*MapRate, ri.GetYmin()*MapRate, ri.GetYmax()*MapRate)
}

func (proc *ClientMsgProc) MsgProc_EffectNotify(content interface{}) {
	msg, ok := content.(*protoMsg.EffectNotify)
	if !ok {
		return
	}
	u, ok := GetClient().GetUser(msg.EntityID)
	if !ok {
		return
	}
	u.AddEffect(msg)
}

func (proc *ClientMsgProc) RPC_SyncFriendList(*protoMsg.SyncFriendList) {

}
func (proc *ClientMsgProc) RPC_SyncApplyList(msg *protoMsg.SyncFriendApplyList) {

}
func (proc *ClientMsgProc) RPC_InitOwnGoodsInfo(*protoMsg.OwnGoodsInfo) {

}

func (proc *ClientMsgProc) RPC_InitNotifyMysqlDbAddr(addr string) {

}

func (proc *ClientMsgProc) RPC_SrvTime(svrTime int64) {

}

// ClientVertifyReq 验证消息
type ClientVertifyReq struct {
	// Source: 消息来源, 分客户端或者服务器(ClientMSG/服务器类型)
	Source uint8 //data[0]
	// UID: 玩家UID或者服务器ID
	UID uint64 //data[1:9]
	// Token: 客户端登录时需要携带Token
	Token string //data[9:41]
}
