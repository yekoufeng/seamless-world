package msgdef

import (
	"encoding/binary"
	"fmt"
)

// TestBinMsg 会话状态通知, 交给业务层处理
type TestBinMsg struct {
	State uint32
}

func (msg *TestBinMsg) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取名字
func (msg *TestBinMsg) Name() string {
	return "TestBinMsg"
}

// MarshalTo 序列化
func (msg *TestBinMsg) MarshalTo(data []byte) (n int, err error) {
	binary.LittleEndian.PutUint32(data, msg.State)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *TestBinMsg) Unmarshal(data []byte) error {
	msg.State = binary.LittleEndian.Uint32(data)
	return nil
}

// Size 获取长度
func (msg *TestBinMsg) Size() (n int) {
	return 4
}
