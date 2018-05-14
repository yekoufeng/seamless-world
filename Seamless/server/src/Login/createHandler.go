package main

import (
	"net/http"
	"zeus/dbservice"
	"zeus/login"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// CreateHandler 创建新帐号的接口
func (s *Server) createHandler(w rest.ResponseWriter, r *rest.Request) {
	msg := login.UserCreateReq{}
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	// 检查帐号是否存在
	var uid uint64
	var resultStr string
	if uid, err = dbservice.GetUID(msg.User); uid == 0 {
		// 用户不存在, 创建用户, 设置密码为用户输入的密码
		uid, err = s.DoCreateNewUser(msg.User, msg.Password, uint32(s.initGrade))
		if err != nil {
			log.Error(err)
			rest.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		resultStr = "创建帐号成功"
	} else if err != nil {
		log.Error(err)
		rest.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	} else {
		resultStr = "用户已存在"
	}

	ret := login.UserCreateAck{
		UID:       uid,
		Result:    0,
		ResultMsg: resultStr,
	}
	w.WriteJson(&ret)
	return
}
