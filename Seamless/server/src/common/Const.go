package common

//玩家在线状态
const (
	//Online_off 离线
	Online_off = 0
	//Online_on 在线
	Online_on = 1
	//Online_team 组队中
	Online_team = 2
	//Online_game 游戏中
	Online_game = 3
)

const (
	//SpaceMaxTime 一局游戏的最大时间
	SpaceMaxTime = 3000
	//SpaceWaitBalance 等待全部结算请求的时间
	SpaceWaitBalance = 2
	//QueueLimit 在队列中等待匹配的上限
	QueueLimit    = 30
	SelectTime    = 3
	MatchWaitTime = 3
	RoomUserLimit = 20
	PlaneHeight   = 7.0
)

const (
	//SpaceStatusInit 初始状态
	SpaceStatusInit int = iota + 1
	//SpaceStatusSelect 玩家选角
	SpaceStatusSelect
	//SpaceStatusCreating room游戏初始化
	SpaceStatusCreating
	SpaceStatusInitRoom
	//SpaceStatusWaitIn 玩家进入
	SpaceStatusWaitIn
	//SpaceStatusBegin 游戏开始
	SpaceStatusBegin
	//SpaceStatusBalance 游戏结束
	SpaceStatusBalance
	//SpaceStatusVerify 验证结果中
	SpaceStatusVerify
	// SpaceStatusBalanceDone 结算完成
	SpaceStatusBalanceDone
	//SpaceStatusClose 游戏关闭
	SpaceStatusClose
)

const (
	System_TeamMemberNum      = 1
	System_InitCell           = 3
	System_RefreshMax         = 4
	System_InitRadius         = 5
	System_RoomUserLimit      = 6
	System_MatchWait          = 7
	System_RoomUserMin        = 8
	System_RefreshBomb        = 11
	System_RefreshSafeArea    = 12
	System_SafeAreaNotify     = 13
	System_SafeAreaRadiusRate = 14
	System_CrouchToDown       = 23
	System_CrouchToStand      = 24
	System_EnergyMax          = 25
	System_FallDamageA        = 26
	System_FallDamageB        = 27
	System_FallDamage         = 28
	System_MailOverdue        = 29
	System_MailMax            = 30
	System_VehicleDamageA     = 35
	System_VehicleDamageB     = 36
	System_VehicleHitA        = 42
	System_VehicleHitB        = 43
	System_AddHpLimit         = 44
	System_VehicleCollisionA  = 45
	System_VehicleCollisionB  = 46
	System_VehicleCollisionC  = 47
	System_FallHeight         = 49
	System_MatchMinWait       = 55
	System_InitBagPack        = 57
	System_ForceEjectNotify   = 58
	System_MinLoadingTime     = 59
	System_MaxLoadingTime     = 60
	System_SummonAiSum        = 63

	System_RefreshBoxHeight   = 65
	System_GetBoxCold         = 66
	System_RefreshBoxSpeed    = 67
	System_RefreshBoxID       = 68
	System_RefreshFakeBoxID   = 69
	System_KillGetCoin        = 70
	System_DayCoinLimit       = 71
	System_VehicleBombRadius  = 72
	System_VehicleBombDamageB = 73
	System_VehicleBombDamageA = 74
	System_VehicleFireTime    = 75
	System_VehicleFireRadius  = 76
	System_VehicleFireDamage  = 77
	System_VehicleAfterStop   = 85
)

const (
	Status_Init        = 0
	Status_ShrinkBegin = 1
)

const (
	Bomb_Status_Init        = 0
	Bomb_Status_ShrinkBegin = 1
	Bomb_Status_Disappear   = 2
)

const (
	Terrian_Area_Water = 20
)

const (
	Terrain_Area_WalkableMax = 19
)

const (
	// notifyCommon 普通通知
	NotifyCommon = 1
	// notifyeError 错误通知
	NotifyError = 2
)

const (

	// StateFree 空闲
	StateFree = 0
	// StateMatchWaiting 匹配等待中
	StateMatchWaiting = 1
	// StateMatching 匹配中
	StateMatching = 2
	// StateGame 游戏中
	StateGame = 3
	// StateOffline 不在线
	StateOffline = 4
)

const (
	//ErrCodeIDInviteJoin 邀请加入房间失效
	ErrCodeIDInviteJoin = 9
	//ErrCodeIDPackCell 背包空间不足
	ErrCodeIDPackCell = 11
)

const (
	// TeamMgr 组队管理器
	TeamMgr = "TeamMgr"
)

// 玩家所在队伍状态
const (
	// PlayerTeamMatching 玩家所在队伍匹配中
	PlayerTeamMatching = 1
)

// QQAppID ID
const QQAppID = 1106393072

const QQAppIDStr = "1106393072"

// WXAppID ID
const WXAppID = "wxa916d09c4b4ef98f"

// GAppID 游客登录时
const GAppID = "G_1106393072"

// MSDKKey KEY
const MSDKKey = "7eead96a3fdb063615b181d7c01480e4"

type ClothesPos int // 服装部位
// 服装部位定义。须与客户端定义一致。
const (
	// 头部服装
	ClothesPos_Head ClothesPos = 1
	// 面部
	ClothesPos_Face ClothesPos = 2
	// 上衣
	ClothesPos_Tops ClothesPos = 3
	// 裤子
	ClothesPos_Pant ClothesPos = 4
	// 鞋子
	ClothesPos_Shoes ClothesPos = 5
)

//cell里人数有多少开始拆分
const SplitEntityNum = 100

//cell里有多少人数开始合并
const MergeEntityNum = SplitEntityNum * 0.1

//被合并时，cell的压力要小于这个值
const BeMergeEntityNum = SplitEntityNum * 0.6

//cell切割百分比
const SplitPercent = 0.1
