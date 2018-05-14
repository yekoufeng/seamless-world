package main

const (
	// UserStatusInRoomGaming 玩家游戏中
	UserStatusInRoomGaming = 1
	// UserStatusInRoomQuit 玩家离线
	UserStatusInRoomQuit = 2
	// UserStatusInRoomAi 玩家在被托管状态
	UserStatusInRoomAi = 3
)

// 匹配容器状态
const (
	MatchModelInit = 0
	MatchModelDel  = 1
)

const (
	// TeamCreate 队伍组建
	TeamCreate = 1
	// TeamStart 队伍开始使用
	TeamStart = 2
	// TeamEnd 队伍结束
	TeamEnd = 3
)

//匹配类型
const (
	//MatchKind_Single 单人匹配
	MatchKindSingle = 1
	//MatchKind_Team 组队匹配
	MatchKindTeam = 2
)
