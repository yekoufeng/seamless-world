package msgdef

import (
	"fmt"
)

// ProtoSync 同步Proto消息号和名字
type ProtoSync struct {
	Data []byte
}

func (msg *ProtoSync) String() string {
	return fmt.Sprintf("%x\n", msg.Data)
}

// Name 消息名
func (msg *ProtoSync) Name() string {
	return "ProtoSync"
}

// MarshalTo 序列化
func (msg *ProtoSync) MarshalTo(data []byte) (n int, err error) {
	copy(data[0:], msg.Data)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *ProtoSync) Unmarshal(data []byte) error {
	msg.Data = data
	return nil
}

// Size 消息长度
func (msg *ProtoSync) Size() (n int) {
	return len(msg.Data)
}
