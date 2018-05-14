package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	//"crypto/hmac"
	//"crypto/sha1"
	"net/url"
	"sort"

	log "github.com/cihub/seelog"
)

const (
	MidasAppKey = "Mhpn1vP6EEdrOVzlVp2zWs4io5bn6vMM"
)

const (
	// MSDK KEY
	MSDK_KEY = "7eead96a3fdb063615b181d7c01480e4"
)

type SigPara []string

func (a SigPara) Len() int           { return len(a) }
func (a SigPara) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SigPara) Less(i, j int) bool { return a[i] < a[j] }

func (srv *Server) getSig(urlPath string, httpReqType string, para SigPara) string {

	// 构建源串
	// 第1步
	urlPathEncode := url.QueryEscape(urlPath)

	// 第2步
	sort.Sort(para)
	paraStr := ""
	for i, j := range para {
		if i != 0 {
			paraStr += "&"
		}
		paraStr += j
	}

	// 第3步
	paraEncode := url.QueryEscape(paraStr)

	// 第4步
	reqStr := httpReqType + "&" + urlPathEncode + "&" + paraEncode

	// 构造秘钥
	secretKey := MidasAppKey + "&"

	// 生成签名值
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(reqStr))

	s := []byte(mac.Sum(nil))
	uEnc := base64.StdEncoding.EncodeToString(s)

	urlSig := url.QueryEscape(uEnc)

	return urlSig
}

// GetpfAndpfKey 获取pf和pfKey
func (srv *Server) GetpfAndpfKey(data *GetpfAndpfkeyData, curTime string) (string, string, error) {

	log.Debugf("GetpfAndpfKey data:%+v\n", data)

	cmdData, DataErr := json.Marshal(data)
	if DataErr != nil {
		log.Error("json.Marshal err:", DataErr)
		return "", "", DataErr
	}

	h := md5.New()
	sigPara := MSDK_KEY + curTime

	io.WriteString(h, sigPara)
	token := fmt.Sprintf("%x", h.Sum(nil))

	client := &http.Client{}
	urlinfo := "http://msdktest.qq.com/auth/get_pfval?" +
		"timestamp=" + curTime +
		"&appid=" + data.Appid +
		"&sig=" + token +
		"&openid=" + data.Openid +
		"&encode=2"
	req, reqErr := http.NewRequest("POST", urlinfo, strings.NewReader(string(cmdData)))
	if reqErr != nil {
		log.Error("GetpfAndpfKey NewRequest err: ", reqErr)
		return "", "", reqErr
	}

	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Error("GetpfAndpfKey err:", respErr)
		return "", "", respErr
	}

	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Errorf("GetpfAndpfKey ret ReadAll err:%+v", bodyErr)
		return "", "", bodyErr
	}

	respInfo := &GetpfAndpfkeyRet{}
	err := json.Unmarshal(body, respInfo)
	if err != nil {
		log.Error("GetpfAndpfKey ret err: ", err)
		log.Error(err)
		return "", "", err
	}

	log.Infof("GetpfAndpfKey result :%+v \n", respInfo)

	if respInfo.Ret != 0 {
		return "", "", errors.New("getpf fail")
	}

	return respInfo.Pf, respInfo.PfKey, nil
}
