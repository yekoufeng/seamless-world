package db

import (
	"fmt"
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"
)

const (
	playerDayDataPrefix = "PlayerDayData"
)

type playerDayDataUtil struct {
	uid    uint64
	season int
}

func PlayerDayDataUtil(uid uint64, season int) *playerDayDataUtil {
	return &playerDayDataUtil{
		uid:    uid,
		season: season,
	}
}

// GetRoundData 获取对局数据
func (u *playerDayDataUtil) GetDayData(data interface{}) error {
	c := dbservice.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", u.key()))
	if err != nil || len(values) == 0 {
		return nil
	}

	err = redis.ScanStruct(values, data)
	return err
}

// GetValueField 获取单个value
func (u *playerDayDataUtil) GetValueField(field string) (uint32, error) {
	reply, err := u.getValue(field)
	if err != nil || reply == nil {
		return 0, err
	}

	v, err := redis.Uint64(reply, nil)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

// GetValues 获取多个values
func (u *playerDayDataUtil) GetValues(args []interface{}) ([]interface{}, error) {
	c := dbservice.Get()
	defer c.Close()

	return redis.Values(c.Do("HMGET", append([]interface{}{u.key()}, args...)...))
}

// SetRoundData 设置多个values
func (u *playerDayDataUtil) SetDayData(args interface{}) error {
	c := dbservice.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(u.key()).AddFlat(args)...)
	return err
}

func (u *playerDayDataUtil) getValue(field string) (interface{}, error) {
	c := dbservice.Get()
	defer c.Close()
	return c.Do("HGET", u.key(), field)
}

func (u *playerDayDataUtil) setValue(field string, value interface{}) error {
	c := dbservice.Get()
	defer c.Close()
	_, err := c.Do("HSET", u.key(), field, value)
	return err
}

func (u *playerDayDataUtil) key() string {
	return fmt.Sprintf("%s:%d:%d", playerDayDataPrefix, u.season, u.uid)
}
