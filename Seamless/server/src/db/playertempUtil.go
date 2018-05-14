package db

import (
	"fmt"
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"
)

// 玩家在线情况下各种临时信息的存储, 不需要持久化

const (
	playerTempPrefix = "PlayerTemp"
)

type playerTempUtil struct {
	uid uint64
}

// PlayerTempUtil 获取工具类
func PlayerTempUtil(uid uint64) *playerTempUtil {
	return &playerTempUtil{
		uid: uid,
	}
}

// GetPlayerTeamID 获取玩家所在队伍id
func (u *playerTempUtil) GetPlayerTeamID() uint64 {
	reply, err := dbservice.CacheHGET(u.key(), "teamid")
	if err != nil {
		return 0
	}
	var teamid uint64
	teamid, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}
	return teamid
}

// SetPlayerTeamID 设置玩家所在队伍id
func (u *playerTempUtil) SetPlayerTeamID(teamID uint64) error {
	return dbservice.CacheHSET(u.key(), "teamid", teamID)
}

// GetPlayerJumpAir 获取玩家跳过飞机
func (u *playerTempUtil) GetPlayerJumpAir() uint64 {
	reply, err := dbservice.CacheHGET(u.key(), "jumpair")
	if err != nil {
		return 0
	}
	var jump uint64
	jump, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}
	return jump
}

// SetPlayerJumpAir 设置是否跳过飞机
func (u *playerTempUtil) SetPlayerJumpAir(value uint64) error {
	return dbservice.CacheHSET(u.key(), "jumpair", value)
}

// GetGameState 获取玩家游戏状态
func (u *playerTempUtil) GetGameState() uint64 {
	reply, err := dbservice.CacheHGET(u.key(), "gamestate")
	if err != nil {
		return 0
	}
	var state uint64
	state, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}
	return state
}

// SetGameState 设置玩家游戏状态
func (u *playerTempUtil) SetGameState(state uint64) error {
	return dbservice.CacheHSET(u.key(), "gamestate", state)
}

// GetEnterGameTime 获取玩家进入游戏时间
func (u *playerTempUtil) GetEnterGameTime() uint64 {
	reply, err := dbservice.CacheHGET(u.key(), "entertime")
	if err != nil {
		return 0
	}
	var entertime uint64
	entertime, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}
	return entertime
}

// playerTempUtil 设置玩家进入游戏时间
func (u *playerTempUtil) SetEnterGameTime(timestamp uint64) error {
	return dbservice.CacheHSET(u.key(), "entertime", timestamp)
}

func (u *playerTempUtil) key() string {
	return fmt.Sprintf("%s:%d", playerTempPrefix, u.uid)
}
