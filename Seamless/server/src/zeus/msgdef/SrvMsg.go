package msgdef

import (
	"fmt"
	"zeus/common"
)

// SrvMsgTransport 分布式实体之间传递消息用
type SrvMsgTransport struct {
	CellID     uint64
	MsgContent []byte
}

func (msg *SrvMsgTransport) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *SrvMsgTransport) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *SrvMsgTransport) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 8 + 2 + len([]byte)
func (msg *SrvMsgTransport) Size() (n int) {
	return 10 + len(msg.MsgContent)
}

// Name 获取名字
func (msg *SrvMsgTransport) Name() string {
	return "SrvMsgTransport"
}
