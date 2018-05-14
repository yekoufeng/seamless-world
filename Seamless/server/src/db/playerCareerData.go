package db

import (
	"fmt"
	"math"
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"
)

const (
	playerCareerDataPrefix = "PlayerCareerData"
)

type playerCareerDataUtil struct {
	uid    uint64
	season int
}

func PlayerCareerDataUtil(uid uint64, season int) *playerCareerDataUtil {
	return &playerCareerDataUtil{
		uid:    uid,
		season: season,
	}
}

// GetRoundData 获取对局数据
func (u *playerCareerDataUtil) GetRoundData(data interface{}) error {
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
func (u *playerCareerDataUtil) GetValueField(field string) (uint32, error) {
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
func (u *playerCareerDataUtil) GetValues(args []interface{}) ([]interface{}, error) {
	c := dbservice.Get()
	defer c.Close()

	return redis.Values(c.Do("HMGET", append([]interface{}{u.key()}, args...)...))
}

// SetRoundData 设置多个values
func (u *playerCareerDataUtil) SetRoundData(args interface{}) error {
	c := dbservice.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(u.key()).AddFlat(args)...)
	return err
}

// SetValueField 设置单个value属性
func (u *playerCareerDataUtil) SetValueField(field string, value string) error {
	c := dbservice.Get()
	defer c.Close()

	_, err := c.Do("HSET", u.key(), field, value)
	return err
}

func (u *playerCareerDataUtil) getValue(field string) (interface{}, error) {
	c := dbservice.Get()
	defer c.Close()
	return c.Do("HGET", u.key(), field)
}

func (u *playerCareerDataUtil) setValue(field string, value interface{}) error {
	c := dbservice.Get()
	defer c.Close()
	_, err := c.Do("HSET", u.key(), field, value)
	return err
}

func (u *playerCareerDataUtil) key() string {
	return fmt.Sprintf("%s:%d:%d", playerCareerDataPrefix, u.season, u.uid)
}

// FetchGameID 获取mysql索引唯一id
func FetchGameID() uint64 {
	c := dbservice.Get()
	defer c.Close()

	id, err := redis.Uint64(c.Do("INCRBY", "gameid", 1))
	if err != nil {
		return math.MaxUint64
	}
	return id
}
