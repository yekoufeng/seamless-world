package dbservice

import (
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

type entityTypeUtil struct {
	entityType string
}

const (
	entityTypePrefix = "entitytypeinfo"
)

// EntityTypeUtil 获得Entity工具类
func EntityTypeUtil(typ string) *entityTypeUtil {
	enUtil := &entityTypeUtil{}
	enUtil.entityType = typ
	return enUtil
}

// RegSrvType 注册服务器类型
func (util *entityTypeUtil) RegSrvType(srvType uint8) {
	c := GetServerRedis()
	defer c.Close()
	_, err := c.Do("SADD", util.key(), srvType)
	if err != nil {
		log.Error(err)
	}
}

// GetSrvType 获取服务器类型列表
func (util *entityTypeUtil) GetSrvType() ([]uint8, error) {
	c := GetServerRedis()
	defer c.Close()

	r, err := c.Do("SMEMBERS", util.key())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	values, err := redis.Values(r, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tl := make([]uint8, 0, 10)
	for i := 0; i < len(values); i++ {

		v, e := redis.Uint64(values[i], nil)
		if e != nil {
			log.Error(e)
			return nil, e
		}

		tl = append(tl, uint8(v))
	}

	return tl, nil
}

func (util *entityTypeUtil) key() string {
	return entityTypePrefix + ":" + util.entityType
}
