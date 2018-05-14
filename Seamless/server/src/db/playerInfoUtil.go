package db

import (
	"fmt"
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"
)

const (
	playerInfoPrefix = "PlayerInfo"
)

const (
	fieldSadness = "Sadness"
)

//========================================
//playerInfoUtil管理玩家的各种信息
//========================================

type playerInfoUtil struct {
	uid uint64
}

//PlayerInfoUtil 生成用于管理玩家信息的工具
func PlayerInfoUtil(uid uint64) *playerInfoUtil {
	return &playerInfoUtil{
		uid: uid,
	}
}

// SetRegisterTime 设置注册时间, 设置成功返回true, 失败返回false
func (u *playerInfoUtil) SetRegisterTime(regTime int64) (bool, error) {
	c := dbservice.Get()
	defer c.Close()

	reply, err := c.Do("HSETNX", u.key(), "registertime", regTime)
	if err != nil {
		return false, err
	}
	v, err := redis.Int(reply, nil)
	if err != nil {
		return false, err
	}
	if v == 1 {
		return true, nil
	}
	return false, err
}

// GetRegisterTime 获取注册时间
func (u *playerInfoUtil) GetRegisterTime() (int64, error) {
	regTime, err := u.getValue("registertime")
	if err != nil {
		return 0, err
	}

	return redis.Int64(regTime, nil)
}

func (u *playerInfoUtil) key() string {
	return fmt.Sprintf("%s:%d", playerInfoPrefix, u.uid)
}

func (u *playerInfoUtil) getValue(field string) (interface{}, error) {
	c := dbservice.Get()
	defer c.Close()
	return c.Do("HGET", u.key(), field)
}

func (u *playerInfoUtil) setValue(field string, value interface{}) error {
	c := dbservice.Get()
	defer c.Close()
	_, err := c.Do("HSET", u.key(), field, value)
	return err
}

func (u *playerInfoUtil) GetLoginAward(str string) (int64, error) {
	if !hExists(u.key(), str) {
		return 0, nil
	}

	regTime, err := u.getValue(str)
	if err != nil {
		return 0, err
	}

	return redis.Int64(regTime, nil)
}

// GetSadness 获取玩家比赛失败而累积的沮丧值. 出错则返回0.
func (u *playerInfoUtil) GetSadness() float32 {
	sadness, err := u.getValue(fieldSadness)
	result, _ := redis.Float64(sadness, err)
	return float32(result)
}

// IncSadness 在玩家比赛结束后累积沮丧值.
func (u *playerInfoUtil) IncSadness(sadness float32) {
	c := dbservice.Get()
	defer c.Close()
	c.Do("HINCRBYFLOAT", u.key(), fieldSadness, sadness)
}

// ResetSadness 在玩家进入前几名，或进入安慰局后清空沮丧值.
func (u *playerInfoUtil) ResetSadness() {
	c := dbservice.Get()
	defer c.Close()
	c.Do("HDEL", u.key(), fieldSadness)
}
