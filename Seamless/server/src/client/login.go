package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"protoMsg"
	"strings"
	"sync/atomic"
	"time"
	"zeus/msgdef"
	"zeus/sess"

	"github.com/spf13/viper"
)

var username = flag.String("username", "", "Use -username <username>")
var httpport = flag.String("httpport", "", "Use -httpport <httpport>")

// UserLoginReq 登录消息格式
type UserLoginReq struct {
	User     string
	Password string
}

// UserLoginAck 登录消息返回格式
type UserLoginAck struct {
	UID       uint64
	Token     string
	LobbyAddr string
	Result    int
	ResultMsg string
}

func (cli *Client) Login() error {
	loginMsg := UserLoginReq{cli.GetUserName(), ""}
	data, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	addr := viper.GetString("Login.Addr") + ":" + viper.GetString("Login.Port")
	url := "http://" + addr + "/login"
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	resMsg := UserLoginAck{}
	json.Unmarshal(body, &resMsg)

	//判断返回消息有效性
	if resMsg.LobbyAddr == "" {
		fmt.Println("LobbyAddr无效")
		return err
	}
	if resMsg.Token == "" {
		fmt.Println("Token无效")
		return err
	}

	fmt.Println(resMsg)

	cli.UID = resMsg.UID
	cli.Token = resMsg.Token
	cli.LobbyAddr = resMsg.LobbyAddr

	msgClient, err := sess.Dial("tcp", resMsg.LobbyAddr)
	if err != nil {
		fmt.Println(err, resMsg.LobbyAddr)
		return err
	}
	cli.sess = msgClient
	msgClient.RegMsgProc(&ClientMsgProc{})
	msgClient.Start()
	// go cli.DoMsg()

	verifyMsg := new(msgdef.ClientVertifyReq)
	verifyMsg.Source = msgdef.ClientMSG
	verifyMsg.UID = resMsg.UID
	verifyMsg.Token = resMsg.Token
	msgClient.Send(verifyMsg)

	return nil
}

func (cli *Client) DetectCell() {
	now := time.Now().UnixNano() / 1e6
	if now-cli.detectCellTime > 1000 {
		cli.SendMsgToCell(&protoMsg.DetectCell{})
		cli.detectCellTime = now
	}
}

func (cli *Client) SetCellOk() {
	atomic.StoreInt64(&cli.cellOkConfirmTime, time.Now().UnixNano()/1e6)
}

func (cli *Client) IsCellOk() bool {
	return time.Now().UnixNano()/1e6-atomic.LoadInt64(&cli.cellOkConfirmTime) < 2000
}
