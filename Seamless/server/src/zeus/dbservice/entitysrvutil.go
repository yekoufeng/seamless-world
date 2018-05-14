package dbservice

import (
	"errors"
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// EntitySrvInfo 服务器信息
type EntitySrvInfo struct {
	SrvID  uint64
	CellID uint64
}

/*
 entitysrvinfo:* 表工具类
 存储entity分布在哪些服务器上的信息
 id为entityid
*/

type entitySrvUtil struct {
	id uint64
}

const (
	entitySrvInfoPrefix = "entitysrvinfo"
	fieldType           = "entitytype"
	fieldDBID           = "dbid"
	fieldExistedCnt     = "existedcnt"
)

// EntitySrvUtil 获得Entity工具类
func EntitySrvUtil(eid uint64) *entitySrvUtil {
	enUtil := &entitySrvUtil{}
	enUtil.id = eid
	return enUtil
}

// IsExist 这个ID号的Entity是否存在
func (util *entitySrvUtil) IsExist() bool {
	c := GetServerRedis()
	defer c.Close()

	r, err := c.Do("EXISTS", util.key())
	if err != nil {
		log.Error(err)
		return false
	}

	v, err := redis.Bool(r, nil)
	if err != nil {
		log.Error(err)
		return false
	}
	return v
}

// RegSrvID 注册类型和服务器信息
func (util *entitySrvUtil) RegSrvID(srvType uint8, srvID uint64, cellID uint64, entityType string, dbID uint64) error {
	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("HMSET", util.key(), srvType, util.joinSrvInfo(srvID, cellID), fieldType, entityType, fieldDBID, dbID)
	if err != nil {
		log.Error(err)
	}
	_, err = c.Do("HINCRBY", util.key(), fieldExistedCnt, 1)
	return err
}

// UnRegSrvID 删除注册信息
func (util *entitySrvUtil) UnRegSrvID(srvType uint8, srvID uint64, cellID uint64) error {
	c := GetServerRedis()
	defer c.Close()

	reply, err := c.Do("HINCRBY", util.key(), fieldExistedCnt, -1)
	if err != nil {
		log.Error(err)
	}

	cnt, err := redis.Int(reply, nil)
	if err != nil {
		log.Error(err)
	}
	if cnt <= 0 {
		_, err = c.Do("DEL", util.key())
		return err
	}

	_srvID, _cellID, err := util.GetSrvInfo(srvType)
	if err != nil {
		return err
	}

	if _srvID == srvID && _cellID == cellID {
		_, err = c.Do("HDEL", util.key(), srvType)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

// GetEntityInfo 获取entityType 和 dbID
func (util *entitySrvUtil) GetEntityInfo() (string, uint64, error) {
	c := GetServerRedis()
	defer c.Close()

	ret, err := c.Do("HMGET", util.key(), fieldType, fieldDBID)
	if err != nil {
		return "", 0, err
	}

	retValue, err := redis.Values(ret, nil)
	if err != nil {
		return "", 0, err
	}

	if len(retValue) != 2 {
		return "", 0, errors.New("wrong")
	}

	entityType, err := redis.String(retValue[0], nil)
	if err != nil {
		return "", 0, err
	}

	dbID, err := redis.Uint64(retValue[1], nil)
	if err != nil {
		return "", 0, err
	}

	return entityType, dbID, nil
}

// GetSrvInfo 获取特定服务器类型的 服务器 ID以及 CellID
func (util *entitySrvUtil) GetSrvInfo(srvType uint8) (srvID uint64, cellID uint64, err error) {
	c := GetServerRedis()
	defer c.Close()

	ret, err := c.Do("HGET", util.key(), srvType)
	if err != nil {
		return
	}

	retStr, err := redis.String(ret, nil)
	if err != nil {
		return
	}

	srvID, cellID = util.splitSrvInfo(retStr)

	return
}

// GetSrvIDs 获取Entity的分布式信息
func (util *entitySrvUtil) GetSrvIDs() (map[uint8]*EntitySrvInfo, error) {
	c := GetServerRedis()
	defer c.Close()

	reply, err := c.Do("HGETALL", util.key())
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}

	values, err := redis.Values(reply, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[uint8]*EntitySrvInfo)

	for i := 0; i < len(values); i += 2 {
		srvType, err := redis.Uint64(values[i], nil)
		if err != nil {
			continue
		}

		s, _ := redis.String(values[i+1], nil)

		srvID, CellID := util.splitSrvInfo(s)

		srvInfo := EntitySrvInfo{srvID, CellID}

		result[uint8(srvType)] = &srvInfo
	}
	return result, nil
}

// func (util *entitySrvUtil) SetTransporting(b bool) {

// 	c := GetServerRedis()
// 	defer c.Close()

// 	if b {
// 		_, err := c.Do("HSET", util.key(), "IsTransport", 0)
// 		if err != nil {
// 			log.Error("set transport error")
// 		}
// 	} else {
// 		_, err := c.Do("HDEL", util.key(), "IsTransport")
// 		if err != nil {
// 			log.Error("del transport error")
// 		}
// 	}

// }

// func (util *entitySrvUtil) IsTransporting() bool {

// 	c := GetServerRedis()
// 	defer c.Close()

// 	r, err := c.Do("HEXISTS", util.key(), "IsTransport")
// 	if err != nil {
// 		log.Error("get is transport error")
// 		return true
// 	}

// 	b, _ := redis.Bool(r, nil)
// 	return b
// }

// GetCellInfo 获取包含Cell的srvID 和 cellID
func (util *entitySrvUtil) GetCellInfo() (uint64, uint64, error) {

	srvInfos, err := util.GetSrvIDs()
	if err != nil {
		return 0, 0, err
	}

	for _, info := range srvInfos {
		if info.CellID != 0 {
			return info.SrvID, info.CellID, nil
		}
	}

	return 0, 0, nil
}

func (util *entitySrvUtil) joinSrvInfo(t uint64, id uint64) string {
	return fmt.Sprintf("%d:%d", t, id)
}

func (util *entitySrvUtil) splitSrvInfo(s string) (srvID uint64, cellID uint64) {
	fmt.Sscanf(s, "%d:%d", &srvID, &cellID)
	return
}

func (util *entitySrvUtil) key() string {
	return fmt.Sprintf("%s:%d", entitySrvInfoPrefix, util.id)
}
