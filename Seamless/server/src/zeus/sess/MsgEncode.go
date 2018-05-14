package sess

import (
	"encoding/binary"
	"errors"
	"reflect"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
	"github.com/golang/snappy"
)

/*
消息定义文件
消息格式


流式结构如下：
消息头| 消息 | 消息头 | 消息

消息头，前三位表示长度，最后一位表示压缩方式

单条消息组成结构

msg_id 消息号 2 byte
msg_body 消息体


消息ID号


服务器有一份消息id列表，客户端每次上线都同步一下
消息id列表自动生成，大致结构如下 ：
消息号，消息名称
1 号消息就是验证消息
Send()

*/

// GetMsgID 获取消息ID
func GetMsgID(buf []byte) uint16 {
	if len(buf) < msgIDSize {
		return 0
	}

	return uint16(buf[0]) | uint16(buf[1])<<8
}

var encrypt_key = []byte{41, 253, 1, 56, 52, 62, 176, 42}
var decrypt_key = []byte{41, 253, 1, 56, 52, 62, 176, 42}

// EncryptData 加密算法
func EncryptData(buf []byte) []byte {
	buflen := len(buf)
	key := encrypt_key
	keylen := len(key)

	for i := 0; i < buflen; i++ {
		n := byte(i%7 + 1)                       //移位长度(1-7)
		b := (buf[i] << n) | (buf[i] >> (8 - n)) // 向左循环移位

		buf[i] = b ^ key[i%keylen]
	}

	return buf
}

// DecryptData 解密算法
func DecryptData(buf []byte) []byte {

	buflen := len(buf)
	key := decrypt_key
	keylen := len(key)

	for i := 0; i < buflen; i++ {

		b := buf[i] ^ key[i%keylen]

		n := byte(i%7 + 1)                 //移位长度(1-7)
		buf[i] = (b >> n) | (b << (8 - n)) // 向右循环移位
	}
	return buf
}

// DecodeMsg 返回消息名称及反序列化后的消息对象
func DecodeMsg(flag byte, buf []byte) (msgdef.IMsg, error) {

	if len(buf) < msgIDSize {
		log.Error("数据格式错误, buf:", len(buf))
		return nil, errors.New("长度错误")
	}
	_, msgContent, err := msgdef.GetMsgDef().GetMsgInfo(GetMsgID(buf))
	if err != nil {
		return nil, err
	}

	msgBody := buf[msgIDSize:]

	encryptFlag := flag & 0x2
	if encryptFlag > 0 {
		msgBody = DecryptData(msgBody)
	}

	compressFlag := flag & 0x1

	if compressFlag == 0 {

		if err = msgContent.Unmarshal(msgBody); err != nil {
			return nil, err
		}

	} else if compressFlag == 1 {

		msgBuf := make([]byte, MaxMsgBuffer)
		unCompressBuf, err := snappy.DecodeWithBuf(msgBuf, msgBody)
		if err != nil {
			return nil, err
		}

		if err = msgContent.Unmarshal(unCompressBuf); err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("error compression flag")
	}

	return msgContent, nil
}

func EncodeMsg(msg msgdef.IMsg, buf []byte, forceNoCompress bool) ([]byte, error) {
	return EncodeMsgWithEncrypt(msg, buf, forceNoCompress, false)
}

// EncodeMsg 序列化消息
func EncodeMsgWithEncrypt(msg msgdef.IMsg, buf []byte, forceNoCompress bool, encryptEnabled bool) ([]byte, error) {

	if msg == nil {
		return nil, errors.New("消息错误，消息不能为nil")
	}

	msgID, err := msgdef.GetMsgDef().GetMsgIDByName(reflect.TypeOf(msg).Elem().Name())
	if err != nil {
		return nil, err
	}

	var msgbuf []byte

	size := msg.Size()
	if size >= minCompressSize && forceNoCompress == false {
		msgbuf, err = _snappyCompressCmd(msgID, msg, size, buf)
	}
	msgbuf, err = _noCompressCmd(msgID, msg, size, buf)
	if err != nil {
		return nil, err
	}

	if encryptEnabled {
		data := msgbuf[MsgHeadSize:]
		EncryptData(data)
		msgbuf[3] = msgbuf[3] | 0x2
	}

	return msgbuf, err
}

func _snappyCompressCmd(cmd uint16, msg msgdef.IMsg, msgSize int, buf []byte) ([]byte, error) {

	msgdata := make([]byte, msgSize)

	n, err := msg.MarshalTo(msgdata)
	if err != nil {
		log.Error("[协议] 编码错误 ", err)
		return nil, err
	}
	data := msgdata[:n]
	maxLen := snappy.MaxEncodedLen(len(data))

	if maxLen+MsgHeadSize > len(buf) {
		return nil, errors.New("message size too large")
	}

	p := buf[:maxLen+MsgHeadSize]

	mbuff := snappy.Encode(p[MsgHeadSize:], data)
	cmdsize := len(mbuff) + msgIDSize
	p[0] = byte(cmdsize)
	p[1] = byte(cmdsize >> 8)
	p[2] = byte(cmdsize >> 16)
	p[3] = 1
	binary.LittleEndian.PutUint16(p[4:], cmd)
	return p[:len(mbuff)+MsgHeadSize], nil
}

func _noCompressCmd(cmd uint16, msg msgdef.IMsg, msgSize int, buf []byte) ([]byte, error) {

	if msgSize+MsgHeadSize > len(buf) {
		return nil, errors.New("message size to large ")
	}

	data := buf
	n, err := msg.MarshalTo(data[MsgHeadSize:])
	if err != nil {
		log.Error("[协议] 编码错误 ", err)
		return nil, err
	}
	cmdsize := n + msgIDSize
	data[0] = byte(cmdsize)
	data[1] = byte(cmdsize >> 8)
	data[2] = byte(cmdsize >> 16)
	data[3] = 0
	binary.LittleEndian.PutUint16(data[4:], cmd)
	return data[:n+MsgHeadSize], nil
}
