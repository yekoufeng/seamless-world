package main

// LoginConfig 登录配置
type LoginConfig struct {
	HBEnable bool
}

// AccountActiveReq 账号激活请求格式
type AccountActiveReq struct {
	Activatetype uint64 // 激活类型 0、手机号激活 1、激活码激活
	User         string
	Password     string
	CDKEY        string // 激活码
}

// AccountActiveRet 账号激活返回格式
type AccountActiveRet struct {
	RetCode int // 返回码
}

/*
// 激活码状态
const (
	CDKEYEXIST  = 0 // 激活存在未使用
	CDKEYINUSE  = 1 // 激活码已被使用
	CDKEYABSEND = 2 // 激活码不存在
)
*/

// ActiveRetType	账号激活返回类型
const (
	ACTIVESUCCUESS     = 0 // 激活成功
	ACTIVECDKEYINUSE   = 1 // 激活码已被使用
	ACTIVEABSEND       = 2 // 激活码不存在
	ACTIVEACCOUNTEXIST = 3 // 账号已存在
)

// QQLogin QQ登录json结构
type QQLogin struct {
	AppID   int    `json:"appid"`
	OpenID  string `json:"openid"`
	OpenKey string `json:"openkey"`
	UserIP  string `json:"userip"`
	Enter   string `json:"gameEnter"`
}

// QQLoginRet QQ登录返回结构
type QQLoginRet struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

// WXLogin 微信登录json结构
type WXLogin struct {
	OpenID      string `json:"openid"`
	AccessToken string `json:"accessToken"`
	Enter       string `json:"gameEnter"`
}

// WXLoginRet QQ登录返回结构
type WXLoginRet struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

// GuestLogin 游客登录json结构
type GuestLogin struct {
	GuestID     string `json:"guestid"`
	AccessToken string `json:"accessToken"`
	Enter       string `json:"gameEnter"`
}

// GuestLoginRet QQ登录返回结构
type GuestLoginRet struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}
