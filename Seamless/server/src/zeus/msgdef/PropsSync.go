package msgdef

import (
	"fmt"
	"zeus/common"
)

// PropsSync 玩家属性变化，服务器之间交换属性信息
type PropsSync struct {
	Num  uint32
	Data []byte
}

func (msg *PropsSync) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *PropsSync) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *PropsSync) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 4 + 2 + len([]byte)
func (msg *PropsSync) Size() (n int) {
	return 6 + len(msg.Data)
}

// Name 获取名字
func (msg *PropsSync) Name() string {
	return "PropsSync"
}

/////////////////////////////////////////////////////////////////

// PropsSyncClient 玩家属性变化，客户端关注的属性变化都由该消息发出
type PropsSyncClient struct {
	EntityID uint64
	Num      uint32
	Data     []byte
}

func (msg *PropsSyncClient) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *PropsSyncClient) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *PropsSyncClient) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 8 + 4 + 2 + len([]byte)
func (msg *PropsSyncClient) Size() (n int) {
	return 14 + len(msg.Data)
}

// Name 获取名字
func (msg *PropsSyncClient) Name() string {
	return "PropsSyncClient"
}

/////////////////////////////////////////////////////////////////

// MRolePropsSyncClient 主角属性变化，客户端主角关注的属性变化由该消息发出
type MRolePropsSyncClient struct {
	EntityID uint64
	Num      uint32
	Data     []byte
}

func (msg *MRolePropsSyncClient) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *MRolePropsSyncClient) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *MRolePropsSyncClient) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 8 + 4 + 2 + len([]byte)
func (msg *MRolePropsSyncClient) Size() (n int) {
	return 14 + len(msg.Data)
}

// Name 获取名字
func (msg *MRolePropsSyncClient) Name() string {
	return "MRolePropsSyncClient"
}
