package main

import (
	"zeus/msgdef"
	"zeus/sess"
)

// GateUserMsgProc GateUser的消息处理函数
type GateUserMsgProc struct {
	user *GateUser
}

func (p *GateUserMsgProc) MsgProc_SessClosed(content interface{}) {
	uid := content.(uint64)
	GetUserMgr().logout(uid)
}

func (p *GateUserMsgProc) MsgProc_MsgForward(content interface{}) {
	rawMsg := content.(*sess.RawMsg)

	/*
		messageID 本身在消息中只占两个字节，也就是最大范围是65565
		所以消息号的取值范围是 1 - 65565之间

		11000 以内保留，用作框架内部消息定义
		11000 - 60000 之间用来上层协议使用
	*/

	srvType := int(rawMsg.MsgID / 1000)
	if srvType <= 10 || srvType >= 60 {
		return
	}

	if err := p.user.PostGateRaw(uint8(srvType), rawMsg.Msg); err != nil {
		p.user.Error("MsgForward failed ", err)
	}
}

func (proc *GateUserMsgProc) MsgProc_TestBinMsg(content msgdef.IMsg) {
	proc.user.GetClientSess().Send(content)
}

// MsgProc_CellEntityMsg 客户端发送给服务器cell场景的消息都通过该消息包装
func (p *GateUserMsgProc) MsgProc_CellEntityMsg(content msgdef.IMsg) {
	msg := content.(*msgdef.CellEntityMsg)
	if err := p.user.PostGateRawToCell(msg.Data); err != nil {
		p.user.Error(err)
	} else {
		// p.user.Debug("MsgProc_CellEntityMsg")
	}
}
