package sess

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"zeus/dbservice"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/msghandler"

	log "github.com/cihub/seelog"
)

//IConnNotify 连接状态发生改变得到通知
type IConnNotify interface {
	connErr(uint64)
}

func newSessMgr(msghandlers msghandler.IMsgHandlers) *SessMgr {
	return &SessMgr{
		msgHandlers:     msghandlers,
		sessesByAddr:    &sync.Map{},
		sessesByID:      &sync.Map{},
		oldSessesByID:   &sync.Map{},
		_encryptEnabled: false,
	}
}

// RawMsg 对象
type RawMsg struct {
	MsgID int
	Msg   []byte
}

// SessMgr 管理 Sess
type SessMgr struct {
	msgHandlers msghandler.IMsgHandlers

	sessesByAddr  *sync.Map
	sessesByID    *sync.Map
	oldSessesByID *sync.Map

	_encryptEnabled bool
}

// SetEncryptEnabled 设置加密模式
func (mgr *SessMgr) SetEncryptEnabled() {
	mgr._encryptEnabled = true
}

func (mgr *SessMgr) acceptConn(conn net.Conn) {

	msg, _, _, err := readARQMsgForward(conn)

	//如果是连接扫描就强制关掉
	if err != nil || msg == nil {
		//log.Error("accept conn and  read first message error ", err, conn.RemoteAddr())
		conn.Close()
		return
	}

	if msg.Name() != "ClientVertifyReq" {
		log.Error("read first message , but message is not client vertify req ", conn.RemoteAddr())
		conn.Close()
		return
	}

	m := msg.(*msgdef.ClientVertifyReq)

	sess := mgr.newSess(conn)
	if err := mgr.vertifySess(m, sess); err != nil {
		sess.Close()
		log.Info("client vertify failed", conn.RemoteAddr(), err)
		return
	}

	mgr.putSess(sess)
	sess.SetServerType(m.Source)

	sess.DoNormalMsg("SessVertified", msg.(*msgdef.ClientVertifyReq).UID)
	sess.Touch()
	sess.Start()
}

func (mgr *SessMgr) vertifySess(msgContent msgdef.IMsg, sess iserver.ISess) error {

	if sess.IsVertified() {
		return errors.New("sess had vertified but receive a clientvetifyreq message ")
	}

	msg := msgContent.(*msgdef.ClientVertifyReq)
	var verify = false
	// 验证
	if msg.Source != msgdef.ClientMSG {
		// 服务器间连接暂时不验证 fixme
		verify = true
	} else {
		verify = dbservice.SessionUtil(msg.UID).VerifyToken(msg.Token)
	}

	if !verify {
		return errors.New(fmt.Sprintln("sess vertified failed ", msg.UID))
	}

	sess.SetVertify()
	sess.SetID(msg.UID)

	return nil
}

func (mgr *SessMgr) close() {
	//此处会不会死锁？
	mgr.sessesByID.Range(func(_, sess interface{}) bool {
		sess.(iserver.ISess).Close()
		return true
	})

	mgr.sessesByAddr = &sync.Map{}
	mgr.sessesByID = &sync.Map{}
}

func (mgr *SessMgr) newSess(conn net.Conn) iserver.ISess {
	return newSess(conn, mgr, mgr._encryptEnabled)
}

func (mgr *SessMgr) putSess(sess iserver.ISess) {

	_, ok := mgr.sessesByID.Load(sess.GetID())
	if ok {
		mgr.removeSess(sess.GetID())
	}

	sess.SetMsgHandler(mgr.msgHandlers)

	mgr.sessesByID.Store(sess.GetID(), sess)
	mgr.sessesByAddr.Store(sess.RemoteAddr(), sess)

	log.Debug("putSess ", sess.RemoteAddr(), "ID: ", sess.GetID())

	if os, ok := mgr.oldSessesByID.Load(sess.GetID()); ok {
		sess.FetchBacklog(os.(iserver.ISess))
		mgr.oldSessesByID.Delete(sess.GetID())
	}
}

func (mgr *SessMgr) removeSess(id uint64) {

	if id == 0 {
		return
	}

	ios, ok := mgr.sessesByID.Load(id)
	if !ok {
		// log.Error("remove sess but sess id is not in list ", id)
		return
	}
	mgr.sessesByID.Delete(id)
	mgr.sessesByAddr.Delete(ios.(iserver.ISess).RemoteAddr())
	mgr.oldSessesByID.Store(id, ios)

	log.Debug("remove sess ", ios.(iserver.ISess).RemoteAddr())

	time.AfterFunc(90*time.Second, func() {
		mgr.oldSessesByID.Delete(id)
	})

	ios.(iserver.ISess).Close()
}

func (mgr *SessMgr) getSessByID(id uint64) iserver.ISess {

	is, ok := mgr.sessesByID.Load(id)
	if !ok {
		return nil
	}

	return is.(iserver.ISess)
}

func (mgr *SessMgr) connErr(id uint64) {
	mgr.removeSess(id)
}

// Count 获取Sess的数量
func (mgr *SessMgr) Count() int {
	return 0
}
