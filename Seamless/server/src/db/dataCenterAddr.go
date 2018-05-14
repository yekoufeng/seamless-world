package db

import (
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"
)

const (
	dataCenterAddr = "DataCenterAddr"
)

// GetDataCenterAddr 获取DataCenter地址
func GetDataCenterAddr(field string) (string, error) {
	reply, err := getValue(field)
	if err != nil || reply == nil {
		return "", err
	}

	v, err := redis.String(reply, nil)
	if err != nil {
		return "", err
	}
	return v, nil
}

func getValue(field string) (interface{}, error) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	return c.Do("HGET", dataCenterAddr, field)
}

// SetDataCenterAddr 设置DataCenter地址
func SetDataCenterAddr(field1 string, value1 interface{}, field2 string, value2 interface{}) error {
	c := dbservice.GetServerRedis()
	defer c.Close()

	_, err := c.Do("HMSET", dataCenterAddr, field1, value1, field2, value2)
	return err
}

// DelDataCenterAddr 删除DataCenter地址
func DelDataCenterAddr(field1, field2 string) error {
	c := dbservice.GetServerRedis()
	defer c.Close()

	_, err := c.Do("HDEL", dataCenterAddr, field1, field2)
	return err
}
