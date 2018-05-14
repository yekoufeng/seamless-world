package dbservice

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

type sessionUtil struct {
	uid uint64
}

// SessionInfo 字段表
type SessionInfo struct {
	Token    string `redis:"token"`
	EntityID uint64 `redis:"entityID"`
	IP       string `redis:"ip"`
}

const (
	// SessionPrefix redis key后缀
	sessionPrefix = "session"
	// redisEntityID Session表中存放entityID的字段
	redisEntityID = "entityID"
	// loginIP 玩家登录IP
	loginIP = "ip"
)

const (
	_ = iota
	userStateLogouting
)

var TokenVerifyError = errors.New("Token Error")

// SessionUtil 获得Session表的工具类
func SessionUtil(uid uint64) *sessionUtil {
	sessUtil := &sessionUtil{}
	sessUtil.uid = uid
	return sessUtil
}

// SaveSessionInfo 保存session信息, 包括entityID等
func (util *sessionUtil) SaveSessionInfo(info *SessionInfo) error {
	c := GetServerRedis()
	defer c.Close()

	if _, err := c.Do("HMSET", redis.Args{}.Add(util.key()).AddFlat(info)...); err != nil {
		return err
	}
	return nil
}

// SetToken 生成token并将token保存至redis
// 返回token和err
func (util *sessionUtil) SetToken() (string, error) {
	t := util.createToken()
	if err := util.setValue("token", t); err != nil {
		return "", err
	}
	if err := util.setExpire(30); err != nil {
		return "", err
	}
	return t, nil
}

// GetToken 返回token
// 返回token和err
func (util *sessionUtil) GetToken() (string, error) {
	return redis.String(util.getValue("token"))
}

// VerifyToken 验证Token有效性
// 返回验证结果: true有效, false无效
func (util *sessionUtil) VerifyToken(token string) bool {
	t, err := util.GetToken()
	if err != nil || t != token {
		return false
	}
	if err := util.removeExpire(); err != nil {
		log.Error(err)
	}
	return true
}

// DelToken 删除Token
// 返回err
func (util *sessionUtil) DelToken() error {
	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("HDEL", util.key(), "token")
	if err != nil {
		return err
	}

	return nil
}

// DelSession 删除Session表
func (util *sessionUtil) DelSession(token string) error {
	if !util.VerifyToken(token) {
		return TokenVerifyError
	}

	c := GetServerRedis()
	defer c.Close()

	_, err := c.Do("DEL", util.key())
	return err
}

// IsExisted Session表是否存在
func (util *sessionUtil) IsExisted() (bool, error) {
	c := GetServerRedis()
	defer c.Close()

	reply, err := c.Do("EXISTS", util.key())
	if err != nil {
		return false, err
	}
	return redis.Bool(reply, nil)
}

// SetUserEntityID 设置玩家的EntityID
// func (util *sessionUtil) SetUserEntityID(entityID uint64) error {
// 	return util.setValue(redisEntityID, entityID)
// }

// GetUserEntityID 获得玩家的EntityID
func (util *sessionUtil) GetUserEntityID() (uint64, error) {
	reply, err := util.getValue(redisEntityID)
	if err != nil {
		return 0, err
	}
	if reply == nil {
		return 0, errors.New("not user login")
	}
	return redis.Uint64(reply, err)
}

// GetIP 获取玩家的登录IP
func (util *sessionUtil) GetIP() (string, error) {
	reply, err := util.getValue(loginIP)
	if err != nil {
		return "", err
	}
	return redis.String(reply, nil)
}

func (util *sessionUtil) getValue(field string) (interface{}, error) {
	c := GetServerRedis()
	defer c.Close()
	return c.Do("HGET", util.key(), field)
}

func (util *sessionUtil) setValue(field string, value interface{}) error {
	c := GetServerRedis()
	defer c.Close()
	_, err := c.Do("HSET", util.key(), field, value)
	return err
}

func (util *sessionUtil) setExpire(seconds int) error {
	c := GetServerRedis()
	defer c.Close()
	_, err := c.Do("EXPIRE", util.key(), seconds)
	return err
}

func (util *sessionUtil) removeExpire() error {
	c := GetServerRedis()
	defer c.Close()
	_, err := c.Do("PERSIST", util.key())
	return err
}

// getSessionKey 获取session表的Key
func (util *sessionUtil) key() string {
	return fmt.Sprintf("%s:%d", sessionPrefix, util.uid)
}

// createToken 方法根据uid和当前时间生成唯一的token
// md5(curtime+uid)
func (util *sessionUtil) createToken() string {
	curtime := time.Now().Unix()
	h := md5.New()

	io.WriteString(h, strconv.FormatInt(curtime, 10))
	io.WriteString(h, strconv.FormatUint(util.uid, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}
