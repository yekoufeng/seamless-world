package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

/*
 以hash表存储entity数据
 key: type:dbid
 如 player:1000
*/

type entityUtil struct {
	typ  string
	dbid uint64
}

// EntityUtil 获得Entity工具类
func EntityUtil(typ string, dbid uint64) *entityUtil {
	enUtil := &entityUtil{}
	enUtil.typ = typ
	enUtil.dbid = dbid
	return enUtil
}

func (util *entityUtil) GetValues(args []interface{}) ([]interface{}, error) {
	c := Get()
	defer c.Close()

	return redis.Values(c.Do("HMGET", append([]interface{}{util.key()}, args...)...))
}

func (util *entityUtil) SetValues(args []interface{}) error {
	c := Get()
	defer c.Close()

	_, err := c.Do("HMSET", append([]interface{}{util.key()}, args...)...)
	return err
}

func (util *entityUtil) GetValue(k string) (interface{}, error) {
	c := Get()
	defer c.Close()

	return c.Do("HGET", util.key(), k)
}

func (util *entityUtil) SetValue(k string, v interface{}) error {
	c := Get()
	defer c.Close()

	_, err := c.Do("HSET", util.key(), k, v)
	return err
}

func (util *entityUtil) key() string {
	return fmt.Sprintf("%s:%d", util.typ, util.dbid)
}
