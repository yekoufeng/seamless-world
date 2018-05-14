package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type cellUtil struct {
	id uint64
}

const (
	cellPrefix = "space"
)

// CellUtil 获取Util
func CellUtil(id uint64) *cellUtil {
	return &cellUtil{id: id}
}

// RegSrvID 设置服务器ID
func (util *cellUtil) RegSrvID(srvID uint64) error {
	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("HSET", util.key(), "SrvID", srvID)
	if err != nil {
		return err
	}

	return nil
}

// UnReg 设置服务器
func (util *cellUtil) UnReg() error {
	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("DEL", util.key())
	if err != nil {
		return err
	}
	return nil
}

func (util *cellUtil) IsExist() (bool, error) {

	c := GetServerRedis()
	defer c.Close()

	r, err := c.Do("EXISTS", util.key())
	if err != nil {
		return false, err
	}

	ret, err := redis.Bool(r, nil)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// GetSrvID 获取服务器ID
func (util *cellUtil) GetSrvID() (uint64, error) {

	c := GetServerRedis()
	defer c.Close()

	r, err := c.Do("HGET", util.key(), "SrvID")
	if err != nil {
		return 0, err
	}

	srvID, err := redis.Uint64(r, nil)
	if err != nil {
		return 0, err
	}

	return srvID, nil
}

func (util *cellUtil) key() string {
	return fmt.Sprintf("%s:%d", cellPrefix, util.id)
}
