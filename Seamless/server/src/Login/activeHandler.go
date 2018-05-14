package main

import (
	"db"
	"excel"
	"strings"
	"zeus/dbservice"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// AccountActive 账号激活
func (s *Server) AccountActive(w rest.ResponseWriter, r *rest.Request) {

	msg := AccountActiveReq{}
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", 400)
		return
	}

	log.Info("请求账号激活 ", msg.Activatetype, msg.User, msg.Password, msg.CDKEY)

	ret := AccountActiveRet{
		RetCode: ACTIVECDKEYINUSE,
	}

	result := s.ChannelVerifyCDKEY(&msg)
	if result == db.CDKEYEXIST {
		// 判断账号是否已经存在
		if uid, _ := dbservice.GetUID(msg.User); uid != 0 {
			ret.RetCode = ACTIVEACCOUNTEXIST
		} else {
			_, err = s.App.DoCreateNewUser(msg.User, msg.Password, uint32(s.initGrade))
			if err != nil {
				log.Error(err)
				rest.Error(w, "服务器错误", 500)
				return
			}
			ret.RetCode = ACTIVESUCCUESS

			if msg.Activatetype == 0 { // 手机号激活
				db.SetCDKEYUse(msg.User, msg.User)
			} else if msg.Activatetype == 1 { // 激活码激活
				db.SetCDKEYUse(msg.CDKEY, msg.User)
			}
		}
	} else if result == db.CDKEYINUSE {
		ret.RetCode = ACTIVECDKEYINUSE
	} else if result == db.CDKEYABSEND {
		ret.RetCode = ACTIVEABSEND
	}

	w.WriteJson(&ret)

	log.Info("账号激活结果 ", msg.Activatetype, msg.User, msg.Password, msg.CDKEY, ret.RetCode)
}

// ChannelVerifyCDKEY 验证激活码
func (s *Server) ChannelVerifyCDKEY(msg interface{}) int {

	info, ok := msg.(*AccountActiveReq)
	if !ok {
		return db.CDKEYABSEND
	}

	if info.Activatetype == 0 { // 手机号激活
		if db.JudgeCDKEYInUse(info.User) {
			return db.CDKEYINUSE
		}

		for _, accData := range excel.GetPhonenumberMap() {
			if accData.Phonenumber == info.User {
				return db.CDKEYEXIST
			}
		}
	} else if info.Activatetype == 1 { // 激活码激活
		if db.JudgeCDKEYInUse(info.CDKEY) {
			return db.CDKEYINUSE
		}

		for _, accData := range excel.GetAdkeyMap() {
			if accData.Adkey == info.CDKEY {
				return db.CDKEYEXIST
			}
		}
	}

	return db.CDKEYABSEND
}

// cdkeyVerify 激活码验证
func (s *Server) cdkeyVerify(w rest.ResponseWriter, r *rest.Request) {

	msg := AccountActiveReq{}
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", 400)
		return
	}
	msg.CDKEY = strings.ToUpper(msg.CDKEY)

	ret := AccountActiveRet{
		RetCode: ACTIVECDKEYINUSE,
	}

	result := db.GetCDKEYState(msg.CDKEY)
	if result == db.CDKEYEXIST {

		// 判断账号是否已经存在
		if uid, _ := dbservice.GetUID(msg.User); uid != 0 {

			if db.SetCDKEYBeUsed(msg.CDKEY, msg.User) {
				dbservice.Account(uid).SetGrade(2)
				ret.RetCode = ACTIVESUCCUESS
			}

		}

	} else if result == db.CDKEYINUSE {
		ret.RetCode = ACTIVECDKEYINUSE
	} else if result == db.CDKEYABSEND {
		ret.RetCode = ACTIVEABSEND
	}

	w.WriteJson(&ret)

	log.Info("cdkeyVerify ", msg, ret.RetCode)
}
