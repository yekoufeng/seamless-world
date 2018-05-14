package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// ServerCmd /cmds处理函数
func (srv *AdminServer) ServerCmd(w rest.ResponseWriter, r *rest.Request) {
	var cmd struct {
		ServerID uint64
		Command  string
	}
	err := r.DecodeJsonPayload(&cmd)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", 400)
		return
	}

	srv.serverMapRW.RLock()
	targerServer, ok := srv.serverMap[cmd.ServerID]
	srv.serverMapRW.RUnlock()
	if !ok {
		log.Error("目标服务器不存在")
		rest.Error(w, "目标服务器不存在", 400)
		return
	}

	// 转发命令至目标服务器
	var innerCmd struct {
		ServerID uint64
		Command  string
		Token    string
	}
	innerCmd.ServerID = cmd.ServerID
	innerCmd.Command = cmd.Command
	innerCmd.Token = targerServer.Token
	cmdData, err := json.Marshal(innerCmd)
	if err != nil {
		log.Error(err)
		rest.Error(w, "服务器内部错误", 500)
		return
	}
	ip := strings.Split(targerServer.InnerAddress, ":")[0]
	targerAddr := fmt.Sprintf("http://%s:%d/exec", ip, targerServer.Console)

	client := &http.Client{}
	req, err := http.NewRequest("POST", targerAddr, strings.NewReader(string(cmdData)))
	if err != nil {
		log.Error(err)
		rest.Error(w, "服务器内部错误", 500)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", r.Header.Get("Origin"))
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		rest.Error(w, "服务器内部错误", 500)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		rest.Error(w, "服务器内部错误", 500)
		return
	}
	w.(http.ResponseWriter).Write(body)
}
