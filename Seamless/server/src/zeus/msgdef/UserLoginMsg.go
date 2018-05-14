package msgdef

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"zeus/common"
)

/*
	msgdef 包下面主要定义二进制消息流
	由于二进制消息流个数固定，所以在代码中直接定义，序列化和反序列化
*/

// ClientVertifyReq 验证消息
type ClientVertifyReq struct {
	// Source: 消息来源, 分客户端或者服务器(ClientMSG/服务器类型)
	Source uint8 //data[0]
	// UID: 玩家UID或者服务器ID
	UID uint64 //data[1:9]
	// Token: 客户端登录时需要携带Token
	Token string //data[9:41]
}

func (msg *ClientVertifyReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 消息名
func (msg *ClientVertifyReq) Name() string {
	return "ClientVertifyReq"
}

// MarshalTo 序列化
func (msg *ClientVertifyReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *ClientVertifyReq) Unmarshal(data []byte) error {

	sw := common.NewByteStream(data)
	msg.Source, _ = sw.ReadByte()
	msg.UID, _ = sw.ReadUInt64()
	msg.Token, _ = sw.ReadStr()

	return nil
}

// Size 获取长度
func (msg *ClientVertifyReq) Size() (n int) {
	return 1 + 8 + len(msg.Token) + 2
}

//////////////////////////////////////////////////////////////////////////////////////////

// ClientVertifySucceedRet 验证成功返回消息
type ClientVertifySucceedRet struct {
	// Source: 消息来源, 分客户端或者服务器(ClientMSG/服务器类型)
	Source uint8
	// UID: 客户端验证时返回UID
	UID uint64
	// SourceID: 验证成功消息的来源
	SourceID uint64
	// Type: 连接类型
	Type uint8
}

func (msg *ClientVertifySucceedRet) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 名字
func (msg *ClientVertifySucceedRet) Name() string {
	return "ClientVertifySucceedRet"
}

// MarshalTo 序列化
func (msg *ClientVertifySucceedRet) MarshalTo(data []byte) (n int, err error) {
	data[0] = byte(msg.Source)
	binary.LittleEndian.PutUint64(data[1:9], msg.UID)
	binary.LittleEndian.PutUint64(data[9:17], msg.SourceID)
	data[17] = byte(msg.Type)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *ClientVertifySucceedRet) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, msg)
}

// Size 获取长度
func (msg *ClientVertifySucceedRet) Size() (n int) {
	return 18 //binary.Size(msg)
}

//////////////////////////////////////////////////////////////////////////////////////////

// ClientVertifyFailedRet 验证成功返回消息
type ClientVertifyFailedRet struct {
}

func (msg *ClientVertifyFailedRet) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 名字
func (msg *ClientVertifyFailedRet) Name() string {
	return "ClientVertifyFailedRet"
}

// MarshalTo 序列化
func (msg *ClientVertifyFailedRet) MarshalTo(data []byte) (n int, err error) {
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *ClientVertifyFailedRet) Unmarshal(data []byte) error {
	return nil
}

// Size 获取长度
func (msg *ClientVertifyFailedRet) Size() (n int) {
	return 0
}

//UserDuplicateLoginNotify 玩家重复登录通知
type UserDuplicateLoginNotify struct {
}

func (msg *UserDuplicateLoginNotify) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 名字
func (msg *UserDuplicateLoginNotify) Name() string {
	return "UserDuplicateLoginNotify"
}

// MarshalTo 序列化
func (msg *UserDuplicateLoginNotify) MarshalTo(data []byte) (n int, err error) {
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *UserDuplicateLoginNotify) Unmarshal(data []byte) error {
	return nil
}

// Size 获取长度
func (msg *UserDuplicateLoginNotify) Size() (n int) {
	return 0
}
