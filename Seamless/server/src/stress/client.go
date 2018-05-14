package main

import (
	"common"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"protoMsg"
	"strings"
	"time"
	"zeus/iserver"
	"zeus/login"
	"zeus/msgdef"
	"zeus/serializer"
	"zeus/sess"
	"zeus/space"

	"github.com/spf13/viper"
)

type Client struct {
	msgClient     iserver.ISess
	msgRoomClient iserver.ISess

	// loginC          chan bool
	enterGatewaySig chan bool
	enterRoomSig    chan bool
	matchSuccSig    chan bool
	enterSpaceSig   chan bool
	inRoom          bool
	inGame          bool

	gateway    string
	verifyMsg  *msgdef.ClientVertifyReq
	verifyTime time.Time
	queueTime  time.Time

	retChan chan *TranRet

	tableID uint64
	uid     uint64

	enterSpaceMsg *msgdef.EnterCell
	states        *space.EntityStates
	curState      *RoomPlayerState
}

// NewClient 新建客户端
func NewClient() *Client {
	c := &Client{}
	// c.loginC = make(chan bool, 1)
	c.enterGatewaySig = make(chan bool, 1)
	c.enterRoomSig = make(chan bool, 1)
	c.matchSuccSig = make(chan bool, 1)
	c.enterSpaceSig = make(chan bool, 1)

	c.states = space.NewEntityStates()
	c.curState = NewRoomPlayerState().(*RoomPlayerState)
	return c
}

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

func (c *Client) SetTranChan(retChan chan *TranRet) {
	c.retChan = retChan
}

func (c *Client) Run(user, pwd string) {
	stopC := make(chan os.Signal, 1)
	// c.loginC = make(chan bool, 1)
	c.enterGatewaySig = make(chan bool, 1)
	c.enterRoomSig = make(chan bool, 1)

	c.login(user, pwd)
	for {
		select {
		case <-stopC:
			return
		// case <-c.loginC:
		//c.enterGateway()
		case <-c.enterGatewaySig:
			//c.RPCCall(11, 1000, "SetRole", uint64(3))
			//c.doQueue()
		case <-c.enterRoomSig:
			//c.startGame()
		}
	}

}

func doMatch(ctx context.Context, user, pwd string, retChan chan *TranRet, team bool) {
	client := NewClient()

	// Login 事务
	loginRet := &TranRet{}
	loginRet.Action = "Login"
	start := time.Now()
	loginRet.Err = client.login(user, pwd)
	loginRet.Cost = time.Now().Sub(start).Nanoseconds() / 1000000
	retChan <- loginRet
	if loginRet.Err != nil {
		return
	}

	// UserVerity 事务
	// <-client.loginC
	verifyRet := &TranRet{}
	verifyRet.Action = "UserVerity"
	verifyRet.Err = client.enterGateway()
	if verifyRet.Err != nil {
		verifyRet.Cost = 0
		retChan <- verifyRet
	}

	for {
		client.inMatchAction(ctx, retChan, team)

		interval := viper.GetInt("Match.Interval")
		if interval == 0 {
			interval = 5
		}
		time.Sleep(time.Duration(interval) * time.Second)

		if team {
			client.doTeamQueue()
		} else {
			client.doQueue(uint32(1))
		}
	}
}

func (c *Client) inMatchAction(ctx context.Context, retChan chan *TranRet, team bool) {
	ticker := time.NewTicker(33 * time.Millisecond)
	c.SetTranChan(retChan)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.msgClient != nil {
				c.msgClient.DoMsg()
			}
		case <-c.enterSpaceSig:
		case <-c.matchSuccSig:
			return
		case <-c.enterGatewaySig:
			time.Sleep(5 * time.Second)
			if team {
				c.doTeamQueue()
			} else {
				c.doQueue(uint32(1))
			}
		}
	}
}

func doGameRepeat(ctx context.Context, user, pwd string, retChan chan *TranRet, typ int, mapid uint32) {
	client := NewClient()

	// Login 事务
	loginRet := &TranRet{}
	loginRet.Action = "Login"
	start := time.Now()
	loginRet.Err = client.login(user, pwd)
	loginRet.Cost = time.Now().Sub(start).Nanoseconds() / 1000000
	retChan <- loginRet
	if loginRet.Err != nil {
		return
	}

	// UserVerity 事务
	// <-client.loginC
	verifyRet := &TranRet{}
	verifyRet.Action = "UserVerity"
	verifyRet.Err = client.enterGateway()
	if verifyRet.Err != nil {
		verifyRet.Cost = 0
		retChan <- verifyRet
	}

	for {
		if client.inGameAction(ctx, retChan, typ, mapid) {
			if client.msgClient != nil {
				client.msgClient.Close()
			}
			return
		}

		time.Sleep(15 * time.Second)

		client.doQueue(mapid)
	}
}

func doGame(ctx context.Context, user, pwd string, retChan chan *TranRet, typ int, mapid uint32) {
	client := NewClient()

	// Login 事务
	loginRet := &TranRet{}
	loginRet.Action = "Login"
	start := time.Now()
	loginRet.Err = client.login(user, pwd)
	loginRet.Cost = time.Now().Sub(start).Nanoseconds() / 1000000
	retChan <- loginRet
	if loginRet.Err != nil {
		return
	}

	// UserVerity 事务
	// <-client.loginC
	verifyRet := &TranRet{}
	verifyRet.Action = "UserVerity"
	verifyRet.Err = client.enterGateway()
	if verifyRet.Err != nil {
		verifyRet.Cost = 0
		retChan <- verifyRet
	}

	ret := client.inGameAction(ctx, retChan, typ, mapid)
	if client.msgClient != nil {
		fmt.Println("Done, close sess", ret)
		client.msgClient.Close()
	}
}

func (c *Client) inGameAction(ctx context.Context, retChan chan *TranRet, typ int, mapid uint32) bool {
	ticker := time.NewTicker(33 * time.Millisecond)
	deadTime := rand.Intn(1700) + 100
	closeTime := deadTime + 5
	deadTimer := time.NewTimer(time.Duration(deadTime) * time.Second)
	closeTimer := time.NewTimer(time.Duration(closeTime) * time.Second)
	c.inGame = false
	c.inRoom = false
	c.SetTranChan(retChan)
	for {
		select {
		case <-ctx.Done():
			return true
		case <-deadTimer.C:
			c.RPCCall(common.ServerTypeRoom, 0, "GmDeath", uint32(0))
			c.inGame = false
		case <-closeTimer.C:
			return false
		case <-ticker.C:
			if c.msgClient != nil {
				c.msgClient.DoMsg()
			}
			if !c.inRoom && c.msgRoomClient != nil {
				c.msgRoomClient.DoMsg()
			}
			if c.inGame {
				c.curState.TimeStamp++
				// if typ == 1 {
				c.sendUserMove()
				// } else if typ == 2 {
				// 	c.sendShootReq()
				// } else if typ == 3 {
				// 	c.sendAttackReq()
				// } else if typ == 4 {
				// 	c.sendPickupItem()
				// }
			}
		case <-c.enterGatewaySig:
			queueTimer := time.NewTimer(15 * time.Second)
			select {
			case <-queueTimer.C:
				c.doQueue(mapid)
			}
		case <-c.matchSuccSig:
		case <-c.enterSpaceSig:
			if err := c.enterRoom(); err != nil {
				fmt.Println(err)
			}
		case <-c.enterRoomSig:
			c.inRoom = true
			msg := &msgdef.SpaceUserConnect{}
			msg.UID = c.uid
			msg.SpaceID = c.enterSpaceMsg.SpaceID
			c.msgRoomClient.Send(msg)
			c.msgRoomClient.SetMsgHandler(c.msgClient)
			c.msgRoomClient.Start()
			// fmt.Println("EnterRoom Succ")
		}
	}
}

func doLoginAll(ctx context.Context, user, pwd string, retChan chan *TranRet) {
	client := NewClient()

	// Login 事务
	loginRet := &TranRet{}
	loginRet.Action = "Login"
	start := time.Now()
	loginRet.Err = client.login(user, pwd)
	loginRet.Cost = time.Now().Sub(start).Nanoseconds() / 1000000
	retChan <- loginRet
	if loginRet.Err != nil {
		return
	}

	// UserVerity 事务
	// <-client.loginC
	client.SetTranChan(retChan)
	verifyRet := &TranRet{}
	verifyRet.Action = "UserVerity"
	verifyRet.Err = client.enterGateway()
	if verifyRet.Err == nil {
		ticker := time.NewTicker(33 * time.Millisecond)
		loginTime := viper.GetInt("Login.LoginTime")
		timer := time.NewTimer(time.Duration(loginTime) * time.Second)
		defer func() {
			if client.msgClient != nil {
				client.msgClient.Close()
			}

			ticker.Stop()
			timer.Stop()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				return
			case <-ticker.C:
				if client.msgClient != nil {
					client.msgClient.DoMsg()
				}
			case <-client.enterGatewaySig:
			}
		}
	} else {
		verifyRet.Cost = 0
		retChan <- verifyRet
	}
}

func doLogin(ctx context.Context, user, pwd string, retChan chan *TranRet) {
	client := NewClient()

	// Login 事务
	loginRet := &TranRet{}
	loginRet.Action = "Login"
	start := time.Now()
	loginRet.Err = client.login(user, pwd)
	loginRet.Cost = time.Now().Sub(start).Nanoseconds() / 1000000
	retChan <- loginRet
	if loginRet.Err != nil {
		return
	}

	// UserVerity 事务
	// <-client.loginC
	client.SetTranChan(retChan)
	verifyRet := &TranRet{}
	verifyRet.Action = "UserVerity"
	verifyRet.Err = client.enterGateway()
	if verifyRet.Err == nil {
		ticker := time.NewTicker(33 * time.Millisecond)
		loginTime := viper.GetInt("Login.LoginTime")
		timer := time.NewTimer(time.Duration(loginTime) * time.Second)
		defer func() {
			if client.msgClient != nil {
				client.msgClient.Close()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				return
			case <-ticker.C:
				if client.msgClient != nil {
					client.msgClient.DoMsg()
				}
			case <-client.enterGatewaySig:
			}
		}
	} else {
		verifyRet.Cost = 0
		retChan <- verifyRet
	}
}

func doLoginRepeat(ctx context.Context, user, pwd string, retChan chan *TranRet) {
	for {
		doLogin(ctx, user, pwd, retChan)
		interval := viper.GetInt("Login.ReloginInterval")
		timer := time.NewTimer(time.Duration(interval) * time.Second)
		select {
		case <-timer.C:
			fmt.Println("Relogin, go")
		}
	}
}

func doRegister(user, pwd string, retChan chan *TranRet) {
	tranRet := &TranRet{}
	tranRet.Action = "CreateAccnt"
	tranRet.Err = nil
	var start time.Time
	defer func() {
		tranRet.Cost = time.Now().Sub(start).Nanoseconds() / 1000000
		retChan <- tranRet
	}()

	regMsg := login.UserCreateReq{
		User:     user,
		Password: pwd,
	}
	data, err := json.Marshal(regMsg)
	if err != nil {
		fmt.Println(err)
		tranRet.Err = err
		return
	}

	addr := fmt.Sprintf("http://%s/create", viper.GetString("Config.Addr"))
	start = time.Now()
	resp, err := http.Post(addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		tranRet.Err = err
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		tranRet.Err = err
		return
	}
	retMsg := login.UserCreateAck{}
	json.Unmarshal(body, &retMsg)
	if retMsg.Result != 0 {
		fmt.Println("注册失败", retMsg.ResultMsg)
	} else {
		fmt.Println("注册成功", retMsg)
	}
	return
}

func (c *Client) login(user, pwd string) error {
	// HTTP
	loginMsg := UserLoginReq{user, pwd}
	data, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	addr := fmt.Sprintf("http://%s/login", viper.GetString("Config.Addr"))

	resp, err := http.Post(addr, "application/json", strings.NewReader(string(data)))
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

	// 保存验证信息
	c.gateway = resMsg.LobbyAddr
	c.verifyMsg = new(msgdef.ClientVertifyReq)
	c.verifyMsg.Source = msgdef.ClientMSG
	c.verifyMsg.UID = resMsg.UID
	c.verifyMsg.Token = resMsg.Token
	c.uid = c.verifyMsg.UID

	// c.loginC <- true
	return nil
}

func (c *Client) enterGateway() error {
	if c.verifyMsg == nil {
		fmt.Println("没有Token")
		return fmt.Errorf("没有Token")
	}
	//common.InitMsg()

	c.verifyTime = time.Now()
	var err error

	c.msgClient, err = sess.Dial("tcp", c.gateway)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Connected to", c.gateway)

	// c.msgClient.Close()

	c.msgClient.RegMsgProc(&ClientMsgProc{c})
	fmt.Println(c.verifyMsg)
	c.msgClient.Send(c.verifyMsg)

	c.msgClient.Start()

	return nil
}

func (c *Client) enterRoom() error {
	if c.verifyMsg == nil {
		fmt.Println("没有Token")
		return fmt.Errorf("没有Token")
	}
	if c.enterSpaceMsg == nil {
		fmt.Println("没有Room信息")
		return fmt.Errorf("没有Room信息")
	}

	var err error

	c.msgRoomClient, err = sess.Dial("tcp", c.enterSpaceMsg.Addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.msgRoomClient.RegMsgProc(&ClientRoomMsgProc{c})
	c.msgRoomClient.Send(c.verifyMsg)
	fmt.Println("Connect to", c.enterSpaceMsg.Addr)

	return nil
}

func (c *Client) doQueue(mapid uint32) {
	c.RPCCall(common.ServerTypeLobby, 0, "EnterRoomReq", mapid)
	c.queueTime = time.Now()
	fmt.Println("Start Match")
}

func (c *Client) doTeamQueue() {

	c.RPCCall(common.ServerTypeLobby, uint64(0), "QuickEnterTeam", uint32(1), uint32(c.uid%2), uint32(c.uid%2))

	time.Sleep(1 * time.Second)

	c.RPCCall(common.ServerTypeLobby, uint64(0), "ConfirmTeamMatch")

	c.queueTime = time.Now()
	fmt.Println("Start Team match")
}

func (c *Client) setBaseState(bs uint8) {
	if c.msgRoomClient == nil {
		fmt.Println("msgRoomClient nil")
		return
	}

	ls := c.curState.Clone()
	c.curState.BaseState = bs
	data, ok := ls.Delta(c.curState)
	if !ok {
		msg := &msgdef.SyncUserState{}
		msg.Data = data
		c.msgRoomClient.Send(msg)
	}
}

func (c *Client) sendUserMove() {
	if c.msgRoomClient == nil {
		fmt.Println("msgRoomClient nil")
		return
	}

	pos := c.curState.GetPos()
	step := rand.Float32() * 0.1
	pos.X += step
	pos.Y += step

	ls := c.curState.Clone()
	c.curState.SetPos(pos)
	data, ok := ls.Delta(c.curState)
	if !ok {
		msg := &msgdef.SyncUserState{}
		msg.Data = data
		c.msgRoomClient.Send(msg)
	}
}

func (c *Client) sendShootReq() {
	if c.msgRoomClient == nil {
		fmt.Println("msgClient nil")
		return
	}

	msg := &protoMsg.ShootReq{}
	msg.Issuc = true
	msg.Attackid = c.uid

	c.RoomRPCCall(common.ServerTypeRoom, 0, "ShootReq", msg)
	// c.msgRoomClient.Send(msg)
}

func (c *Client) sendAttackReq() {
	if c.msgRoomClient == nil {
		fmt.Println("msgRoomClient nil")
		return
	}

	msg := &protoMsg.AttackReq{}
	msg.Defendid = uint64(rand.Intn(50000))
	msg.Ishead = false
	msg.Origion = &protoMsg.Vector3{rand.Float32() * 1000, rand.Float32() * 50, rand.Float32() * 1000}
	msg.Dir = &protoMsg.Vector3{0, rand.Float32(), 0}
	msg.Firetime = 1
	msg.Distance = 100
	msg.WheelIndex = 0

	c.RoomRPCCall(common.ServerTypeRoom, 0, "AttackReq", msg)
	// c.msgClient.Send(msg)
}

func (c *Client) sendPickupItem() {
	c.RoomRPCCall(common.ServerTypeRoom, c.uid, "PickupItem", uint64(rand.Intn(40000)))
}

func (c *Client) startGame() {

}

func (c *Client) doBalance() {

}

func (c *Client) RPCCall(srvType uint8, srcEntityID uint64, methodName string, args ...interface{}) {
	if c.msgClient == nil {
		fmt.Println("msgClient nil")
		return
	}

	data := serializer.Serialize(args...)

	msg := &msgdef.RPCMsg{}
	msg.ServerType = srvType
	msg.SrcEntityID = srcEntityID
	msg.MethodName = methodName
	msg.Data = data

	c.msgClient.Send(msg)
	return
}

func (c *Client) RoomRPCCall(srvType uint8, srcEntityID uint64, methodName string, args ...interface{}) {
	if c.msgRoomClient == nil {
		fmt.Println("msgRoomClient nil")
		return
	}

	data := serializer.Serialize(args...)

	msg := &msgdef.RPCMsg{}
	msg.ServerType = srvType
	msg.SrcEntityID = srcEntityID
	msg.MethodName = methodName
	msg.Data = data

	c.msgRoomClient.Send(msg)
	return
}
