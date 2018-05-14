package main

import (
	"zeus/dbservice"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
)

// sess 登录
func (mgr *GateUserMgr) login(uid uint64) {
	sess := GetSrvInst().clientSrv.GetSession(uid)
	if sess == nil {
		log.Error("add user but sess is not existed", uid)
		return
	}

	//此处为断线重连
	user := mgr.getUserByUID(uid)
	if user != nil && dbservice.SessionUtil(uid).VerifyToken(user.token) {
		user.cancelLogoutTick()
		user.SetClient(sess)
		sess.Send(&msgdef.ClientVertifySucceedRet{
			Source:   msgdef.ClientMSG,
			UID:      user.GetDBID(),
			SourceID: GetSrvInst().GetSrvID(),
			Type:     1,
		})
		user.saveSessionInfo(sess.RemoteAddr())

		sess.FlushBacklog()

		log.Info("user reconnected ", uid)
	} else {
		id, err := dbservice.SessionUtil(uid).GetUserEntityID()
		if err == nil {
			GetSrvInst().DestroyEntityAll(id)
		}
		user = mgr.addUser(uid)
		if user == nil {
			sess.Close()
			return
		}
	}
}

// sess 下线 , 启动倒计时
func (mgr *GateUserMgr) logout(uid uint64) {
	user := mgr.getUserByUID(uid)
	if user == nil || user.isLogoutTick() {
		return
	}

	user.startLogoutTick()
}

// 真正下线
func (mgr *GateUserMgr) _logout(uid uint64) {
	user := mgr.getUserByUID(uid)
	if user == nil {
		return
	}

	GetSrvInst().DestroyEntityAll(user.GetID())
}
