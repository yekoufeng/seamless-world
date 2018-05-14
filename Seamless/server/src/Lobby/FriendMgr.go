package main

import (
	"common"
	"db"
	"protoMsg"
	"strconv"
	"time"
	"zeus/dbservice"
	"zeus/entity"
	"zeus/iserver"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// FriendMgr 好友管理器
type FriendMgr struct {
	user       *LobbyUser
	platFrient map[uint64]string
}

// NewFriendMgr 获取好友管理器
func NewFriendMgr(user *LobbyUser) *FriendMgr {
	mgr := &FriendMgr{
		user:       user,
		platFrient: make(map[uint64]string),
	}

	return mgr
}

// addFriend 添加好友
func (mgr *FriendMgr) addFriend(id uint64) {

	if !db.GetFriendUtil(mgr.user.GetDBID()).InApplyListByID(id) {
		// 请求申请id无效
		mgr.user.AdviceNotify(common.NotifyCommon, 2)
		mgr.syncApplyList()
		return
	}

	if db.GetFriendUtil(mgr.user.GetDBID()).IsReachLimit() {
		// 你的好友已达上限
		mgr.user.AdviceNotify(common.NotifyCommon, 3)
		return
	}

	if db.GetFriendUtil(id).IsReachLimit() {
		// 对方好友已达上限
		mgr.user.AdviceNotify(common.NotifyCommon, 4)
		return
	}

	applyInfo := db.GetFriendUtil(mgr.user.GetDBID()).GetSigleApplyReq(id)
	if applyInfo == nil {
		return
	}

	// 自己添加好友
	db.GetFriendUtil(mgr.user.GetDBID()).AddFriend(db.FriendInfo{ID: id, Name: applyInfo.Name})
	mgr.user.friendMgr.syncFriendList()

	// 删除申请请求
	mgr.delApplyReq(id)

	// 目标玩家添加好友
	db.GetFriendUtil(id).AddFriend(db.FriendInfo{ID: mgr.user.GetDBID(), Name: mgr.user.GetName()})

	// 删除目标玩家请求列表中申请信息
	db.GetFriendUtil(id).DelApply(mgr.user.GetDBID())

	// 更新目标玩家好友列表和申请列表
	mgr.SendProxyInfo(id, "SyncFriendList")
	mgr.SendProxyInfo(id, "SyncApplyList")

}

// SendProxyInfo 发送代理请求
func (mgr *FriendMgr) SendProxyInfo(targetID uint64, event string) {
	entityID, err := dbservice.SessionUtil(targetID).GetUserEntityID()
	if err != nil {
		return
	}

	srvID, cellID, errSrv := dbservice.EntitySrvUtil(entityID).GetSrvInfo(iserver.ServerTypeGateway)
	if errSrv != nil {
		return
	}

	proxy := entity.NewEntityProxy(srvID, cellID, entityID)
	proxy.RPC(common.ServerTypeLobby, event)
}

//WXAddFriend 微信添加好友
func (mgr *FriendMgr) WXAddFriend(targetUser *LobbyUser) {

	if targetUser == nil {
		return
	}

	// 你的好友已达上限
	if db.GetFriendUtil(mgr.user.GetDBID()).IsReachLimit() {
		mgr.friendApplyReq(targetUser.GetName())
		return
	}

	// 对方好友已达上限
	if db.GetFriendUtil(targetUser.GetDBID()).IsReachLimit() {
		mgr.friendApplyReq(targetUser.GetName())
		return
	}

	// 自己添加好友
	db.GetFriendUtil(mgr.user.GetDBID()).AddFriend(db.FriendInfo{ID: targetUser.GetDBID(), Name: targetUser.GetName()})
	mgr.user.friendMgr.syncFriendList()
	// 删除申请请求
	mgr.delApplyReq(targetUser.GetDBID())

	// 目标玩家添加好友
	db.GetFriendUtil(targetUser.GetDBID()).AddFriend(db.FriendInfo{ID: mgr.user.GetDBID(), Name: mgr.user.GetName()})

	// 更新目标玩家好友列表
	targetUser.friendMgr.delApplyReq(mgr.user.GetDBID())
	targetUser.friendMgr.syncFriendList()
}

// delFriend 删除好友
func (mgr *FriendMgr) delFriend(id uint64) {

	// 删除自己
	if db.GetFriendUtil(mgr.user.GetDBID()).DelFriend(id) {
		mgr.syncFriendList()

		//  删除对方
		if db.GetFriendUtil(id).DelFriend(mgr.user.GetDBID()) {
			mgr.SendProxyInfo(id, "SyncFriendList")
		}
	} else {
		log.Warn("delFriend id fail id = ", id)
	}
}

// syncFriendList rpc同步好友信息
func (mgr *FriendMgr) syncFriendList() {

	retMsg := &protoMsg.SyncFriendList{}

	list := make([]*db.FriendInfo, 0)
	list = db.GetFriendUtil(mgr.user.GetDBID()).GetFriendList()

	for _, info := range list {

		args := []interface{}{
			"LogoutTime",
			"Picture",
			"QQVIP",
			"NickName",
			"GameEnter",
		}
		values, valueErr := dbservice.EntityUtil("Player", info.ID).GetValues(args)
		if valueErr != nil {
			continue
		}

		item := protoMsg.FriendInfo{
			Id:        info.ID,
			Name_:     info.Name,
			State:     0,
			Time:      0,
			Url:       "",
			Enterplat: "platform",
			Qqvip:     0,
			Nickname:  "nickname",
		}

		tmpUrl, urlErr := redis.String(values[1], nil)
		if urlErr == nil {
			item.Url = tmpUrl
		}
		tmpQQvip, qqvipErr := redis.Int64(values[2], nil)
		if qqvipErr == nil {
			item.Qqvip = uint32(tmpQQvip)
		}
		tmpNickname, nicknameErr := redis.String(values[3], nil)
		if nicknameErr == nil {
			item.Nickname = tmpNickname
		}
		tmpEnterplat, platformErr := redis.String(values[4], nil)
		if platformErr == nil {
			item.Enterplat = tmpEnterplat
		}

		// 根据是否存在Session表判断是否在线
		isOnline, err := dbservice.SessionUtil(info.ID).IsExisted()
		if err == nil && isOnline == false {
			item.State = common.StateOffline

			timeStr, timeTmp := redis.String(values[0], nil)
			if timeTmp == nil {
				tmpTime, erro := strconv.ParseInt(timeStr, 10, 64)
				if erro == nil {
					item.Time = uint32(tmpTime)
				}
			}

		} else if err == nil {
			item.State = uint32(db.PlayerTempUtil(info.ID).GetGameState())
			if item.State == common.StateGame {
				item.Time = uint32(db.PlayerTempUtil(info.ID).GetEnterGameTime())
			}
		}
		log.Info("friendinfo item = ", item)

		retMsg.Item = append(retMsg.Item, &item)
	}

	if err := mgr.user.RPC(iserver.ServerTypeClient, "SyncFriendList", retMsg); err != nil {
		log.Error(err)
	}

}

// addApplyReq 添加申请请求
func (mgr *FriendMgr) friendApplyReq(name string) {

	// 判断是否存在名字为name的玩家
	result := 0

	targetID := db.GetIDByName(name)
	if targetID == 0 || targetID == mgr.user.GetDBID() {
		result = 1
		mgr.user.AdviceNotify(common.NotifyCommon, 6)
	} else {
		// 判断对方是否已经是好友
		if db.GetFriendUtil(mgr.user.GetDBID()).IsFriendByID(targetID) {
			result = 2
			mgr.user.AdviceNotify(common.NotifyCommon, 7)
		}

		// 已在对方申请列表中
		if db.GetFriendUtil(targetID).InApplyListByID(mgr.user.GetDBID()) {
			result = 3
			mgr.user.AdviceNotify(common.NotifyCommon, 8)
		}

		// 判断是否到达申请列表上线
		if db.GetFriendUtil(targetID).IsReachApplyLimit() {
			result = 4
			mgr.user.AdviceNotify(common.NotifyCommon, 10)
		}

	}

	if result == 0 {
		mgr.user.AdviceNotify(common.NotifyCommon, 5)

		info := db.ApplyInfo{
			ID:        mgr.user.GetDBID(),
			Name:      mgr.user.GetName(),
			ApplyTime: time.Now().Unix(),
		}

		// 申请信息添加至数据库
		db.GetFriendUtil(targetID).AddApplyInfo(info)

		// 更新目标玩家申请列表
		mgr.SendProxyInfo(targetID, "SyncApplyList")

	}

	// 好友申请请求结果 0申请成功 1用户名不存在 2已是好友 3已经申请 4达到申请列表上限
	mgr.user.RPC(iserver.ServerTypeClient, "FriendApplyReqRet", uint32(result))

}

// delApplyReq 删除申请请求
func (mgr *FriendMgr) delApplyReq(id uint64) {

	if !db.GetFriendUtil(mgr.user.GetDBID()).DelApply(id) {
		log.Warn("delApplyReq lose efficacy")
	}

	mgr.syncApplyList()
}

// syncApplyList 同步申请列表
func (mgr *FriendMgr) syncApplyList() {

	retMsg := &protoMsg.SyncFriendApplyList{}

	list := make([]*db.ApplyInfo, 0)
	list = db.GetFriendUtil(mgr.user.GetDBID()).GetApplyList()

	for _, info := range list {

		item := protoMsg.FriendApplyInfo{
			Id:        info.ID,
			Name_:     info.Name,
			ApplyTime: info.ApplyTime,
		}

		retMsg.Item = append(retMsg.Item, &item)
	}

	if err := mgr.user.RPC(iserver.ServerTypeClient, "SyncApplyList", retMsg); err != nil {
		log.Error(err)
	}

}

// InitPlatFriendList 初始化好友平台好友
func (mgr *FriendMgr) InitPlatFriendList(msg *protoMsg.PlatFriendStateReq) {
	if msg == nil {
		return
	}

	retMsg := &protoMsg.PlatFriendStateRet{}
	platFrientID := make([]uint64, 0)

	for _, openid := range msg.Openid {

		log.Debug("info  openid:", openid)
		uid, err := dbservice.GetUID(openid)
		if err != nil || uid == 0 {
			continue
		}

		item := &protoMsg.PlatFriendState{
			Openid: openid,
			Uid:    uid,
			State:  common.StateOffline,
			Time:   0,
			Name_:  "",
		}

		args := []interface{}{
			"LogoutTime",
			"Name",
		}

		values, valueErr := dbservice.EntityUtil("Player", uid).GetValues(args)
		if valueErr != nil || len(values) != 2 {
			continue
		}

		// 游戏中名称
		nameStr, nameErr := redis.String(values[1], nil)
		if nameErr == nil {
			item.Name_ = nameStr
		}

		if item.Name_ == "" {
			continue
		}

		// 根据是否存在Session表判断是否在线
		isOnline, err := dbservice.SessionUtil(uid).IsExisted()
		if err == nil && isOnline == false {
			item.State = common.StateOffline

			// 登出时间
			timeStr, timeTmp := redis.String(values[0], nil)
			if timeTmp == nil {
				tmpTime, erro := strconv.ParseInt(timeStr, 10, 64)
				if erro == nil {
					item.Time = uint32(tmpTime)
				}
			}

		} else if err == nil {
			item.State = uint32(db.PlayerTempUtil(uid).GetGameState())

			// 游戏时间
			if item.State == common.StateGame {
				item.Time = uint32(db.PlayerTempUtil(uid).GetEnterGameTime())
			}
		}

		retMsg.Data = append(retMsg.Data, item)
		platFrientID = append(platFrientID, item.Uid)

		mgr.platFrient[uid] = openid

		//log.Debugf("item info: %+v", item)
	}

	// 平台好友数据存入数据库
	db.GetFriendUtil(mgr.user.GetDBID()).UpdatePlatFrientInfo(platFrientID)

	if err := mgr.user.RPC(iserver.ServerTypeClient, "PlatFriendStateRet", retMsg); err != nil {
		log.Error(err)
	}

	curState := db.PlayerTempUtil(mgr.user.GetDBID()).GetGameState()
	mgr.syncFriendState(curState)
	//log.Debug("初始化好友平台好友数量 reqSum:", len(msg.Openid), " sum:", len(mgr.platFrient))
}

// syncFriendList rpc同步好友信息
func (mgr *FriendMgr) syncFriendState(state uint64) {

	list := make([]*db.FriendInfo, 0)
	list = db.GetFriendUtil(mgr.user.GetDBID()).GetFriendList()

	for _, info := range list {
		entityID, err := dbservice.SessionUtil(info.ID).GetUserEntityID()
		if err != nil {
			continue
		}

		srvID, cellID, errSrv := dbservice.EntitySrvUtil(entityID).GetSrvInfo(iserver.ServerTypeGateway)
		if errSrv != nil {
			continue
		}

		proxy := entity.NewEntityProxy(srvID, cellID, entityID)

		proxy.RPC(common.ServerTypeLobby, "SyncFriendState", mgr.user.GetDBID(), state)
	}

	for id, _ := range mgr.platFrient {
		entityID, err := dbservice.SessionUtil(id).GetUserEntityID()
		if err != nil {
			continue
		}

		srvID, cellID, errSrv := dbservice.EntitySrvUtil(entityID).GetSrvInfo(iserver.ServerTypeGateway)
		if errSrv != nil {
			continue
		}

		proxy := entity.NewEntityProxy(srvID, cellID, entityID)

		proxy.RPC(common.ServerTypeLobby, "SyncFriendState", mgr.user.GetDBID(), state)
	}

}
