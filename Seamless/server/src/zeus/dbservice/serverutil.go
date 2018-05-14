package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

/*
 server:* 表工具类
*/

const (
	serverPrefix = "server:"
	servers      = "servers"
)

var serverInfoHashFields = []interface{}{
	"serverid",
	"type",
	"outeraddr",
	"inneraddr",
	"console",
	"load",
	"token",
	"status",
}

type serverUtil struct {
	uid uint64
}

// ServerUtil 获取server表的工具类
func ServerUtil(uid uint64) *serverUtil {
	srvutil := &serverUtil{}
	srvutil.uid = uid
	return srvutil
}

// GetServerList 获得服务器列表
func GetServerList(list interface{}) error {
	c := GetSingletonRedis()
	defer c.Close()

	ids, err := redis.Ints(c.Do("SMEMBERS", servers))
	if err != nil {
		return err
	}

	var values []interface{}

	for _, id := range ids {
		key := fmt.Sprintf("%s%d", serverPrefix, id)
		args := []interface{}{key}
		args = append(args, serverInfoHashFields...)
		vs, err := redis.Values(c.Do("HMGET", args...))
		if err != nil {
			continue
		}

		values = append(values, vs...)
	}

	if len(values) == 0 {
		return nil
	}

	if err = redis.ScanSlice(values, list); err != nil {
		return err
	}

	return nil
}

// SetStatus 设置服务器状态
func (util *serverUtil) SetStatus(status int) error {
	return util.setValue("status", status)
}

// SetLoad 设置服务器负载
func (util *serverUtil) SetLoad(load int) error {
	return util.setValue("load", load)
}

// Load 获取服务器负载
func (util *serverUtil) Load() (int, error) {
	return redis.Int(util.getValue("load"))
}

// Register 注册服务器信息, 设置过期时间60秒
func (util *serverUtil) Register(value interface{}) error {
	c := GetSingletonRedis()
	defer c.Close()
	if _, err := c.Do("HMSET", redis.Args{}.Add(util.key()).AddFlat(value)...); err != nil {
		return err
	}
	if _, err := c.Do("SADD", servers, util.uid); err != nil {
		return err
	}
	if _, err := c.Do("EXPIRE", util.key(), 60); err != nil {
		return err
	}
	return nil
}

// Update 更新服务器信息, 设置过期时间60秒
func (util *serverUtil) Update(value interface{}) error {
	c := GetSingletonRedis()
	defer c.Close()
	if _, err := c.Do("HMSET", redis.Args{}.Add(util.key()).AddFlat(value)...); err != nil {
		return err
	}
	if _, err := c.Do("EXPIRE", util.key(), 60); err != nil {
		return err
	}
	return nil
}

// IsExist 当前服务器是否已经注册过
func (util *serverUtil) IsExist() bool {
	c := GetSingletonRedis()
	defer c.Close()

	if isExist, err := redis.Bool(c.Do("EXISTS", util.key())); err == nil {
		return isExist
	}

	return false
}

// RefreshExpire 刷新过期时间
func (util *serverUtil) RefreshExpire(expire uint32) error {
	c := GetSingletonRedis()
	defer c.Close()

	_, err := c.Do("EXPIRE", util.key(), expire)
	return err
}

// GetToken 获取Token
func (util *serverUtil) GetToken() (string, error) {
	return redis.String(util.getValue("token"))
}

// Delete 删除key
func (util *serverUtil) Delete() error {
	c := GetSingletonRedis()
	defer c.Close()
	if _, err := c.Do("DEL", util.key()); err != nil {
		return err
	}
	if _, err := c.Do("SREM", servers, util.uid); err != nil {
		return err
	}
	return nil
}

func (util *serverUtil) setValue(field string, value interface{}) error {
	c := GetSingletonRedis()
	defer c.Close()
	_, err := c.Do("HSET", util.key(), field, value)
	return err
}

func (util *serverUtil) getValue(field string) (interface{}, error) {
	c := GetSingletonRedis()
	defer c.Close()
	return c.Do("HGET", util.key(), field)
}

func (util *serverUtil) key() string {
	return fmt.Sprintf("%s%d", serverPrefix, util.uid)
}
