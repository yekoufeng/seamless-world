package main

import (
	"common"
	"db"
	"entitydef"
	"fmt"
	"msdk"
	"protoMsg"
	"time"
	"zeus/dbservice"
	"zeus/entity"
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
	"zeus/serializer"

	log "github.com/cihub/seelog"
)

// LobbyUser 大厅玩家
type LobbyUser struct {
	entitydef.PlayerDef
	entity.Entity
	tempCellID uint64 //所属cellID
	macthTime  int64  //unix时间戳
	//teamid  uint64 // 组队id

	// 压测数据相关
	enterRoomStamp time.Time

	friendMgr *FriendMgr
	gm        *GmMgr
	storeMgr  *StoreMgr

	loginMsg  *protoMsg.PlayerLogin
	loginTime int64
	isReg     bool

	onlineReportTicker *time.Ticker
	lastBatchTime      int64
	token              string
}

// Init 初始化调用
func (user *LobbyUser) Init(initParam interface{}) {
	user.RegMsgProc(&LobbyUserMsgProc{user: user})
	user.friendMgr = NewFriendMgr(user)
	user.storeMgr = NewStoreMgr(user)
	user.gm = NewGmMgr(user)
	user.createInit()

	user.loginTime = time.Now().Unix()
	user.SetLoginTime(user.loginTime)

	var err error
	user.isReg, err = db.PlayerInfoUtil(user.GetDBID()).SetRegisterTime(user.loginTime)
	if err != nil {
		log.Error(err, user)
	}

	tmpToken, errToken := dbservice.SessionUtil(user.GetDBID()).GetToken()
	if errToken != nil {
		log.Error(errToken)
	} else {
		user.token = tmpToken
	}

	user.onlineReportTicker = time.NewTicker(5 * time.Second)

	user.RPC(iserver.ServerTypeClient, "SrvTime", time.Now().UnixNano())

	// user.GetCellInfo()
	log.Info("LobbyUser Inited ", user)
}

//新玩家的数据初始化
func (user *LobbyUser) createInit() {
	if user.GetRoleModel() == 0 {
		user.SetRoleModel(uint32(common.GetTBSystemValue(32)))
	}

	//发送用户主数据
	maindata := &protoMsg.UserMainDataNotify{}
	maindata.Uid = user.GetID()
	maindata.Name_ = user.GetName()

	if err := user.Post(iserver.ServerTypeClient, maindata); err != nil {
		log.Error(err, user)
		return
	}
	user.SetPlayerGameState(common.StateFree)

	// 同步好友列表和好友申请列表
	user.friendMgr.syncFriendList()
	user.friendMgr.syncApplyList()

	checkMail(user.GetDBID())
	user.checkGlobalMail()
	user.MailNotify()

	// 初始化购买物品信息
	user.storeMgr.initPropInfo()

	user.InitAnnounceData()
	dataCenterOuterAddr, err := db.GetDataCenterAddr("DataCenterOuterAddr")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("dataCenterOuterAddr:", dataCenterOuterAddr)
	if err := user.RPC(iserver.ServerTypeClient, "InitNotifyMysqlDbAddr", "http://"+dataCenterOuterAddr); err != nil {
		log.Error(err)
	}
}

// Destroy 析构时调用
func (user *LobbyUser) Destroy() {
	if user.tempCellID != 0 {
		user.clear()
	}

	//user.setGameTime(0)
	user.setOnline(common.Online_off)

	// 检查是否被顶号
	if dbservice.SessionUtil(user.GetDBID()).VerifyToken(user.token) {
		user.SetPlayerGameState(common.StateOffline)
	} else {
		if err := user.Post(iserver.ServerTypeClient, &msgdef.UserDuplicateLoginNotify{}); err != nil {
			log.Error(err, user)
		}
	}

	user.onlineReportTicker.Stop()

	log.Debug("Lobby User Destroy ", user)
}

func (user *LobbyUser) clear() {
	user.tempCellID = 0
}

//setOnline 设置在线状态
func (user *LobbyUser) setOnline(s uint64) {
}

// AdviceNotify 提示信息通知
func (user *LobbyUser) AdviceNotify(notifyType uint32, id uint64) {
	if err := user.RPC(iserver.ServerTypeClient, "AdviceNotify", notifyType, id); err != nil {
		log.Error(err)
	}
}

// GetLoginChannel 1: 微信  2: QQ 3: 游客
func (user *LobbyUser) GetLoginChannel() uint32 {
	return user.GetPlayerLogin().LoginChannel
}

// 玩家注册初始化数据
func (user *LobbyUser) onUserRegister() {
	if user.loginMsg != nil && user.loginMsg.LoginChannel == 2 {
		var lst []*msdk.Param
		if user.isReg {
			//qqscorebatch 用户注册时间
			lst = append(lst, &msdk.Param{
				Tp:      25,
				BCover:  1,
				Data:    fmt.Sprintf("%v", user.loginTime),
				Expires: "不过期",
			})
			// msdk.QQScoreBatch(common.QQAppIDStr, common.MSDKKey, openid, user.GetAccessToken(), loginMsg.PlatID,
			// 	user.GetName(), user.GetDBID(), 25, 1, fmt.Sprintf("%v", user.loginTime), "不过期")
		}

		//qqscorebatch 最近登录时间
		lst = append(lst, &msdk.Param{
			Tp:      8,
			BCover:  1,
			Data:    fmt.Sprintf("%v", user.loginTime),
			Expires: "不过期",
		})

		msdk.QQScoreBatchList(common.QQAppIDStr, common.MSDKKey, user.loginMsg.VOpenID, user.GetAccessToken(), user.loginMsg.PlatID,
			user.GetName(), user.GetDBID(), lst)
	}
}

func (user *LobbyUser) InitAnnounceData() {

	// 初始玩家公告数据
	data := db.GetAllAnnuoncingData()
	if data == nil {
		return
	}

	if err := user.RPC(iserver.ServerTypeClient, "InitAnnounceData", data); err != nil {
		log.Error(err)
	}

	log.Info("InitAnnounceData ", user.GetID(), data)

}

// SetPlayerGameState 设置玩家游戏状态
func (user *LobbyUser) SetPlayerGameState(state uint64) {
	user.friendMgr.syncFriendState(state)
	db.PlayerTempUtil(user.GetDBID()).SetGameState(state)
}

func (user *LobbyUser) Loop() {
	select {
	case <-user.onlineReportTicker.C:
		if user.loginMsg == nil || user.loginMsg.LoginChannel != 2 { //不是手Q不上报
			return
		}

		now := time.Now()
		nowUnix := now.Unix()
		if nowUnix < user.lastBatchTime+300 { //5分钟上报一次
			zeroUnixTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
			if nowUnix < zeroUnixTime+24*3600-5 { //上报的时间还没到，但是下次触发还是在今天就返回，不然今天就再触发一次
				return
			}
		}
		user.lastBatchTime = nowUnix

		// isOnline, err := dbservice.SessionUtil(user.GetDBID()).IsExisted()
		// if err != nil {
		// 	log.Error(err)
		// 	return
		// }

		// if !isOnline { //不在线了，不上报了
		// 	return
		// }

		todayOnlineTime := user.CalTodayOnlineTime(user.GetTodayOnlineTime(), user.loginTime, now)
		//qqscorebatch 累计在线时间
		msdk.QQScoreBatch(common.QQAppIDStr, common.MSDKKey, user.loginMsg.VOpenID, user.GetAccessToken(), user.loginMsg.PlatID,
			user.GetName(), user.GetDBID(), 6000, 1, fmt.Sprintf("%v", todayOnlineTime), "不过期")
	default:
	}
}

func (user *LobbyUser) CalTodayOnlineTime(todayOnlineTime, loginTime int64, now time.Time) int64 {
	zeroUnixTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	if zeroUnixTime > loginTime {
		return now.Unix() - zeroUnixTime
	} else {
		return todayOnlineTime + (now.Unix() - loginTime)
	}
}

//获取所在的cell
func (user *LobbyUser) GetCellInfo() {

	// if user.tempCellID != 0 {
	// 	log.Error("GetCellInfo 已经进入了地图")
	// 	return
	// }

	playerMapData := db.PlayerMapUtil(user.GetDBID()).GetPlayerMapData()
	if playerMapData == nil {
		playerMapData = &db.PlayerMapData{Pos: linmath.Vector3{50, 0, 50}}
	}

	log.Info("GetCellInfo 玩家初始位置, X: ", playerMapData.Pos.X, ", Y:", playerMapData.Pos.Y, ", Z", playerMapData.Pos.Z)

	msg := &protoMsg.CellInfoReq{
		EntityID: user.GetID(),
		MapName:  playerMapData.MapName,
		Pos:      &protoMsg.Vector3{playerMapData.Pos.X, playerMapData.Pos.Y, playerMapData.Pos.Z},
		SrvID:    iserver.GetSrvInst().GetSrvID(),
	}

	serverID, err := iserver.GetSrvInst().GetSrvIDBySrvType(common.ServerTypeCellAppMgr)
	if err != nil {
		log.Error(err)
		return
	}

	if err = iserver.GetSrvInst().PostMsgToCell(serverID, 0, msg); err != nil {
		log.Error(err)
	}
}

//EnterCell 进入cell
func (user *LobbyUser) EnterCell(cellInfo *protoMsg.CellInfoRet) {
	log.Info("EnterCell cellSrvID: ", cellInfo.GetCellSrvID(), ", cellID: ", cellInfo.GetCellID(), ", entityID: ", cellInfo.GetEntityID())

	//避免重复进入场景
	if user.GetCellID() != 0 {
		log.Error("EnterCell cellID: ", user.GetCellID())
		return
	}

	user.tempCellID = cellInfo.GetCellID()

	sendFunc := func() {
		msg := &msgdef.EnterCellReq{
			SrvID:      cellInfo.GetCellSrvID(),
			CellID:     cellInfo.GetCellID(),
			EntityType: "Player",
			EntityID:   user.GetID(),
			DBID:       user.GetDBID(),
			InitParam:  serializer.Serialize(user.GetInitParam()),
			OldSrvID:   0,
			OldCellID:  0,
			Pos:        linmath.Vector3{cellInfo.GetPos().GetX(), cellInfo.GetPos().GetY(), cellInfo.GetPos().GetZ()},
		}

		log.Debug("EnterSpaceReq:", msg)
		if err := iserver.GetSrvInst().PostMsgToCell(cellInfo.GetCellSrvID(), cellInfo.GetCellID(), msg); err != nil {
			user.Error("Enter cell failed: ", err)
		}
	}

	sendFunc()
}
