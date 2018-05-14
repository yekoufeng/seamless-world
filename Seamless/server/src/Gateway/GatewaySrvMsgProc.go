package main

// GatewaySrvMsgProc 网关服务器消息处理类
type GatewaySrvMsgProc struct {
}

func (p *GatewaySrvMsgProc) MsgProc_SessVertified(content interface{}) {

	uid := content.(uint64)
	GetUserMgr().login(uid)
}
