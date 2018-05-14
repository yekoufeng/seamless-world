package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/pool"
	"zeus/serverMgr"
	"zeus/sess"

	log "github.com/cihub/seelog"
)

// SrvNet 建立一个服务器的网关结构
type SrvNet struct {
	srvType   uint8
	srvID     uint64
	addr      string
	outerAddr string //服务网外网地址
	token     string
	status    int
	console   uint64

	msgSrv sess.IMsgServer //负责监听其它Server连接

	pendingSesses *sync.Map
	clientSesses  *sync.Map

	srvInfo *iserver.ServerInfo

	// srvInfos *sync.Map
}

// NewSrvNet 创建一个新的服务器网
func NewSrvNet(srvType uint8, srvID uint64, addr string, outerAddr string) *SrvNet {

	if iserver.GetSrvInst() != nil {
		log.Error("Server existed")
		return nil
	}

	srv := &SrvNet{
		srvType:       srvType,
		srvID:         srvID,
		addr:          addr,
		outerAddr:     outerAddr,
		token:         "",
		msgSrv:        nil,
		pendingSesses: &sync.Map{},
		clientSesses:  &sync.Map{},

		// srvInfos: &sync.Map{},
	}

	return srv
}

func (srv *SrvNet) init() error {
	srv.msgSrv = sess.NewMsgServer("tcp", srv.addr, 0)

	srv.msgSrv.RegMsgProc(srv)

	if err := srv.msgSrv.Start(); err != nil {
		return err
	}

	if err := srv.registerSrvInfo(); err != nil {
		return err
	}

	go srv.refresh()

	return nil
}

// 10秒一次刷新过期时间, 25~35秒一次刷新一次服务器列表
func (srv *SrvNet) refresh() {
	srv.RefreshSrvInfo()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	srvNetTicker := time.NewTicker(time.Duration(5) * time.Second)
	defer srvNetTicker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := serverMgr.GetServerMgr().Update(srv.srvInfo); err != nil {
				log.Error(err)
			}
		case <-srvNetTicker.C:
			srv.RefreshSrvInfo()
		}
	}
}

func (srv *SrvNet) regMsgProc(proc interface{}) {
	if srv.msgSrv != nil {
		srv.msgSrv.RegMsgProc(proc)
	}
}

//注册当前服务器信息到redis，包括生成token等等
func (srv *SrvNet) registerSrvInfo() error {
	srv.token = srv.genToken()
	srv.regSrvInfo()
	return nil
}

func (srv *SrvNet) genToken() string {
	curtime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(curtime, 10))
	io.WriteString(h, strconv.FormatUint(srv.srvID, 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//更新服务器信息, 并刷新
func (srv *SrvNet) regSrvInfo() {
	srv.srvInfo = &iserver.ServerInfo{
		ServerID:     srv.srvID,
		Type:         srv.srvType,
		OuterAddress: srv.outerAddr,
		InnerAddress: srv.addr,
		Console:      srv.console,
		Token:        srv.token,
		Status:       srv.status,
	}

	succeed := false

	for i := 0; i < 5; i++ {
		if err := serverMgr.GetServerMgr().RegState(srv.srvInfo); err != nil {
			log.Error("regist server info fail , try again after seconds ", srv.srvInfo)
			time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)

		} else {
			succeed = true
			break
		}
	}

	if !succeed {
		log.Error("server register to srvnet failed ", srv.srvInfo)
		panic("register server failed")
	}

	return
}

// RefreshSrvInfo 刷新当前服务器信息, 最多一秒钟调用一次, 防止服务器同时拉起时大量重复请求
func (srv *SrvNet) RefreshSrvInfo() {

	remoteSrvList, err := serverMgr.GetServerMgr().GetServerList()
	if err != nil {
		log.Error("fetch server info failed", err)
		return
	}
	// srv.srvInfos = &sync.Map{}
	for _, srvInfo := range remoteSrvList {
		srv.tryConnectToSrv(srvInfo)
		// srv.srvInfos.Store(srvInfo.ServerID, srvInfo)
	}
}

func (srv *SrvNet) isNeedConnectedToSrv(info *iserver.ServerInfo) bool {
	return srv.isClientSess(info) && !srv.checkInConnectionList(info.ServerID)
}

func (srv *SrvNet) checkInConnectionList(srvID uint64) bool {

	_, ok := srv.pendingSesses.Load(srvID)
	if ok {
		log.Warn("Server contecting, waiting response", srvID)
		return true
	}

	if is, ok := srv.clientSesses.Load(srvID); ok {
		s := is.(iserver.ISess)
		if s.IsClosed() {
			log.Warnf("Server disconnected, ID:%d \n", srvID)
			srv.clientSesses.Delete(srvID)
		} else {
			return true
		}
	}

	return false
}

// GetSrvType 获取服务器类型
func (srv *SrvNet) GetSrvType() uint8 {
	return srv.srvType
}

// GetSrvID 获取服务器ID
func (srv *SrvNet) GetSrvID() uint64 {
	return srv.srvID
}

// GetSrvAddr 获取服务器内网地址
func (srv *SrvNet) GetSrvAddr() string {
	return srv.addr
}

func (srv *SrvNet) isClientSess(info *iserver.ServerInfo) bool {
	return srv.srvID < info.ServerID /*&& srv.srvType != info.Type*/
}

func (srv *SrvNet) tryConnectToSrv(info *iserver.ServerInfo) {

	if !srv.isNeedConnectedToSrv(info) {
		return
	}

	srv.pendingSesses.Store(info.ServerID, nil)
	go func() {

		s, err := sess.Dial("tcp", info.InnerAddress)
		if err != nil {
			srv.pendingSesses.Delete(info.ServerID)

			log.Errorf("Connect failed. %d to %d Addr %s, error:%v", srv.srvID, info.ServerID, info.InnerAddress, err)
			return
		}

		srv.pendingSesses.Store(info.ServerID, s)

		s.SetID(info.ServerID)
		s.SetMsgHandler(srv.msgSrv)
		s.SetServerType(info.Type)

		s.Send(&msgdef.ClientVertifyReq{
			Source: srv.srvType,
			UID:    srv.srvID,
			Token:  srv.token,
		})
		s.Start()

		log.Info("SrvNet try connect to ", info.ServerID)

	}()
}

func (srv *SrvNet) onServerConnected(srvID uint64) {
	log.Info("Connected to Server succeed !!!  ", srvID)
	log.Debug("Connected to Server succeed !!!  ", srvID)

	if _, ok := srv.clientSesses.Load(srvID); ok {
		log.Error("Session existed, server id:", srvID)
		return
	}

	sess, ok := srv.pendingSesses.Load(srvID)
	if !ok {
		log.Error("Session is pending, server id", srvID)
		return
	}

	srv.pendingSesses.Delete(srvID)
	srv.clientSesses.Store(srvID, sess)

	iserver.GetSrvInst().OnServerConnect(srvID, sess.(iserver.ISess).GetServerType())
}

func (srv *SrvNet) destroy() {
	info := &iserver.ServerInfo{
		ServerID: srv.srvID,
	}

	serverMgr.GetServerMgr().Unregister(info)
}

// MainLoop 主循环
func (srv *SrvNet) MainLoop() {
	srv.msgSrv.MainLoop()

	srv.clientSesses.Range(func(k, e interface{}) bool {
		sess := e.(iserver.ISess)
		if sess != nil {
			sess.DoMsg()
		}

		return true
	})

	srv.pendingSesses.Range(func(k, e interface{}) bool {
		if e == nil {
			return true
		}

		sess := e.(iserver.ISess)
		if sess != nil {
			sess.DoMsg()
		}
		return true
	})

}

// PostMsgToSrv 根据srvID号把消息投递到相应的服务器上
func (srv *SrvNet) PostMsgToSrv(srvID uint64, msg msgdef.IMsg) error {
	if srvID == iserver.GetSrvInst().GetSrvID() {
		srv.msgSrv.FireMsg(reflect.TypeOf(msg).Elem().Name(), msg)
		return nil
	}

	if srv.srvID < srvID {
		isess, ok := srv.clientSesses.Load(srvID)
		if ok {
			isess.(iserver.ISess).Send(msg)
			return nil
		}
	} else {
		isess := srv.msgSrv.GetSession(srvID)
		if isess != nil {
			isess.Send(msg)
			return nil
		}
	}

	return fmt.Errorf("SrvNet server %d  Server not existed, id:%d", srv.srvID, srvID)
}

// PostMsgToCell 将消息投递给某个Cell
func (srv *SrvNet) PostMsgToCell(srvID uint64, cellID uint64, msg msgdef.IMsg) error {
	srvMsg, err := srv.packSrvMsg(cellID, msg)
	if err != nil {
		return err
	}

	return srv.PostMsgToSrv(srvID, srvMsg)
}

// GetSrvIDBySrvType 获取一个SrvID
func (srv *SrvNet) GetSrvIDBySrvType(srvType uint8) (uint64, error) {
	srvInfo, err := serverMgr.GetServerMgr().GetServerByType(srvType)
	if err != nil {
		return 0, err
	}

	return srvInfo.ServerID, nil
}

func (srv *SrvNet) packSrvMsg(cellID uint64, msg msgdef.IMsg) (msgdef.IMsg, error) {
	//buf := make([]byte, sess.MaxMsgBuffer)
	buf := pool.Get(sess.MaxMsgBuffer)
	encBuf, err := sess.EncodeMsg(msg, buf, true)
	if err != nil {
		return nil, err
	}
	msgContent := make([]byte, len(encBuf))
	copy(msgContent, encBuf)
	pool.Put(buf)
	return &msgdef.SrvMsgTransport{
		CellID:     cellID,
		MsgContent: msgContent,
	}, nil
}

// GetToken 获取服务器token
func (srv *SrvNet) GetToken() string {
	return srv.token
}

// GetSrvInfo 获取服务器信息
func (srv *SrvNet) GetCurSrvInfo() *iserver.ServerInfo {
	return srv.srvInfo
}

// MsgProc_SessClosed 会话关闭
func (srv *SrvNet) MsgProc_SessClosed(content interface{}) {
	uid := content.(uint64)

	ii, ok := srv.pendingSesses.Load(uid)
	if ok {
		ii.(iserver.ISess).Close()
		srv.pendingSesses.Delete(uid)
		log.Info("SessClose, remove from pending list ", uid)
	}

	ii, ok = srv.clientSesses.Load(uid)
	if ok {
		ii.(iserver.ISess).Close()
		srv.clientSesses.Delete(uid)
		log.Info("SessClose, remove from client list ", uid)
	}

	log.Info("SrvNet server ", srv.srvID, " the connect with ", uid, " is closed")
}

//MsgProc_ClientVertifySucceedRet 连入其它服务器
func (srv *SrvNet) MsgProc_ClientVertifySucceedRet(content msgdef.IMsg) {
	msg := content.(*msgdef.ClientVertifySucceedRet)
	srv.onServerConnected(msg.SourceID)

	//seelog.Debug("SrvNet server ", srv.srvID, " connect to sserver ", msg.SourceID, "  succeed")
}

//MsgProc_SessVertified 连接已经验证通过
func (srv *SrvNet) MsgProc_SessVertified(content interface{}) {

	uid := content.(uint64)

	sess := srv.msgSrv.GetSession(uid)
	sess.Send(&msgdef.ClientVertifySucceedRet{
		Source:   0,
		UID:      0,
		SourceID: srv.srvID,
		Type:     0,
	})

	log.Info("SrvNet server ", srv.srvID, "  recevice a connect from server ", uid)
	//触发建立连接回调
	iserver.GetSrvInst().OnServerConnect(sess.GetID(), sess.GetServerType())
}
