package main

import (
	"time"
	"zeus/dbservice"
	"zeus/entity"
	"zeus/msgdef"
)

// GateUser 网关玩家
type GateUser struct {
	entity.Entity
	token       string
	logoutTimer *time.Timer
}

// Init 初始化调用
func (user *GateUser) Init(initParam interface{}) {
	user.RegMsgProc(&GateUserMsgProc{user: user})

	var err error
	// 获取并保存token
	if user.token, err = dbservice.SessionUtil(user.GetDBID()).GetToken(); err != nil {
		user.Error("获取Token失败 ", err)
	}

	sess := GetSrvInst().clientSrv.GetSession(user.GetDBID())
	user.SetClient(sess)
	user.onLoginSuccessFinal()

	user.Info("GateUser Inited")
}

// Destroy 析构时调用
func (user *GateUser) Destroy() {
	sess := user.GetClientSess()
	// 清理redis中的临时信息
	err := dbservice.SessionUtil(user.GetDBID()).DelSession(user.token)
	if err == dbservice.TokenVerifyError {
		if sess != nil {
			sess.Send(&msgdef.UserDuplicateLoginNotify{})

			// 被顶号之后等一秒后关闭连接, 让这条消息能到达客户端
			time.Sleep(1 * time.Second)
		}
	} else if err != nil {
		user.Error("删除Session失败 ", err)
	}
	if sess != nil {
		sess.SetMsgHandler(sess)
		sess.Close()
	}

	user.Info("GateUser Destroy")
}

func (user *GateUser) isLogoutTick() bool {
	return user.logoutTimer != nil
}

func (user *GateUser) startLogoutTick() {
	if user.logoutTimer != nil {
		return
	}

	user.logoutTimer = time.AfterFunc(1*time.Minute, user.logout)

	user.Info("Start logout tick")
}

func (user *GateUser) cancelLogoutTick() {
	if user.logoutTimer == nil {
		return
	}

	user.logoutTimer.Stop()
	user.logoutTimer = nil

	user.Info("Cancel logout tick")
}

//退出登陆
func (user *GateUser) logout() {
	GetSrvInst().DestroyEntityAll(user.GetID())
}

// 登录完全成功之后的操作
// typ: 登录类型, 正常登陆/断线重连/重复登录
func (user *GateUser) onLoginSuccessFinal() {
	sess := user.GetClientSess()
	if sess == nil {
		user.startLogoutTick()
		return
	}

	sess.Send(&msgdef.ClientVertifySucceedRet{
		Source:   msgdef.ClientMSG,
		UID:      user.GetDBID(),
		SourceID: GetSrvInst().GetSrvID(),
		Type:     0,
	})
	sess.Send(GetSrvInst().protoSync)

	user.saveSessionInfo(sess.RemoteAddr())
	user.sendMRoleProp()
}

func (user *GateUser) saveSessionInfo(ip string) {
	info := dbservice.SessionInfo{
		Token:    user.token,
		EntityID: user.GetID(),
		IP:       ip,
	}
	if err := dbservice.SessionUtil(user.GetDBID()).SaveSessionInfo(&info); err != nil {
		user.Error("保存Session信息失败", err)
	}
}

//客户端拿到这条消息是创建主角，而不是更新属性
func (user *GateUser) sendMRoleProp() {
	num, data := user.PackMRoleProps()
	msg := &msgdef.MRolePropsSyncClient{}
	msg.EntityID = user.GetID()
	msg.Num = uint32(num)
	msg.Data = data

	if sess := user.GetClientSess(); sess != nil {
		sess.Send(msg)
	} else {
		user.Error("发送主角属性失败, Sess为空")
	}
}
