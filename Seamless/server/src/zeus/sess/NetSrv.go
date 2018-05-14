package sess

import (
	"context"
	"net"

	log "github.com/cihub/seelog"
	"golang.org/x/net/netutil"
)

func newNetSrv(protocal string, msgDeliver IConnReceiver, maxConns int) INetSrv {
	if protocal == "tcp" {
		return newARQNetSrv(protocal, msgDeliver, maxConns)
	}
	return nil
}

// INetSrv 网络监听
type INetSrv interface {
	start(addr string) error
	close()
}

// IConnReceiver 连接接收
type IConnReceiver interface {
	acceptConn(net.Conn)
}

///////////////////////////////////////////////////////////////////

// ARQNetSrv 网络服务器
type ARQNetSrv struct {
	protocal  string
	listener  net.Listener
	ctx       context.Context
	ctxCancel context.CancelFunc
	receiver  IConnReceiver

	maxConns int
}

func newARQNetSrv(protocal string, receiver IConnReceiver, maxConns int) INetSrv {
	srv := &ARQNetSrv{
		protocal: protocal,
		listener: nil,
		receiver: receiver,
		maxConns: maxConns,
	}
	srv.ctx, srv.ctxCancel = context.WithCancel(context.Background())
	return srv
}

func (srv *ARQNetSrv) start(addr string) error {
	var err error

	switch srv.protocal {
	case "tcp":
		srv.listener, err = net.Listen("tcp", addr)
	default:
		panic("WRONG PROTOCAL")
	}
	if err != nil {
		return err
	}

	if srv.maxConns > 0 {
		srv.listener = netutil.LimitListener(srv.listener, srv.maxConns)
	}

	go srv.acceptConn()

	return nil
}

func (srv *ARQNetSrv) close() {
	srv.ctxCancel()
	srv.listener.Close()
}

func (srv *ARQNetSrv) acceptConn() {
	for {
		select {
		case <-srv.ctx.Done():
			return
		default:
			{
				conn, err := srv.listener.Accept()
				if err != nil {
					log.Error("accept connection error ", err)
					continue
				}
				go srv.receiver.acceptConn(conn)
			}
		}
	}
}
