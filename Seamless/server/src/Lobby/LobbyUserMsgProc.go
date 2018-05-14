package main

import (
	"common"
	"db"
	"fmt"
	"msdk"
	"protoMsg"
	"strings"
	"time"
	"zeus/iserver"
	"zeus/tlog"

	log "github.com/cihub/seelog"

	"zeus/tsssdk"
)

// LobbyUserMsgProc LobbyUser的消息处理函数
type LobbyUserMsgProc struct {
	user *LobbyUser
}

// RPC_EnterCellReq C->S 客户端进入地图
func (proc *LobbyUserMsgProc) RPC_EnterCellReq() {

	proc.user.GetCellInfo()

	log.Debug("RPC_EnterCellReq ", proc.user)
}

// 设置角色名称
func (proc *LobbyUserMsgProc) RPC_SetName(name string) {
	if proc.user == nil {
		return
	}

	result := 0
	if ret, err := tsssdk.JudgeUserInputName(name); !ret {
		if err != nil {
			log.Warn(err.Error())
		}
		result = 2
	} else if !db.JudgeNameInUse(name) {
		proc.user.SetName(name)
		db.AddUsedName(name, proc.user.GetDBID())
		result = 1
	}

	//qqscorebatch 角色名称
	if proc.user.loginMsg != nil && proc.user.loginMsg.LoginChannel == 2 {
		msdk.QQScoreBatch(common.QQAppIDStr, common.MSDKKey, proc.user.loginMsg.VOpenID, proc.user.GetAccessToken(), proc.user.loginMsg.PlatID,
			name, proc.user.GetDBID(), 0, 0, "", "")
	}

	proc.user.RPC(iserver.ServerTypeClient, "SetNameResult", uint64(result)) // 0 失败(名字已被使用) 1 成功

	log.Infof("设置角色名称 userdbid(%d) userid(%d) name(%s) result(%d)", proc.user.GetDBID(), proc.user.GetID(), name, result)
}

// RPC_AddFriend rpc添加好友
func (proc *LobbyUserMsgProc) RPC_AddFriend(id uint64) {
	proc.user.friendMgr.addFriend(id)
}

// RPC_DelFriend rpc删除好友
func (proc *LobbyUserMsgProc) RPC_DelFriend(id uint64) {
	proc.user.friendMgr.delFriend(id)
}

// RPC_SyncFriendList 同步好友列表
func (proc *LobbyUserMsgProc) RPC_SyncFriendList() {
	proc.user.friendMgr.syncFriendList()
}

// RPC_FriendApplyReq rpc添加申请请求
func (proc *LobbyUserMsgProc) RPC_FriendApplyReq(name string) {
	proc.user.friendMgr.friendApplyReq(name)
}

// RPC_DelApplyReq rpc删除申请请求
func (proc *LobbyUserMsgProc) RPC_DelApplyReq(id uint64) {
	proc.user.friendMgr.delApplyReq(id)
}

// RPC_SyncApplyList rpc同步申请列表
func (proc *LobbyUserMsgProc) RPC_SyncApplyList() {
	proc.user.friendMgr.syncApplyList()
}

// RPC_SyncFriendState 同步好友状态
func (proc *LobbyUserMsgProc) RPC_SyncFriendState(friendID uint64, state uint64) {
	curtime := time.Now().Unix()
	proc.user.RPC(iserver.ServerTypeClient, "SyncFriendState", friendID, uint32(state), uint32(curtime))
	//log.Debug("SyncFriendState ", friendID, " state:", state)
}

// RPC_GmLobby Gm命令
func (p *LobbyUserMsgProc) RPC_GmLobby(paras string) {
	log.Info("调用gm命令", paras)
	p.user.gm.exec(paras)
}

func (proc *LobbyUserMsgProc) RPC_AddCoin(num uint32) {
	log.Debug(proc.user.GetDBID(), "增加金币", num)
	proc.user.SetCoin(proc.user.GetCoin() + uint64(num))
}

// RPC_PlayerLogin 玩家登录消息
func (proc *LobbyUserMsgProc) RPC_PlayerLogin(msg *protoMsg.PlayerLogin) {
	if msg.VGameAppID == "" {
		if msg.LoginChannel == 1 {
			msg.VGameAppID = "wxa916d09c4b4ef98f"
		} else if msg.LoginChannel == 2 {
			msg.VGameAppID = "1106393072"
		} else if msg.LoginChannel == 3 {
			msg.VGameAppID = "G_1106393072"
		}
	}

	log.Info("Set loginchannel, accesstoken, platid: ", msg.LoginChannel, proc.user.GetAccessToken(), msg.PlatID)

	if msg.VOpenID == "" {
		msg.VOpenID = "no OpenID"
	}
	if msg.ClientVersion == "" {
		msg.ClientVersion = "1.0.0"
	}
	msg.SystemHardware = strings.Replace(msg.SystemHardware, ",", "-", -1)
	msg.SystemHardware = strings.Replace(msg.SystemHardware, "|", "-", -1)
	if msg.TelecomOper == "" {
		msg.TelecomOper = "no op"
	}
	if msg.Network == "" {
		msg.Network = "no Network"
	}

	msg.GameSvrID = fmt.Sprintf("%d", GetSrvInst().GetSrvID())
	msg.DtEventTime = time.Now().Format("2006-01-02 15:04:05")
	msg.IZoneAreaID = 0
	msg.Level = proc.user.GetLevel()
	msg.PlayerFriendsNum = proc.user.GetFriendsNum()
	msg.VRoleID = fmt.Sprintf("%d", proc.user.GetDBID())
	msg.VRoleName = proc.user.GetName()

	if msg.RegChannel == "" {
		msg.RegChannel = "NoRegChannel"
	} else {
		msg.RegChannel = strings.Replace(msg.RegChannel, ",", "-", -1)
		msg.RegChannel = strings.Replace(msg.RegChannel, "|", "-", -1)
	}

	proc.user.SetPlayerLogin(msg)
	proc.user.SetPlayerLoginDirty()
	proc.user.loginMsg = msg

	if proc.user.isReg {
		regMsg := &protoMsg.PlayerRegister{}
		regMsg.GameSvrID = msg.GameSvrID
		regMsg.DtEventTime = msg.DtEventTime
		regMsg.VGameAppID = msg.VGameAppID
		regMsg.PlatID = msg.PlatID
		regMsg.IZoneAreaID = msg.IZoneAreaID
		regMsg.VOpenID = msg.VOpenID
		regMsg.TelecomOper = msg.TelecomOper
		regMsg.RegChannel = msg.RegChannel
		regMsg.LoginChannel = msg.LoginChannel
		tlog.Format(regMsg)

		proc.user.onUserRegister()
	}

	tlog.Format(msg)

	if err := GetSrvInst().loginCnt(msg.VGameAppID, int(msg.PlatID)); err != nil {
		log.Error(err, proc.user, msg)
	}
}

func (proc *LobbyUserMsgProc) RPC_CheckUserState() {
	state := db.PlayerTempUtil(proc.user.GetDBID()).GetGameState()
	proc.user.RPC(iserver.ServerTypeClient, "RetUserState", uint8(state))
}

// RPC_SyncFlatFriend 同步平台好友
func (proc *LobbyUserMsgProc) RPC_SyncFlatFriend(msg *protoMsg.PlatFriendStateReq) {
	if msg == nil {
		return
	}
	proc.user.friendMgr.InitPlatFriendList(msg)
}

// RPC_KickingPlayer 执行踢出角色下线请求
func (proc *LobbyUserMsgProc) RPC_KickingPlayer(banAccountReason string) {
	proc.user.RPC(iserver.ServerTypeClient, "KickingPlayerMsg", banAccountReason)

	proc.user.RPC(iserver.ServerTypeClient, "NoticeKicking")
	err := GetSrvInst().DestroyEntityAll(proc.user.GetID())
	if err != nil {
		log.Error("RPC_KickingPlayer:", err, proc.user)
	}
}

// RPC_AccessTokenChange 玩家的accessToken变化了
func (proc *LobbyUserMsgProc) RPC_AccessTokenChange(accessToken string) {
	proc.user.SetAccessToken(accessToken)
}
