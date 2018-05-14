package sess

import (
	"zeus/iserver"
	"zeus/msghandler"
)

// IMsgServer 消息服务器接口
type IMsgServer interface {
	msghandler.IMsgHandlers
	IConnReceiver

	Start() error
	Close()
	MainLoop()
	GetSession(id uint64) iserver.ISess
	Count() int
	SetEncryptEnabled()
}

// NewMsgServer 创建一个消息监听器
func NewMsgServer(protocal string, addr string, maxConns int) IMsgServer {
	srv := new(MsgServer)

	mh := msghandler.NewMsgHandlers()

	srv.IMsgHandlers = mh
	srv.SessMgr = newSessMgr(mh)
	srv.INetSrv = newNetSrv(protocal, srv, maxConns)
	srv.listenAddr = addr
	srv.maxConn = maxConns

	return srv
}

// MsgServer 消息服务器
type MsgServer struct {
	msghandler.IMsgHandlers
	*SessMgr
	INetSrv
	maxConn    int
	listenAddr string
}

// Start 启动服务器
func (srv *MsgServer) Start() error {
	return srv.INetSrv.start(srv.listenAddr)
}

// Close 关闭服务器
func (srv *MsgServer) Close() {
	srv.INetSrv.close()
	srv.SessMgr.close()
}

// SetEncryptEnabled 是否加密
func (srv *MsgServer) SetEncryptEnabled() {
	srv.SessMgr.SetEncryptEnabled()
}

// GetSession 获取一个Sess对象
func (srv *MsgServer) GetSession(id uint64) iserver.ISess {
	return srv.SessMgr.getSessByID(id)
}

// MainLoop 每帧调用
func (srv *MsgServer) MainLoop() {
	srv.msgHandlers.DoMsg()
}
