package main

import (
	"common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zeus/dbservice"
	"zeus/login"

	log "github.com/cihub/seelog"
)

// WXInfo /relation/wxprofile返回用户信息数据结构
type WXInfo struct {
	NickName string `json:"nickName"` //仅返回昵称，需要备注的调用wxfriends_profile接口
	Sex      int    `json:"sex"`      //性别1男2女,用户未填则默认1男
	Picture  string `json:"picture"`  //用户头像URL,必须在URL后追加以下参数/0，/132，/96，/64，这样可以分别获得不同规格的图片：原始图片(/0)、132*132(/132)、96*96(/96)、64*64(/64)、46*46(/46)
	Provice  string `json:"provice"`  //省份
	City     string `json:"city"`     //城市
	Country  string `json:"Country"`  //国家
	Language string `json:"language"` //语言
}

// WXProfileReq /relation/wxprofile 请求格式
type WXProfileReq struct {
	AppID       string   `json:"appid"`       //游戏唯一标识
	AccessToken string   `json:"accessToken"` //登录态
	OpenIDs     []string `json:"openids"`     //需要拉取的openid账号列表
}

// WXProfileRet /relation/wxprofile 返回格式
type WXProfileRet struct {
	Ret   int      `json:"ret"` //返回码 0：正确，其它：失败
	Msg   string   `json:"msg"` //ret非0，则表示“错误码，错误提示”，详细注释参见第5节
	Lists []WXInfo `json:"lists"`
}

// QQProfileReq /relation/qqprofile 请求格式
type QQProfileReq struct {
	AppID       string `json:"appid"`       //游戏唯一标识
	OpenID      string `json:"openid"`      //用户唯一标识
	AccessToken string `json:"accessToken"` //登录态
}

// QQProfileRet /relation/qqprofile 返回格式
type QQProfileRet struct {
	Ret            int    `json:"ret"`              //返回码 0：正确，其它：失败
	Msg            string `json:"msg"`              //ret非0，则表示“错误码，错误提示”，详细注释参见第5节
	NickName       string `json:"nickName"`         //用户在QQ空间的昵称（和手机QQ昵称是同步的）
	Gender         string `json:"gender"`           //性别 如果获取不到则默认返回"男"
	Picture40      string `json:"picture40"`        //大小为40×40像素的QQ头像URL
	Picture100     string `json:"picture100"`       //大小为100×100像素的QQ头像URL，需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有
	YellowVIP      int    `json:"yellow_vip"`       //是否是黄钻用户，0表示没有黄钻
	YellowVIPLevel int    `json:"yellow_vip_level"` //黄钻等级
	YellowYearVIP  int    `json:"yellow_year_vip"`  //是否是年费黄钻用户，0表示否
	IsLost         string `json:"is_lost"`
}

// QQVipInfo QQ会员信息
type QQVipInfo struct {
	OpenID      string `json:"openid"`         //用户id
	IsQQVIP     int    `json:"is_qq_vip"`      //是否为QQ会员（0：不是； 1：是）
	QQVIPLevel  int    `json:"qq_vip_level"`   //QQ会员等级（如果是QQ会员才返回）
	IsQQYearVIP int    `json:"is_qq_year_vip"` //是否为年费QQ会员（0：不是； 1：是）
	IsQQSVIP    int    `json:"is_qq_svip"`     //是否为QQ超级会员（0：不是； 1：是）
}

// QQFriendsVIPReq /relation/qqfriends_vip 请求格式
type QQFriendsVIPReq struct {
	AppID       string   `json:"appid"`       //游戏唯一标识
	OpenID      string   `json:"openid"`      //用户唯一标识
	FOpenIDs    []string `json:"fopenids"`    //待查询openid列表，每次最多可输入50个
	Flags       string   `json:"flags"`       //VIP业务查询标识。目前支持查询QQ会员信息:qq_vip,QQ超级会员：qq_svip。后期会支持更多业务的用户VIP信息查询。如果要查询多种VIP业务，通过“,”分隔。
	UserIP      string   `json:"userip"`      //调用方ip信息
	PF          string   `json:"pf"`          //玩家登录平台，默认openmobile，有openmobile_android/openmobile_ios/openmobile_wp等，该值来自客户端手Q登录返回
	AccessToken string   `json:"accessToken"` //登录态
}

// QQFriendsVIPRet /relation/qqfriends_vip 返回格式
type QQFriendsVIPRet struct {
	Ret    int         `json:"ret"` //返回码 0：正确，其它：失败
	Msg    string      `json:"msg"` //ret非0，则表示“错误码，错误提示”
	Lists  []QQVipInfo `json:"lists"`
	IsLost string      `json:"is_lost"`
}

func (s *Server) queryPFInfo(loginMsg *login.UserLoginReq) {
	switch loginMsg.Channel {
	case "QQ":
		s.queryQQInfo(loginMsg.Data)
	case "Weixin":
		s.queryWXInfo(loginMsg.Data)
	}
}

func (s *Server) queryWXInfo(data []byte) {
	msg := &WXLogin{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Error(err, string(data))
		return
	}

	timestamp := time.Now().Unix()
	sig := GenSig(timestamp)
	reqAddr := fmt.Sprintf("%s/relation/wxprofile?timestamp=%d&appid=%s&sig=%s&openid=%s&encode=2",
		s.msdkAddr, timestamp, common.WXAppID, sig, msg.OpenID)
	req := &WXProfileReq{}
	req.AccessToken = msg.AccessToken
	req.AppID = common.WXAppID
	req.OpenIDs = []string{msg.OpenID}
	data, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := http.Post(reqAddr, "application/x-www-form-urlencoded", strings.NewReader(string(data)))
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	respInfo := &WXProfileRet{}
	err = json.Unmarshal(body, respInfo)
	if err != nil {
		log.Error(err)
		return
	}
	if respInfo.Ret != 0 {
		log.Warn("获取微信个人信息失败! [", respInfo.Msg, respInfo.Ret, "]")
		return
	}
	if len(respInfo.Lists) < 1 {
		log.Warn("个人信息列表为空")
		return
	}

	uid, err := dbservice.GetUID(msg.OpenID)
	if err != nil {
		log.Error(err)
		return
	}
	if uid == 0 {
		log.Warn("用户不存在", msg.OpenID)
		return
	}

	args := make([]interface{}, 0, 1)
	args = append(args, "Picture")
	args = append(args, []byte(respInfo.Lists[0].Picture))
	args = append(args, "NickName")
	args = append(args, []byte(respInfo.Lists[0].NickName))
	args = append(args, "Gender")
	if respInfo.Lists[0].Sex == 1 {
		args = append(args, []byte("男"))
	} else if respInfo.Lists[0].Sex == 2 {
		args = append(args, []byte("女"))
	} else {
		args = append(args, []byte("未知"))
	}
	if err = dbservice.EntityUtil("Player", uid).SetValues(args); err != nil {
		log.Error(err)
		return
	}
}

func (s *Server) queryQQInfo(data []byte) {
	msg := &QQLogin{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Error(err, string(data))
		return
	}

	s.queryQQProfile(msg)
	s.queryQQVipInfo(msg)
}

func (s *Server) queryQQProfile(msg *QQLogin) {
	timestamp := time.Now().Unix()
	sig := GenSig(timestamp)
	reqAddr := fmt.Sprintf("%s/relation/qqprofile?timestamp=%d&appid=%d&sig=%s&openid=%s&encode=2",
		s.msdkAddr, timestamp, common.QQAppID, sig, msg.OpenID)
	req := &QQProfileReq{}
	req.AccessToken = msg.OpenKey
	req.AppID = strconv.Itoa(common.QQAppID)
	req.OpenID = msg.OpenID
	data, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := http.Post(reqAddr, "application/x-www-form-urlencoded", strings.NewReader(string(data)))
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	respInfo := &QQProfileRet{}
	err = json.Unmarshal(body, respInfo)
	if err != nil {
		log.Error(err)
		return
	}
	if respInfo.Ret != 0 {
		log.Warn("获取QQ个人信息失败! [", respInfo.Msg, respInfo.Ret, "]")
		return
	}

	uid, err := dbservice.GetUID(msg.OpenID)
	if err != nil {
		log.Error(err)
		return
	}
	if uid == 0 {
		log.Warn("用户不存在", msg.OpenID)
		return
	}

	args := make([]interface{}, 0, 1)
	args = append(args, "Picture")
	args = append(args, []byte(strings.TrimSuffix(respInfo.Picture40, "/40")))
	args = append(args, "NickName")
	args = append(args, []byte(respInfo.NickName))
	args = append(args, "Gender")
	args = append(args, []byte(respInfo.Gender))
	if err = dbservice.EntityUtil("Player", uid).SetValues(args); err != nil {
		log.Error(err)
		return
	}
}

func (s *Server) queryQQVipInfo(msg *QQLogin) {
	timestamp := time.Now().Unix()
	sig := GenSig(timestamp)
	reqAddr := fmt.Sprintf("%s/relation/qqfriends_vip?timestamp=%d&appid=%d&sig=%s&openid=%s&encode=2",
		s.msdkAddr, timestamp, common.QQAppID, sig, msg.OpenID)
	req := &QQFriendsVIPReq{}
	req.AppID = strconv.Itoa(common.QQAppID)
	req.OpenID = msg.OpenID
	req.FOpenIDs = []string{msg.OpenID}
	req.Flags = "qq_vip,qq_svip"
	req.UserIP = msg.UserIP
	req.PF = "openmobile"
	req.AccessToken = msg.OpenKey
	data, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := http.Post(reqAddr, "application/x-www-form-urlencoded", strings.NewReader(string(data)))
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	respInfo := &QQFriendsVIPRet{}
	err = json.Unmarshal(body, respInfo)
	if err != nil {
		log.Error(err)
		return
	}
	if respInfo.Ret != 0 || len(respInfo.Lists) < 1 {
		log.Warn("获取QQ VIP信息失败! [", respInfo.Msg, respInfo.Ret, "]")
		return
	}

	uid, err := dbservice.GetUID(msg.OpenID)
	if err != nil {
		log.Error(err)
		return
	}
	if uid == 0 {
		log.Warn("用户不存在", msg.OpenID)
		return
	}

	var qqvip uint8
	if respInfo.Lists[0].IsQQSVIP == 1 {
		qqvip = 2
	} else if respInfo.Lists[0].IsQQVIP == 1 {
		qqvip = 1
	} else {
		qqvip = 0
	}
	if err = dbservice.EntityUtil("Player", uid).SetValue("QQVIP", qqvip); err != nil {
		log.Error(err)
		return
	}
}
