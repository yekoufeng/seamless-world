package entitydef

import "zeus/iserver"
import "protoMsg"

// PlayerDef 自动生成的属性包装代码
type PlayerDef struct {
	ip iserver.IEntityProps
}

// SetPropsSetter 设置接口
func (p *PlayerDef) SetPropsSetter(ip iserver.IEntityProps) {
	p.ip = ip
}

// SetAccessToken 设置 AccessToken
func (p *PlayerDef) SetAccessToken(v string) {
	p.ip.SetProp("AccessToken", v)
}

// SetAccessTokenDirty 设置AccessToken被修改
func (p *PlayerDef) SetAccessTokenDirty() {
	p.ip.PropDirty("AccessToken")
}

// GetAccessToken 获取 AccessToken
func (p *PlayerDef) GetAccessToken() string {
	v := p.ip.GetProp("AccessToken")
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetAttack 设置 Attack
func (p *PlayerDef) SetAttack(v uint32) {
	p.ip.SetProp("Attack", v)
}

// SetAttackDirty 设置Attack被修改
func (p *PlayerDef) SetAttackDirty() {
	p.ip.PropDirty("Attack")
}

// GetAttack 获取 Attack
func (p *PlayerDef) GetAttack() uint32 {
	v := p.ip.GetProp("Attack")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetCoin 设置 Coin
func (p *PlayerDef) SetCoin(v uint64) {
	p.ip.SetProp("Coin", v)
}

// SetCoinDirty 设置Coin被修改
func (p *PlayerDef) SetCoinDirty() {
	p.ip.PropDirty("Coin")
}

// GetCoin 获取 Coin
func (p *PlayerDef) GetCoin() uint64 {
	v := p.ip.GetProp("Coin")
	if v == nil {
		return uint64(0)
	}

	return v.(uint64)
}

// SetDefence 设置 Defence
func (p *PlayerDef) SetDefence(v uint32) {
	p.ip.SetProp("Defence", v)
}

// SetDefenceDirty 设置Defence被修改
func (p *PlayerDef) SetDefenceDirty() {
	p.ip.PropDirty("Defence")
}

// GetDefence 获取 Defence
func (p *PlayerDef) GetDefence() uint32 {
	v := p.ip.GetProp("Defence")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetFriendsNum 设置 FriendsNum
func (p *PlayerDef) SetFriendsNum(v uint32) {
	p.ip.SetProp("FriendsNum", v)
}

// SetFriendsNumDirty 设置FriendsNum被修改
func (p *PlayerDef) SetFriendsNumDirty() {
	p.ip.PropDirty("FriendsNum")
}

// GetFriendsNum 获取 FriendsNum
func (p *PlayerDef) GetFriendsNum() uint32 {
	v := p.ip.GetProp("FriendsNum")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetGameEnter 设置 GameEnter
func (p *PlayerDef) SetGameEnter(v string) {
	p.ip.SetProp("GameEnter", v)
}

// SetGameEnterDirty 设置GameEnter被修改
func (p *PlayerDef) SetGameEnterDirty() {
	p.ip.PropDirty("GameEnter")
}

// GetGameEnter 获取 GameEnter
func (p *PlayerDef) GetGameEnter() string {
	v := p.ip.GetProp("GameEnter")
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetGender 设置 Gender
func (p *PlayerDef) SetGender(v string) {
	p.ip.SetProp("Gender", v)
}

// SetGenderDirty 设置Gender被修改
func (p *PlayerDef) SetGenderDirty() {
	p.ip.PropDirty("Gender")
}

// GetGender 获取 Gender
func (p *PlayerDef) GetGender() string {
	v := p.ip.GetProp("Gender")
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetHP 设置 HP
func (p *PlayerDef) SetHP(v uint32) {
	p.ip.SetProp("HP", v)
}

// SetHPDirty 设置HP被修改
func (p *PlayerDef) SetHPDirty() {
	p.ip.PropDirty("HP")
}

// GetHP 获取 HP
func (p *PlayerDef) GetHP() uint32 {
	v := p.ip.GetProp("HP")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetLevel 设置 Level
func (p *PlayerDef) SetLevel(v uint32) {
	p.ip.SetProp("Level", v)
}

// SetLevelDirty 设置Level被修改
func (p *PlayerDef) SetLevelDirty() {
	p.ip.PropDirty("Level")
}

// GetLevel 获取 Level
func (p *PlayerDef) GetLevel() uint32 {
	v := p.ip.GetProp("Level")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetLoginTime 设置 LoginTime
func (p *PlayerDef) SetLoginTime(v int64) {
	p.ip.SetProp("LoginTime", v)
}

// SetLoginTimeDirty 设置LoginTime被修改
func (p *PlayerDef) SetLoginTimeDirty() {
	p.ip.PropDirty("LoginTime")
}

// GetLoginTime 获取 LoginTime
func (p *PlayerDef) GetLoginTime() int64 {
	v := p.ip.GetProp("LoginTime")
	if v == nil {
		return int64(0)
	}

	return v.(int64)
}

// SetLogoutTime 设置 LogoutTime
func (p *PlayerDef) SetLogoutTime(v int64) {
	p.ip.SetProp("LogoutTime", v)
}

// SetLogoutTimeDirty 设置LogoutTime被修改
func (p *PlayerDef) SetLogoutTimeDirty() {
	p.ip.PropDirty("LogoutTime")
}

// GetLogoutTime 获取 LogoutTime
func (p *PlayerDef) GetLogoutTime() int64 {
	v := p.ip.GetProp("LogoutTime")
	if v == nil {
		return int64(0)
	}

	return v.(int64)
}

// SetMaxHP 设置 MaxHP
func (p *PlayerDef) SetMaxHP(v uint32) {
	p.ip.SetProp("MaxHP", v)
}

// SetMaxHPDirty 设置MaxHP被修改
func (p *PlayerDef) SetMaxHPDirty() {
	p.ip.PropDirty("MaxHP")
}

// GetMaxHP 获取 MaxHP
func (p *PlayerDef) GetMaxHP() uint32 {
	v := p.ip.GetProp("MaxHP")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetName 设置 Name
func (p *PlayerDef) SetName(v string) {
	p.ip.SetProp("Name", v)
}

// SetNameDirty 设置Name被修改
func (p *PlayerDef) SetNameDirty() {
	p.ip.PropDirty("Name")
}

// GetName 获取 Name
func (p *PlayerDef) GetName() string {
	v := p.ip.GetProp("Name")
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetNickName 设置 NickName
func (p *PlayerDef) SetNickName(v string) {
	p.ip.SetProp("NickName", v)
}

// SetNickNameDirty 设置NickName被修改
func (p *PlayerDef) SetNickNameDirty() {
	p.ip.PropDirty("NickName")
}

// GetNickName 获取 NickName
func (p *PlayerDef) GetNickName() string {
	v := p.ip.GetProp("NickName")
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetOnlineTime 设置 OnlineTime
func (p *PlayerDef) SetOnlineTime(v int64) {
	p.ip.SetProp("OnlineTime", v)
}

// SetOnlineTimeDirty 设置OnlineTime被修改
func (p *PlayerDef) SetOnlineTimeDirty() {
	p.ip.PropDirty("OnlineTime")
}

// GetOnlineTime 获取 OnlineTime
func (p *PlayerDef) GetOnlineTime() int64 {
	v := p.ip.GetProp("OnlineTime")
	if v == nil {
		return int64(0)
	}

	return v.(int64)
}

// SetPicture 设置 Picture
func (p *PlayerDef) SetPicture(v string) {
	p.ip.SetProp("Picture", v)
}

// SetPictureDirty 设置Picture被修改
func (p *PlayerDef) SetPictureDirty() {
	p.ip.PropDirty("Picture")
}

// GetPicture 获取 Picture
func (p *PlayerDef) GetPicture() string {
	v := p.ip.GetProp("Picture")
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetPlayerLogin 设置 PlayerLogin
func (p *PlayerDef) SetPlayerLogin(v *protoMsg.PlayerLogin) {
	p.ip.SetProp("PlayerLogin", v)
}

// SetPlayerLoginDirty 设置PlayerLogin被修改
func (p *PlayerDef) SetPlayerLoginDirty() {
	p.ip.PropDirty("PlayerLogin")
}

// GetPlayerLogin 获取 PlayerLogin
func (p *PlayerDef) GetPlayerLogin() *protoMsg.PlayerLogin {
	v := p.ip.GetProp("PlayerLogin")
	if v == nil {
		return nil
	}

	return v.(*protoMsg.PlayerLogin)
}

// SetQQVIP 设置 QQVIP
func (p *PlayerDef) SetQQVIP(v uint8) {
	p.ip.SetProp("QQVIP", v)
}

// SetQQVIPDirty 设置QQVIP被修改
func (p *PlayerDef) SetQQVIPDirty() {
	p.ip.PropDirty("QQVIP")
}

// GetQQVIP 获取 QQVIP
func (p *PlayerDef) GetQQVIP() uint8 {
	v := p.ip.GetProp("QQVIP")
	if v == nil {
		return uint8(0)
	}

	return v.(uint8)
}

// SetRoleModel 设置 RoleModel
func (p *PlayerDef) SetRoleModel(v uint32) {
	p.ip.SetProp("RoleModel", v)
}

// SetRoleModelDirty 设置RoleModel被修改
func (p *PlayerDef) SetRoleModelDirty() {
	p.ip.PropDirty("RoleModel")
}

// GetRoleModel 获取 RoleModel
func (p *PlayerDef) GetRoleModel() uint32 {
	v := p.ip.GetProp("RoleModel")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetState 设置 State
func (p *PlayerDef) SetState(v uint32) {
	p.ip.SetProp("State", v)
}

// SetStateDirty 设置State被修改
func (p *PlayerDef) SetStateDirty() {
	p.ip.PropDirty("State")
}

// GetState 获取 State
func (p *PlayerDef) GetState() uint32 {
	v := p.ip.GetProp("State")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// SetTodayOnlineTime 设置 TodayOnlineTime
func (p *PlayerDef) SetTodayOnlineTime(v int64) {
	p.ip.SetProp("TodayOnlineTime", v)
}

// SetTodayOnlineTimeDirty 设置TodayOnlineTime被修改
func (p *PlayerDef) SetTodayOnlineTimeDirty() {
	p.ip.PropDirty("TodayOnlineTime")
}

// GetTodayOnlineTime 获取 TodayOnlineTime
func (p *PlayerDef) GetTodayOnlineTime() int64 {
	v := p.ip.GetProp("TodayOnlineTime")
	if v == nil {
		return int64(0)
	}

	return v.(int64)
}

type IPlayerDef interface {
	SetAccessToken(v string)
	SetAccessTokenDirty()
	GetAccessToken() string
	SetAttack(v uint32)
	SetAttackDirty()
	GetAttack() uint32
	SetCoin(v uint64)
	SetCoinDirty()
	GetCoin() uint64
	SetDefence(v uint32)
	SetDefenceDirty()
	GetDefence() uint32
	SetFriendsNum(v uint32)
	SetFriendsNumDirty()
	GetFriendsNum() uint32
	SetGameEnter(v string)
	SetGameEnterDirty()
	GetGameEnter() string
	SetGender(v string)
	SetGenderDirty()
	GetGender() string
	SetHP(v uint32)
	SetHPDirty()
	GetHP() uint32
	SetLevel(v uint32)
	SetLevelDirty()
	GetLevel() uint32
	SetLoginTime(v int64)
	SetLoginTimeDirty()
	GetLoginTime() int64
	SetLogoutTime(v int64)
	SetLogoutTimeDirty()
	GetLogoutTime() int64
	SetMaxHP(v uint32)
	SetMaxHPDirty()
	GetMaxHP() uint32
	SetName(v string)
	SetNameDirty()
	GetName() string
	SetNickName(v string)
	SetNickNameDirty()
	GetNickName() string
	SetOnlineTime(v int64)
	SetOnlineTimeDirty()
	GetOnlineTime() int64
	SetPicture(v string)
	SetPictureDirty()
	GetPicture() string
	SetPlayerLogin(v *protoMsg.PlayerLogin)
	SetPlayerLoginDirty()
	GetPlayerLogin() *protoMsg.PlayerLogin
	SetQQVIP(v uint8)
	SetQQVIPDirty()
	GetQQVIP() uint8
	SetRoleModel(v uint32)
	SetRoleModelDirty()
	GetRoleModel() uint32
	SetState(v uint32)
	SetStateDirty()
	GetState() uint32
	SetTodayOnlineTime(v int64)
	SetTodayOnlineTimeDirty()
	GetTodayOnlineTime() int64
}
