package db

import (
	"fmt"
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"
)

/* 记录玩家积分。
分单人赛，双人赛和四人赛。分赛季。记录3个积分：
WinRating, KillRating, Rating
采用增量方式改写，因为比赛的开始与结算会交错，可能有多个结算并发。
而 PlayerCareerData 中的 Rating 是此处的一个快照，采用改写的方式更新。
*/

const (
	playerRatingPrefix = "PlayerRating"
)

type playerRatingUtil struct {
	uid    uint64
	season int
}

// PlayerRating 录入redis哈希表 PlayerRating 表中
type PlayerRating struct {
	Inited           bool // 是否已初始化
	SoloWinRating    float32
	SoloKillRating   float32
	SoloFinalRating  float32
	DuoWinRating     float32
	DuoKillRating    float32
	DuoFinalRating   float32
	SquadWinRating   float32
	SquadKillRating  float32
	SquadFinalRating float32
}

func PlayerRatingUtil(uid uint64, season int) *playerRatingUtil {
	return &playerRatingUtil{
		uid:    uid,
		season: season,
	}
}

// Get 获取数据
func (u *playerRatingUtil) Get(data *PlayerRating) error {
	c := dbservice.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", u.key()))
	if err != nil || len(values) == 0 {
		return nil
	}

	err = redis.ScanStruct(values, data)
	return err
}

// SetRoundData 设置多个values
func (u *playerRatingUtil) Set(data *PlayerRating) error {
	c := dbservice.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(u.key()).AddFlat(data)...)
	return err
}

// IncrFieldByFloat 单个字段加上浮点数增量
func (u *playerRatingUtil) IncrFieldByFloat(field string, increment float64) error {
	c := dbservice.Get()
	defer c.Close()

	_, err := c.Do("HINCRBYFLOAT", u.key(), field, increment)
	return err
}

func (u *playerRatingUtil) key() string {
	return fmt.Sprintf("%s:%d:%d", playerRatingPrefix, u.season, u.uid)
}
