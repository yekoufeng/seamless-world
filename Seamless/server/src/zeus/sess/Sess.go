package sess

import (
	"context"
	"net"
	"time"
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/msghandler"
	"zeus/pool"
	"zeus/safecontainer"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

func newSess(conn net.Conn, notifier IConnNotify, encryEnabled bool) iserver.ISess {
	sess := &NetSess{
		IMsgHandlers:  msghandler.NewMsgHandlers(),
		connNotifier:  notifier,
		conn:          conn,
		_isVertified:  false,
		sendBuf:       safecontainer.NewSafeList(),
		sendRawBuf:    safecontainer.NewSafeList(),
		backlogMsg:    make([]msgdef.IMsg, 0, 1),
		backlogRawMsg: make([][]byte, 0, 1),

		id:               0,
		hbTimerInterval:  time.Duration(viper.GetInt("Config.HBTimer")) * time.Second,
		hbTickerInterval: time.Duration(viper.GetInt("Config.HBTicker")) * time.Second,

		_isClosed:     false,
		_isError:      false,
		_encryEnabled: encryEnabled,
	}

	sess.hbTimer = time.NewTimer(sess.hbTimerInterval)
	sess.hbTicker = time.NewTicker(sess.hbTickerInterval)

	sess.ctx, sess.ctxCancel = context.WithCancel(context.Background())
	return sess
}

// NetSess 代表一个网络连接
type NetSess struct {
	msghandler.IMsgHandlers
	connNotifier IConnNotify
	conn         net.Conn
	_isVertified bool
	sendBuf      *safecontainer.SafeList
	sendRawBuf   *safecontainer.SafeList
	lastMsg      msgdef.IMsg
	lastRawMsg   []byte

	backlogMsg        []msgdef.IMsg
	backlogRawMsg     [][]byte
	backlogLastMsg    msgdef.IMsg
	backlogLastRawMsg []byte

	ctx       context.Context
	ctxCancel context.CancelFunc
	id            uint64

	//服务器类型
	serverType    uint8

	hbTimer          *time.Timer
	hbTimerInterval  time.Duration
	hbTicker         *time.Ticker
	hbTickerInterval time.Duration

	_isClosed     bool
	_isError      bool
	_encryEnabled bool
}

// SetMsgHandler 设置消息处理器
func (sess *NetSess) SetMsgHandler(handler msghandler.IMsgHandlers) {
	if handler == sess || handler == nil {
		sess.IMsgHandlers = msghandler.NewMsgHandlers()
	} else {
		sess.IMsgHandlers = handler
	}
}

// Start 验证完成
func (sess *NetSess) Start() {
	go sess.recvLoop()
	go sess.sendLoop()
	if viper.GetBool("Config.HeartBeat") {
		go sess.hbLoop()
	}
}

// Send 发送消息
func (sess *NetSess) Send(msg msgdef.IMsg) {
	if sess._isClosed {
		log.Warnf("Send after sess close %s %s %s", sess.conn.RemoteAddr(), msg.Name(), msg)
		return
	}

	sess.sendBuf.Put(msg)
}

// SendRaw 发送原始消息
func (sess *NetSess) SendRaw(rawMsg []byte) {
	if sess._isClosed {
		log.Warn("SendRaw after sess close ", sess.conn.RemoteAddr())
		return
	}

	sess.sendRawBuf.Put(rawMsg)
}

// Touch 记录心跳状态
func (sess *NetSess) Touch() {
	sess.hbTimer.Reset(sess.hbTimerInterval)
}

func (sess *NetSess) hbLoop() {
	for {
		select {
		case <-sess.ctx.Done():
			sess.hbTimer.Stop()
			return
		case <-sess.hbTimer.C:
			log.Error("sess heart tick expired ", sess.conn.RemoteAddr())

			sess._isError = true
			sess.hbTimer.Stop()
			sess.Close()
			return
		case <-sess.hbTicker.C:
			sess.Send(&msgdef.HeartBeat{})
		}
	}
}

func (sess *NetSess) recvLoop() {

	for {

		select {
		case <-sess.ctx.Done():
			return
		default:
			msg, msgID, rawMsg, err := readARQMsgForward(sess.conn)
			if err != nil {
				log.Error("tcp read message error ", err, sess.conn.RemoteAddr())
				if sess.IsClosed() {
					return
				}

				sess._isError = true
				if sess.connNotifier != nil {
					sess.connNotifier.connErr(sess.GetID())
				} else {
					sess.Close()
				}
				return
			}

			if msg != nil && msg.Name() != "HeartBeat" {
				sess.FireMsg(msg.Name(), msg)
			}

			if rawMsg != nil {
				sess.FireMsg("MsgForward", &RawMsg{msgID, rawMsg})
			}

			sess.Touch()
		}
	}
}

func (sess *NetSess) sendLoop() {

	for {
		select {
		case <-sess.ctx.Done():
			return
		case <-sess.sendBuf.C:
			for {
				if sess.IsClosed() {
					return
				}

				data, err := sess.sendBuf.Pop()
				if err != nil {
					break
				}

				m := data.(msgdef.IMsg)
				buf := pool.Get(MaxMsgBuffer)
				msgBuf, err := EncodeMsgWithEncrypt(m, buf, true, sess._encryEnabled)
				if err != nil {
					log.Error("encode message error ", err)
					continue
				}

				_, err = sess.conn.Write(msgBuf)
				pool.Put(buf)
				if err != nil {
					log.Error("send message error ", err)

					sess._isError = true
					sess.lastMsg = m
					if sess.connNotifier != nil {
						sess.connNotifier.connErr(sess.GetID())
					} else {
						sess.Close()
						return
					}
				}
			}
		case <-sess.sendRawBuf.C:
			for {
				if sess.IsClosed() {
					return
				}

				data, err := sess.sendRawBuf.Pop()
				if err != nil {
					break
				}

				_, err = sess.conn.Write(data.([]byte))
				if err != nil {
					log.Error("send message error ", err)

					sess._isError = true
					sess.lastRawMsg = data.([]byte)
					if sess.connNotifier != nil {
						sess.connNotifier.connErr(sess.GetID())
					} else {
						sess.Close()
						return
					}
				}
			}
		}
	}
}

// Close 关闭
func (sess *NetSess) Close() {
	if sess._isClosed {
		return
	}

	sess._isClosed = true

	sess.hbTicker.Stop()
	sess.hbTimer.Stop()

	if sess.connNotifier != nil {
		sess.connNotifier.connErr(sess.GetID())
	}

	if sess.IMsgHandlers != nil {
		sess.DoNormalMsg("SessClosed", sess.GetID())
	}

	if !sess._isError {
		go func() {
			closeTicker := time.NewTicker(100 * time.Millisecond)
			defer closeTicker.Stop()
			for {
				select {
				case <-closeTicker.C:
					if sess.sendBuf.IsEmpty() && sess.sendRawBuf.IsEmpty() {
						sess.ctxCancel()
						sess.conn.Close()
						return
					}
				}
			}
		}()
	} else {
		sess.ctxCancel()
		sess.conn.Close()
	}
}

// RemoteAddr 远程地址
func (sess *NetSess) RemoteAddr() string {
	return sess.conn.RemoteAddr().String()
}

// IsVertified 是否验证的连接
func (sess *NetSess) IsVertified() bool {
	return sess._isVertified
}

// SetVertify 设置已经验证
func (sess *NetSess) SetVertify() {
	sess._isVertified = true
}

// SetID 设置ID
func (sess *NetSess) SetID(id uint64) {
	sess.id = id
}

// GetID 获取ID
func (sess *NetSess) GetID() uint64 {
	return sess.id
}

// SetID 设置servertype
func (sess *NetSess) SetServerType(servertype uint8) {
	sess.serverType = servertype
}

// GetID 获取servertype
func (sess *NetSess) GetServerType() uint8 {
	return sess.serverType
}


// IsClosed 返回sess是否已经关闭
func (sess *NetSess) IsClosed() bool {
	return sess._isClosed
}

// FetchBacklog 获取积压消息
func (sess *NetSess) FetchBacklog(o iserver.ISess) {
	sess.backlogMsg = sess.backlogMsg[0:0]
	sess.backlogRawMsg = sess.backlogRawMsg[0:0]

	sess.backlogLastMsg = sess.lastMsg
	sess.backlogLastRawMsg = sess.lastRawMsg

	os := o.(*NetSess)
	for {
		m, err := os.sendBuf.Pop()
		if err != nil {
			break
		}
		sess.backlogMsg = append(sess.backlogMsg, m.(msgdef.IMsg))
	}

	for {
		m, err := os.sendRawBuf.Pop()
		if err != nil {
			break
		}
		sess.backlogRawMsg = append(sess.backlogRawMsg, m.([]byte))
	}

	var fetchStr string
	if sess.backlogLastMsg != nil {
		fetchStr = sess.backlogLastMsg.Name()
	}

	for _, msg := range sess.backlogMsg {
		fetchStr += " " + msg.Name()
	}

	log.Debug("FetchBacklog:", fetchStr)
}

// FlushBacklog 刷新积压消息
func (sess *NetSess) FlushBacklog() {
	var flushStr string

	if sess.backlogLastMsg != nil {
		flushStr += " " + sess.backlogLastMsg.Name()
		sess.Send(sess.backlogLastMsg)
		sess.backlogLastMsg = nil
	}

	for _, m := range sess.backlogMsg {
		flushStr += " " + m.Name()
		sess.Send(m)
	}

	sess.backlogMsg = sess.backlogMsg[0:0]

	if sess.backlogLastRawMsg != nil {
		sess.SendRaw(sess.backlogLastRawMsg)
		sess.backlogLastRawMsg = nil
	}

	for _, rm := range sess.backlogRawMsg {
		sess.SendRaw(rm)
	}

	sess.backlogRawMsg = sess.backlogRawMsg[0:0]

	log.Debug("FlushBacklog:", flushStr)
}
