package main

import (
	"common"
	"db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"zeus/dbservice"
	"zeus/login"
	"zeus/serverMgr"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// UserLogin 玩家登录处理
func (s *Server) loginHandler(w rest.ResponseWriter, r *rest.Request) {
	if !(dbservice.DBValid && dbservice.SrvRedisValid && dbservice.SingletonRedisValid) {
		log.Error("服务器不可用")
		rest.Error(w, "网络异常, 请稍后再试", http.StatusInternalServerError)
		return
	}

	msg := login.UserLoginReq{}
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	if s.versionCheck && msg.ClientVer != s.allowVersion {
		log.Warn("客户端版本错误 ", msg)
		rest.Error(w, "客户端版本错误", http.StatusForbidden)
		return
	}

	ok, accessToken := s.channelVerify(&msg)
	if !ok {
		log.Warn("渠道验证失败 ", msg)
		rest.Error(w, "渠道验证失败", http.StatusForbidden)
		return
	}

	// 检查帐号是否存在
	var uid uint64
	if uid, err = dbservice.GetUID(msg.User); uid == 0 {
		// 用户不存在, 创建用户, 设置密码为用户输入的密码
		if s.forceCreate {
			uid, err = s.DoCreateNewUser(msg.User, msg.Password, uint32(s.initGrade))
			if err != nil {
				log.Error(err, msg)
				rest.Error(w, "服务器错误", http.StatusInternalServerError)
				return
			}
		} else {
			log.Warn("帐号不存在 ", msg)
			rest.Error(w, "帐号不存在", http.StatusForbidden)
			return
		}
	} else if err != nil {
		log.Error(err, msg)
		rest.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	} else {
		//用户存在, 验证密码
		if !dbservice.Account(uid).VerifyPassword(msg.Password) {
			log.Warn("帐号或密码错误 ", msg)
			rest.Error(w, "帐号或密码错误", http.StatusForbidden)
			return
		}

		// 账号是否被禁用
		isAccountBan, accountBanReason := s.VerifyAccountBan(msg.User)
		if isAccountBan {
			log.Warn("账号被禁用! ", msg)
			rest.Error(w, accountBanReason, http.StatusForbidden)
			return
		}

		// 账号下角色是否被禁用
		isRoleBan, roleBanReason := s.VerifyRoleBan(msg.User)
		if isRoleBan {
			log.Warn("账号下角色被禁用! ", msg)
			rest.Error(w, roleBanReason, http.StatusForbidden)
			return
		}

	}

	//检查白名单
	checkResult := s.checkWhiteList(uid)
	if checkResult == 0 {
		log.Warn("白名单检查失败", msg)
		rest.Error(w, "白名单检查失败", http.StatusForbidden)
		return
	} else if checkResult == 2 {
		log.Warn("服务器未开放")
		rest.Error(w, "服务器未开放", http.StatusForbidden)
		return
	}

	// 设置token, 返回客户端token和大厅服务器信息
	t, err := dbservice.SessionUtil(uid).SetToken()
	if err != nil {
		log.Error(err, msg)
		rest.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	//获取大厅服务器信息
	addr, err := s.GetGatewayAddr()
	if err == serverMgr.ErrServerBusy {
		log.Warn(err, msg)
		rest.Error(w, "服务器忙", http.StatusNotAcceptable)
		return
	}
	if err != nil {
		log.Error(err, msg)
		rest.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	err = dbservice.EntityUtil("Player", uid).SetValue("AccessToken", accessToken)
	if err != nil {
		log.Error(err, msg)
		rest.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	var loginChannel int
	if msg.Channel == "QQ" {
		loginChannel = 2
	} else if msg.Channel == "Weixin" {
		loginChannel = 1
	} else if msg.Channel == "Guest" {
		loginChannel = 3
	}
	err = dbservice.EntityUtil("Player", uid).SetValue("LoginChannel", loginChannel)
	if err != nil {
		log.Error(err, msg)
		rest.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	ret := login.UserLoginAck{
		UID:       uid,
		Token:     t,
		LobbyAddr: addr,
		Result:    0,
		ResultMsg: "登录成功",
		HB:        s.config.HBEnable,
		// Config:    s.config,
	}

	w.WriteJson(&ret)
	s.queryPFInfo(&msg)
}

// ChannelVerify 渠道验证接口
func (s *Server) channelVerify(loginMsg *login.UserLoginReq) (bool, string) {
	switch loginMsg.Channel {
	case "QQ":
		return s.doQQLogin(loginMsg.Data)
	case "Weixin":
		return s.doWXLogin(loginMsg.Data)
	case "Guest":
		return s.doGuestLogin(loginMsg.Data)
	default:
		return !s.forceChannel, ""
	}
}

// CheckWhiteList 白名单检查接口
// 级别：1表示内部 2表示外部玩家, 通过server.json里AllowGrade字段配置
func (s *Server) checkWhiteList(uid uint64) int {
	userGrade, err := dbservice.Account(uid).GetGrade()
	if err != nil {
		log.Error(err)
		return 2
	}
	if userGrade <= uint32(s.allowGrade) {
		return 1
	}

	if userGrade == 2 {
		return 2
	}

	return 0
}

func (s *Server) doQQLogin(data []byte) (bool, string) {
	msg := &QQLogin{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Error(err, string(data))
		return false, ""
	}

	timestamp := time.Now().Unix()
	sig := GenSig(timestamp)
	reqAddr := fmt.Sprintf("%s/auth/verify_login?timestamp=%d&appid=%d&sig=%s&openid=%s&encode=2",
		s.msdkAddr, timestamp, common.QQAppID, sig, msg.OpenID)
	ret := &QQLoginRet{}

	if err := s.doVerify(reqAddr, data, ret); err != nil {
		log.Error(err)
		return false, ""
	}

	if ret.Ret != 0 {
		log.Info("QQ鉴权失败", ret)
		return false, ""
	}

	return true, msg.OpenKey
}

func (s *Server) doWXLogin(data []byte) (bool, string) {
	msg := &WXLogin{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Error(err, string(data))
		return false, ""
	}

	timestamp := time.Now().Unix()
	sig := GenSig(timestamp)
	reqAddr := fmt.Sprintf("%s/auth/check_token?timestamp=%d&appid=%s&sig=%s&openid=%s&encode=2",
		s.msdkAddr, timestamp, common.WXAppID, sig, msg.OpenID)
	ret := &WXLoginRet{}

	if err := s.doVerify(reqAddr, data, ret); err != nil {
		log.Error(err)
		return false, ""
	}

	if ret.Ret != 0 {
		log.Info("微信鉴权失败", ret.Msg)
		return false, ""
	}

	return true, msg.AccessToken
}

func (s *Server) doGuestLogin(data []byte) (bool, string) {
	msg := &GuestLogin{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Error(err)
		return false, ""
	}

	timestamp := time.Now().Unix()
	sig := GenSig(timestamp)
	reqAddr := fmt.Sprintf("%s/auth/guest_check_token?timestamp=%d&appid=G_%d&sig=%s&openid=%s&encode=2",
		s.msdkAddr, timestamp, common.QQAppID, sig, msg.GuestID)
	ret := &GuestLoginRet{}

	if err := s.doVerify(reqAddr, data, ret); err != nil {
		log.Error(err)
		return false, ""
	}

	if ret.Ret != 0 {
		log.Info("游客鉴权失败", ret.Msg)
		return false, ""
	}

	return true, msg.AccessToken
}

func (s *Server) doVerify(reqAddr string, data []byte, ret interface{}) error {
	resp, err := http.Post(reqAddr, "application/x-www-form-urlencoded", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, ret); err != nil {
		return err
	}
	return nil
}

// VerifyAccountBan 验证账号是否封停
func (s *Server) VerifyAccountBan(account string) (bool, string) {
	banData := db.GetBanAccountData(account)
	if banData != nil {
		return true, banData.BanReason
	}

	return false, ""
}

// VerifyRoleBan 验证账号下的角色是否封停
func (s *Server) VerifyRoleBan(account string) (bool, string) {

	// 获取账号uid
	targetID, err := dbservice.GetUID(account)
	if err != nil || targetID == 0 {
		return false, ""
	}

	// 获取账号下的角色名称
	args := []interface{}{
		"Name",
	}
	values, err := dbservice.EntityUtil("Player", targetID).GetValues(args)
	if err != nil {
		return false, ""
	}

	name, nameErr := redis.String(values[0], nil)
	if nameErr != nil {
		return false, ""
	}

	// 判断角色是否被禁停
	data := db.GetBanRoleData(name)
	if data == nil {
		return false, ""
	}

	return true, data.BanReason
}
