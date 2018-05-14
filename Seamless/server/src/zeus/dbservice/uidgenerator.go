package dbservice

import (
	"github.com/garyburd/redigo/redis"
)

type uidGenerator struct {
}

// UIDGenerator UID生成器
func UIDGenerator() *uidGenerator {
	return &uidGenerator{}
}

// Get 获取指定类型的UID, 保证全局唯一, 从DB库中生成
func (util *uidGenerator) Get(field string) (uint64, error) {
	c := Get()
	defer c.Close()

	return redis.Uint64(c.Do("HINCRBY", "uidgenerator", field, 1))
}
