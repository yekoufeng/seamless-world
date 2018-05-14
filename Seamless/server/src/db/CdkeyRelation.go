package db

import (
	//"common"
	"strconv"
	"zeus/dbservice"

	"github.com/garyburd/redigo/redis"

	log "github.com/cihub/seelog"
)

// 判断是否激活码是否被使用
func JudgeCDKEYInUse(cdkey string) bool {
	return hExists("UsedCdkey", cdkey)
}

// 设置激活码被使用
func SetCDKEYUse(cdkey string, account string) {

	if JudgeCDKEYInUse(cdkey) == false {
		hSet("UsedCdkey", cdkey, account)
	}
}

// 激活码状态
const (
	CDKEYEXIST  = 0 // 激活存在未使用
	CDKEYINUSE  = 1 // 激活码已被使用
	CDKEYABSEND = 2 // 激活码不存在
)

// 填充激活码
func FillCDKEY(cdkey string) bool {

	c := dbservice.Get()
	defer c.Close()

	result, err := redis.Int(c.Do("HSETNX", "CDKEYStateLib", cdkey, "false"))
	if err != nil {
		log.Error("填充激活码失败", err)
		return false
	}

	if result == 1 {
		return true
	}

	return false
}

// 设置激活码被使用
func SetCDKEYBeUsed(cdkey string, userName string) bool {

	if !hExists("CDKEYStateLib", cdkey) {
		return false
	}

	state := hGet("CDKEYStateLib", cdkey)
	if state == "true" { // 已被使用
		return false
	}

	hSet("CDKEYStateLib", cdkey, "true")
	hSet("BeUsedCdkey", userName, cdkey)

	return true

}

// 获取激活码状态
func GetCDKEYState(cdkey string) uint32 {

	if !hExists("CDKEYStateLib", cdkey) {
		return CDKEYABSEND
	}

	state := hGet("CDKEYStateLib", cdkey)

	if state == "false" {
		return CDKEYEXIST
	}

	return CDKEYINUSE
}

// 判断角色名是否被使用
func JudgeNameInUse(name string) bool {

	return hExists("UsedNameLib", name)

}

// 添加已用角色名
func AddUsedName(name string, id uint64) {

	if hExists("UsedNameLib", name) {
		return
	}

	hSet("UsedNameLib", name, id)
}

// 根据名称,获取玩家id
func GetIDByName(name string) uint64 {

	if !hExists("UsedNameLib", name) {
		return 0
	}

	idStr := hGet("UsedNameLib", name)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("GetIDByID failed ", err, name, idStr)
		return 0
	}

	return uint64(id)
}
