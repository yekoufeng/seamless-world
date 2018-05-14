package msgdef

import (
	"bytes"
	"encoding/binary"
)

/*
	些文件定义帧同步相关消息
*/

// ClientFrameMsg 客户端上发的帧消息
type ClientFrameMsg struct {
	Data []byte
}

// MarshalTo 序列化
func (msg *ClientFrameMsg) MarshalTo(data []byte) (n int, err error) {
	copy(data[0:], msg.Data)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *ClientFrameMsg) Unmarshal(data []byte) error {
	msg.Data = data
	return nil
}

// Size 获取长度
func (msg *ClientFrameMsg) Size() int {
	return 5
}

func (msg *ClientFrameMsg) Name() string {
	return "ClientFrameMsg"
}

// ClientFrameMsgData 拼接UID之后的单用户帧消息
type ClientFrameMsgData struct {
	// 最高位 0标识正常玩家, 1标识AI托管
	UID  uint8
	Data []byte
}

// MarshalTo 序列化
func (msg *ClientFrameMsgData) MarshalTo(data []byte) (n int, err error) {
	data[0] = byte(msg.UID)
	copy(data[1:6], msg.Data)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *ClientFrameMsgData) Unmarshal(data []byte) error {
	msg.UID = uint8(data[0])
	msg.Data = data[1:6]
	return nil
}

// Size 获取长度
func (msg *ClientFrameMsgData) Size() int {
	return 6
}

func (msg *ClientFrameMsgData) Name() string {
	return "ClientFrameMsgData"
}

// ServerFrameMsg 服务器下发的帧消息
type ServerFrameMsg struct {
	FrameID uint16
	Msgs    []ClientFrameMsgData
}

// MarshalTo 序列化
func (msg *ServerFrameMsg) MarshalTo(data []byte) (n int, err error) {
	binary.LittleEndian.PutUint16(data[0:2], msg.FrameID)
	data[2] = uint8(len(msg.Msgs))
	pos := 3
	for _, v := range msg.Msgs {
		i, _ := v.MarshalTo(data[pos:])
		pos += i
	}
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *ServerFrameMsg) Unmarshal(data []byte) error {
	msg.FrameID = binary.LittleEndian.Uint16(data[0:2])
	pos := 3
	t := &ClientFrameMsgData{}
	len := t.Size()
	for i := 0; i < int(data[2]); i++ {
		v := ClientFrameMsgData{}
		v.Unmarshal(data[pos : pos+len])
		msg.Msgs = append(msg.Msgs, v)
		pos += len
	}
	return nil
}

// Size 获取长度
func (msg *ServerFrameMsg) Size() int {
	i := &ClientFrameMsgData{}
	return 3 + len(msg.Msgs)*i.Size()
}

func (msg *ServerFrameMsg) Name() string {
	return "ServerFrameMsg"
}

// FramesMsg 服务器下发的多帧消息
type FramesMsg struct {
	Frames []ServerFrameMsg
}

// MarshalTo 序列化
func (msg *FramesMsg) MarshalTo(data []byte) (n int, err error) {
	count := uint16(len(msg.Frames))
	binary.LittleEndian.PutUint16(data[0:2], count)
	pos := 2
	for _, v := range msg.Frames {
		i, _ := v.MarshalTo(data[pos:])
		pos += i
	}
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *FramesMsg) Unmarshal(data []byte) error {
	count := binary.LittleEndian.Uint16(data[0:2])
	pos := 2
	for i := uint16(0); i < count; i++ {
		f := new(ServerFrameMsg)
		f.Unmarshal(data[pos:])
		msg.Frames = append(msg.Frames, *f)
		pos += f.Size()
	}

	return nil
}

// Size 获取长度
func (msg *FramesMsg) Size() int {
	i := 2
	for _, v := range msg.Frames {
		i += v.Size()
	}
	return i
}

func (msg *FramesMsg) Name() string {
	return "FramesMsg"
}

// RequireFramesMsg 客户端请求重传帧的消息
type RequireFramesMsg struct {
	FrameID uint16
}

// MarshalTo 序列化
func (msg *RequireFramesMsg) MarshalTo(data []byte) (n int, err error) {
	binary.LittleEndian.PutUint16(data[0:2], msg.FrameID)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *RequireFramesMsg) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, msg)
}

// Size 获取长度
func (msg *RequireFramesMsg) Size() int {
	return 2
}

func (msg *RequireFramesMsg) Name() string {
	return "RequireFramesMsg"
}
