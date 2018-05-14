package dbservice

import (
	"github.com/garyburd/redigo/redis"
)

// CacheHSET 缓存hset
func CacheHSET(key string, field interface{}, value interface{}) error {
	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("HSET", key, field, value)
	return err
}

// CacheHGET 缓存hget
func CacheHGET(key string, field interface{}) (interface{}, error) {
	c := GetServerRedis()
	defer c.Close()

	return c.Do("HGET", key, field)
}

// CacheHDEL 缓存hdel
func CacheHDEL(key string, field interface{}) error {
	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("HDEL", key, field)
	return err
}

// CacheHExists 缓存HExists
func CacheHExists(key string, field interface{}) (bool, error) {
	c := GetServerRedis()
	defer c.Close()

	return redis.Bool(c.Do("HEXISTS", key, field))
}
