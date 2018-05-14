package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	//"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// QueryBalance 查询剩余虚拟币
func (srv *Server) QueryBalance(w rest.ResponseWriter, r *rest.Request) {

	log.Infof("QueryBalance start : %s\n", time.Now().Format("2006-01-02 15:04:05"))

	clientMsg := QueryBalaceByClient{}
	err := r.DecodeJsonPayload(&clientMsg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "解析json出错", ErrPayJSONDecodeFailed)
		return
	}
	clientMsg.Offerid = GetMidasAppidByType(clientMsg.Os)

	log.Infof("QueryBalance request para: %+v\n", clientMsg)

	client := &http.Client{}

	nowTime := time.Now().Unix()
	nowStr := strconv.FormatInt(nowTime, 10)

	getpfAndpfkey := &GetpfAndpfkeyData{
		Appid:          GetMSDKAppidByType(clientMsg.Platform), // 游戏的唯一标识
		Openid:         clientMsg.Openid,                       // 用户的唯一标识
		AccessToken:    clientMsg.AccessToken,                  // AccessToken
		Platform:       clientMsg.Platform,                     // 平台标识(一般情况下：qq对应值为desktop_m_qq，wx对应值为desktop_m_wx)
		RegChannel:     clientMsg.RegChannel,                   // 注册渠道
		Os:             clientMsg.Os,                           // 系统(安卓对应android，ios对应iap)
		Installchannel: clientMsg.Installchannel,               // 安装渠道
		Offerid:        clientMsg.Offerid,                      // 支付的appid
	}

	// 获取pf和pfkey
	Pf, Pfkey, pfErr := srv.GetpfAndpfKey(getpfAndpfkey, nowStr)
	if err != nil {
		log.Error("QueryBalance GetpfAndpfKey info:", pfErr)
		rest.Error(w, "获取pf和pfkey失败", ErrPayGetPfAndPfkeyFailed)
		return
	}

	// 获取sig
	var sigpara SigPara
	urlPath := "/mpay/get_balance_m"
	httpReqType := "GET"

	sigpara = append(sigpara, "openid="+clientMsg.Openid)
	sigpara = append(sigpara, "openkey="+clientMsg.Openkey)
	sigpara = append(sigpara, "pay_token="+clientMsg.PayToken)
	sigpara = append(sigpara, "appid="+GetMidasAppidByType(clientMsg.Os))
	sigpara = append(sigpara, "ts="+nowStr)
	sigpara = append(sigpara, "pf="+Pf)
	sigpara = append(sigpara, "pfkey="+Pfkey)
	sigpara = append(sigpara, "zoneid="+"1")
	sigpara = append(sigpara, "accounttype="+"common")
	Sig := srv.getSig(urlPath, httpReqType, sigpara)

	urlinfo := "http://opensdktest.tencent.com/mpay/get_balance_m?" +
		"openid=" + clientMsg.Openid +
		"&openkey=" + clientMsg.Openkey +
		"&pay_token=" + clientMsg.PayToken +
		"&appid=" + GetMidasAppidByType(clientMsg.Os) +
		"&ts=" + nowStr +
		"&sig=" + Sig +
		"&pf=" + Pf +
		"&pfkey=" + Pfkey +
		"&zoneid=" + "1" +
		"&accounttype=" + "common"

	log.Info("urlinfo :", urlinfo)

	req, reqErr := http.NewRequest("GET", urlinfo, nil)
	if reqErr != nil {
		log.Error("QueryBalance NewRequest err: ", err)

		rest.Error(w, "生成向midas请求失败", ErrPayNewRequestFailed)
		return
	}

	c1 := &http.Cookie{
		Name:     "session_id",
		Value:    "openid",
		HttpOnly: true,
	}
	c2 := &http.Cookie{
		Name:     "session_type",
		Value:    "kp_actoken",
		HttpOnly: true,
	}
	c3 := &http.Cookie{
		Name:     "org_loc",
		Value:    "/mpay/get_balance_m",
		HttpOnly: true,
	}
	/*
		c4 := &http.Cookie{
			Name:     "appip",
			Value:    "",
			HttpOnly: true,
		}
	*/

	req.AddCookie(c1)
	req.AddCookie(c2)
	req.AddCookie(c3)
	//req.AddCookie(c4)

	log.Errorf("Req get_balance_m info :%v", req)
	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Error("/mpay/get_balance_m client.Do() err:", respErr)
		rest.Error(w, "请求midas失败", ErrPayReqMidasFailed)
		return
	}

	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Error("QueryBalance1 ret ReadAll err: ", bodyErr)
		rest.Error(w, "读取midas返回数据失败", ErrPayReadMidasRetFailed)
		return
	}

	respInfo := &MidasQueryBalanceResult{}
	unErr := json.Unmarshal(body, respInfo)
	if unErr != nil {
		log.Error("querybalance midas ret data err: ", unErr)
		rest.Error(w, "读取midas返回数据失败", ErrPayReadMidasRetFailed)
		return
	}

	log.Errorf("result=%+v", respInfo)

	retClientMsg := &RetQueryBalanceResultToClient{
		Ret:     respInfo.Ret,
		Balance: respInfo.Balance,
	}

	log.Debug("req query balance result = %+v", retClientMsg)
	w.WriteJson(retClientMsg)
}
