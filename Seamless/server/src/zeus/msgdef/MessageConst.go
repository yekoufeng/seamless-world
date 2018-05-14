package msgdef

const (
	//ClientVertifyReqMsgID 验证消息的ID号
	ClientVertifyReqMsgID = 1
	//ClientVertifySucceedRetMsgID 验证结果成功返回的ID号
	ClientVertifySucceedRetMsgID = 2
	//ClientVertifyFailedRetMsgID 验证结果失败返回的ID号
	ClientVertifyFailedRetMsgID = 3

	// HeartBeatMsgID 心跳消息ID
	HeartBeatMsgID = 4

	// SessionNotifyID 会话状态消息ID
	//SessionNotifyID = 5
	//SessionCreateID 会话建立时的ID号
	//SessionCreateID = 6
	//ClientVertifySucceedMsgID 验证结果成功的ID号
	//ClientVertifySucceedMsgID = 7

	// TestBinMsgID 测试消息
	TestBinMsgID = 8

	//EntitySrvInfoNotifyID 实体变化消息
	EntitySrvInfoNotifyID = 9

	//ProtoSyncMsgID ProtoSync消息ID
	ProtoSyncMsgID = 10

	// ClientFrameMsgID 客户端过来的帧消息
	ClientFrameMsgID = 20
	// ServerFrameMsgID 服务器下发的单帧消息
	ServerFrameMsgID = 21
	// FramesMsgID 服务器下发的多帧消息
	FramesMsgID = 22
	// RequireFramesMsgID 客户端
	RequireFramesMsgID = 23
	// 	UserDuplicateLoginNotifyMsgID 玩家重复登录的消息
	UserDuplicateLoginNotifyMsgID = 24

	// PropsSyncMsgID 属性消息ID
	PropsSyncMsgID = 30
	// PropsSyncClientMsgID 发给客户端的属性消息ID
	PropsSyncClientMsgID = 31
	// MRolePropsSyncClientMsgID 主角属性同步消息DI
	MRolePropsSyncClientMsgID = 32

	// EnterCellReqMsgID 进入空间的消息
	EnterCellReqMsgID = 40
	// LeaveCellReqMsgID 离开空间的消息
	LeaveCellReqMsgID = 41
	// EnterAOIMsgID 进入AOI消息
	EnterAOIMsgID = 42
	// LeaveAOIMsgID 离开AOI消息
	LeaveAOIMsgID = 43
	// AOIPosChangeMsgID AOI位置改变消息
	AOIPosChangeMsgID = 44
	// EnterCellMsgID 玩家进入场景
	EnterCellMsgID = 45
	// LeaveCellMsgID 玩家离开场景
	LeaveCellMsgID = 46
	// UserMoveMsgID 玩家移动
	UserMoveMsgID = 47
	// CellEntityMsgID 空间实体消息
	CellEntityMsgID = 48
	// EntityPosSetMsgID 玩家移动
	EntityPosSetMsgID = 49

	// ClientTransportMsgID 客户端中转消息ID
	ClientTransportMsgID = 50
	// CreateEntityReqMsgID 请求创建实体消息ID
	CreateEntityReqMsgID = 51
	// CreateEntityRetMsgID 创建实体返回消息ID
	CreateEntityRetMsgID = 52
	// DestroyEntityReqMsgID 请求删除实体消息ID
	DestroyEntityReqMsgID = 53
	// DestroyEntityRetMsgID 销毁实体返回消息ID
	DestroyEntityRetMsgID = 54
	// EntityMsgTransportMsgID 分布式实体之间传递消息用
	EntityMsgTransportMsgID = 55
	// EntityMsgChangeMsgID 分布式实体之间同步数据使用
	EntityMsgChangeMsgID = 56
	// SrvMsgTransportMsgID 服务器间消息转发ID
	SrvMsgTransportMsgID = 57
	// RPCMsgID RPC消息
	RPCMsgID       = 58
	SyncClockMsgID = 59

	SpaceUserConnectMsgID           = 61
	SpaceUserConnectSucceedRetMsgID = 62
	SyncUserStateMsgID              = 63
	AOISyncUserStateMsgID           = 64
	AdjustUserStateMsgID            = 65
	EntityAOISMsgID                 = 66
	EntityBasePropsMsgID            = 67
	EntityEventMsgID                = 68

	// CellMsgTransportMsgID 发消息给某个cell
	CellMsgTransportMsgID = 71

	//创建ghost
	CreateGhostReqMsgID = 72
	//删除ghost
	DeleteGhostReqMsgID = 73
	//切换real
	TransferRealReqMsgID = 74
	//新real通知
	NewRealNotifyMsgID = 75
)

const (
	// ClientMSG 来自客户端的验证消息
	ClientMSG uint8 = 0
)

// ClientVertifySucceedRet 中登录类型
const (
	// Connected 正常连接
	Connected uint8 = 1
	// ReConnect 断线重连
	ReConnect uint8 = 2
	// DupConnect 重复连接
	DupConnect uint8 = 3
)

const (
	// SessHBTimeout 心跳超时
	SessHBTimeout uint8 = 1
	// SessDisconnect 断线
	SessDisconnect uint8 = 2
	// SessError 出错
	SessError uint8 = 3
)
