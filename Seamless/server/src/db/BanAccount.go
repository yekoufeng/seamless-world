package db

import (
	"time"

	"encoding/json"

	log "github.com/cihub/seelog"
)

const (
	BanAccountRedis = "BanAccount"
	BanRoleRedis    = "BanRole"
)

// BanAccount 账号封停信息
type BanAccount struct {
	Uid         uint64
	OpenID      string
	BanDuration int32  // 封停时长，-1表示永久封号
	EndTime     uint32 // 封停结束时间
	BanReason   string // 封号原因：（自定义文字，玩家登录时客户端可见）
}

// AddBanAccount 添加账号封停
func AddBanAccount(data *BanAccount) bool {

	if data == nil {
		log.Warn("AddBanAccount data is nil")
		return false
	}

	if hExists(BanAccountRedis, data.OpenID) {
		log.Info("账号已经进行了封停继续封停 ", data.OpenID)
	}

	d, err := json.Marshal(*data)
	if err != nil {
		log.Info("AddBanAccount err ", err)
		return false
	}

	hSet(BanAccountRedis, data.OpenID, string(d)) // 添加封停账号

	log.Info("AddBanAccount ", data)
	return true
}

// GetBanAccountInfo 获取账号封停信息
func GetBanAccountData(account string) *BanAccount {

	if !hExists(BanAccountRedis, account) {
		return nil
	}

	v := hGet(BanAccountRedis, account)
	var d *BanAccount
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GetBanAccountInfo Failed to Unmarshal ", err)
		return nil
	}

	curtime := time.Now().Unix()
	if int64(d.EndTime) < curtime && d.BanDuration != -1 {
		UnbanAccount(account)
		return nil
	}

	return d
}

// UnbanAccount 解封封停账号
func UnbanAccount(account string) bool {
	if !hExists(BanAccountRedis, account) {
		return false
	}

	hDEL(BanAccountRedis, account)
	return true
}

// BanRole 角色封停信息
type BanRole struct {
	Uid         uint64
	OpenID      string
	RoleID      string
	BanDuration int32  // 封停时长，-1表示永久封号
	EndTime     uint32 // 封停结束时间
	BanReason   string // 封号原因：（自定义文字，玩家登录时客户端可见）
}

// AddBanRole 添加角色封停
func AddBanRole(data *BanRole) bool {

	if data == nil {
		log.Warn("AddBanRole data is nil")
		return false
	}

	if hExists(BanRoleRedis, data.OpenID) {
		log.Info("角色已经进行了封停继续封停 ", data.OpenID)
	}

	d, err := json.Marshal(*data)
	if err != nil {
		log.Info("AddBanRole err ", err)
		return false
	}

	hSet(BanRoleRedis, data.RoleID, string(d)) // 添加封停角色

	log.Info("AddBanRole ", data)
	return true

}

// GetBanRoleInfo 获取角色封停信息
func GetBanRoleData(role string) *BanRole {

	if !hExists(BanRoleRedis, role) {
		return nil
	}

	v := hGet(BanRoleRedis, role)
	var d *BanRole
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GetBanRoleData Failed to Unmarshal ", err)
		return nil
	}

	curtime := time.Now().Unix()
	if int64(d.EndTime) < curtime && d.BanDuration != -1 {
		UnbanRole(role)
		return nil
	}

	return d

}

// UnbanRole 解封封停角色
func UnbanRole(role string) bool {
	if !hExists(BanRoleRedis, role) {
		return false
	}

	hDEL(BanRoleRedis, role)
	return true
}
