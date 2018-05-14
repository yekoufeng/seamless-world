package dbservice

import (
	"fmt"
	"math"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

type accountUtil struct {
	uid uint64
}

const (
	// AccountPrefix 帐号表前缀
	AccountPrefix = "account"

	// AccountOpenID 用户名表前缀, 存储帐号和UID的对应关系
	AccountOpenID = "accountopenid"

	// UIDField UID字段
	UIDField = "uid"
)

// GetUID 通过username获取uid
func GetUID(user string) (uint64, error) {
	c := Get()
	defer c.Close()

	return redis.Uint64(c.Do("HGET", AccountOpenID+":"+user, UIDField))
}

// Account 获取帐号表工具类
func Account(uid uint64) *accountUtil {
	acc := &accountUtil{}
	acc.uid = uid
	return acc
}

// SetPassword 密码保存至redis
func (util *accountUtil) SetPassword(password string) error {
	return util.setValue("password", password)
}

// VerifyPassword 验证密码
func (util *accountUtil) VerifyPassword(password string) bool {
	pwd, err := util.getPassword()
	if err != nil {
		log.Error("Get password failed", err)
		return false
	}

	if pwd != password {
		log.Error("Password not match")
		return false
	}

	return true
}

// SetUsername 保存用户名
func (util *accountUtil) SetUsername(user string) error {
	c := Get()
	defer c.Close()

	if reply, err := redis.Int(c.Do("HSETNX", util.key(), "username", user)); err != nil {
		return err
	} else if reply == 0 {
		return fmt.Errorf("Account existed %s", user)
	}

	_, err := c.Do("HSET", AccountOpenID+":"+user, UIDField, util.uid)
	return err
}

// GetUsername 获得用户名
func (util *accountUtil) GetUsername() (string, error) {
	return redis.String(util.getValue("username"))
}

// SetGrade 设置帐号级别 1表示内部 2表示外部玩家
func (util *accountUtil) SetGrade(grade uint32) error {
	return util.setValue("grade", grade)
}

// GetGrade 获取帐号级别
func (util *accountUtil) GetGrade() (uint32, error) {
	v, err := util.getValue("grade")
	if err != nil {
		return math.MaxUint32, err
	}
	if v == nil {
		return 0, nil
	}
	grade, err := redis.Uint64(v, nil)
	if err != nil {
		return math.MaxUint32, err
	}
	return uint32(grade), nil
}

// getPassword 方法根据用户名返回密码
func (util *accountUtil) getPassword() (string, error) {
	return redis.String(util.getValue("password"))
}

func (util *accountUtil) getValue(field string) (interface{}, error) {
	c := Get()
	defer c.Close()
	return c.Do("HGET", util.key(), field)
}

func (util *accountUtil) setValue(field string, value interface{}) error {
	c := Get()
	defer c.Close()
	_, err := c.Do("HSET", util.key(), field, value)
	return err
}

func (util *accountUtil) key() string {
	return fmt.Sprintf("%s:%d", AccountPrefix, util.uid)
}
