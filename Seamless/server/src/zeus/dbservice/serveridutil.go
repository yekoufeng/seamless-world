package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	serverIDPrefix = "serverid:"
	startIDKey     = "startid"
)

type srvIDUtil struct {
	srvID uint64
}

// SrvIDUtil 获得SrvID工具类
func SrvIDUtil(srvID uint64) *srvIDUtil {
	return &srvIDUtil{
		srvID: srvID,
	}
}

// GetStartID 获取起始ID并累加
func (util *srvIDUtil) GetStartID() (uint64, error) {
	c := GetSingletonRedis()
	defer c.Close()

	r, err := redis.Bool(c.Do("HEXISTS", util.key(), startIDKey))
	if err != nil {
		return 0, err
	}

	if !r {
		_, err := c.Do("HSET", util.key(), startIDKey, 0)
		if err != nil {
			return 0, err
		}
	}

	n, err := redis.Uint64(c.Do("HGET", util.key(), startIDKey))

	if err != nil {
		return 0, err
	}

	n++
	if n >= (1 << 5) {
		n = 1
	}

	_, err = c.Do("HSET", util.key(), startIDKey, n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (util *srvIDUtil) key() string {
	return fmt.Sprintf("%s%d", serverIDPrefix, util.srvID)
}
