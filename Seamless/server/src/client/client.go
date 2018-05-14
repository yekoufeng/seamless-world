package main

import (
	"runtime/debug"
	"time"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/serializer"

	log "github.com/cihub/seelog"
)

var _client *Client

func GetClient() *Client {
	if _client == nil {
		_client = &Client{
			User: NewUser(),
			AOIS: NewAOIS(),
		}
		_client.IsMainRole = true
	}
	return _client
}

type Client struct {
	sess iserver.ISess

	UID       uint64
	Token     string
	LobbyAddr string

	*User
	*AOIS

	detectCellTime    int64
	cellOkConfirmTime int64

	lastEnterCellTime int64

	SelectEntity
	SkillMgr
}

func (cli *Client) GetUserName() string {
	if *username == "" {
		return "jcmiao"
	}
	return *username
}

func (cli *Client) Loop() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Error(string(debug.Stack()))
		}
	}()
loop1:
	for {
		select {
		case <-closeChan:
			log.Debug("choseChan")
			break loop1
		default:
		}
		cli.sess.DoMsg()

		cli.MainLoop()
		time.Sleep(33e6)
	}
}

func (cli *Client) MainLoop() {
	cli.AOISLoop()
	cli.DetectCell()
}

func (cli *Client) RPCCall(srvType uint8, methodName string, args ...interface{}) {
	if cli.sess == nil {
		log.Error("msgClient nil")
		return
	}

	data := serializer.Serialize(args...)

	msg := &msgdef.RPCMsg{}
	msg.ServerType = srvType
	// msg.SrcEntityID = srcEntityID
	msg.MethodName = methodName
	msg.Data = data

	cli.sess.Send(msg)
	return
}

func (cli *Client) GetSess() iserver.ISess {
	return cli.sess
}

func (cli *Client) SendMsgToCell(msg msgdef.IMsg) {
	msgLen := msg.Size()
	data := make([]byte, 4 /*msglen*/ +2 /*msgidsize*/ +msgLen)
	data[0] = byte(msgLen)
	data[1] = byte(msgLen >> 8)
	data[2] = byte(msgLen >> 16)
	data[3] = 0

	msgID, err := msgdef.GetMsgDef().GetMsgIDByName(msg.Name())
	if err != nil {
		log.Error(err)
		return
	}

	data[4] = (byte)(msgID)
	data[5] = (byte)(msgID >> 8)

	msg.MarshalTo(data[6:])
	cellMsg := &msgdef.CellEntityMsg{
		Data: data,
	}
	cli.sess.Send(cellMsg)
}
