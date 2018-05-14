package sess

import (
	"errors"
	"io"
	"net"
	"zeus/msgdef"
)

//////////////////////////////////////////////////////////////////////

func readARQMsgForward(conn net.Conn) (msgdef.IMsg, int, []byte, error) {

	if conn == nil {
		return nil, 0, nil, errors.New("无效连接")
	}

	msgHead := make([]byte, MsgHeadSize-msgIDSize)

	if _, err := io.ReadFull(conn, msgHead); err != nil {
		return nil, 0, nil, err
	}

	msgSize := (int(msgHead[0]) | int(msgHead[1])<<8 | int(msgHead[2])<<16)
	compressFlag := msgHead[3]
	if msgSize > maxTCPPacket || msgSize <= 0 {
		return nil, 0, nil, errors.New("收到的数据长度超过最大值")
	}

	msgData := make([]byte, msgSize)

	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, 0, nil, err
	}

	msgID := GetMsgID(msgData)
	if msgdef.GetMsgDef().IsMsgExist(msgID) {
		msg, err := DecodeMsg(compressFlag, msgData)
		return msg, 0, nil, err
	}

	msgBuf := make([]byte, MsgHeadSize-msgIDSize+msgSize)
	copy(msgBuf, msgHead)
	copy(msgBuf[MsgHeadSize-msgIDSize:], msgData)
	return nil, int(msgID), msgBuf, nil

}
