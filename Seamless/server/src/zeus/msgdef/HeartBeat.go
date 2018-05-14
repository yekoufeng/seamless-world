package msgdef

// HeartBeat 心跳消息
type HeartBeat struct {
}

func (hb *HeartBeat) Name() string {
	return "HeartBeat"
}

// MarshalTo 序列化
func (hb *HeartBeat) MarshalTo(data []byte) (n int, err error) {
	return hb.Size(), nil
}

// Unmarshal 反序列化
func (hb *HeartBeat) Unmarshal(data []byte) error {
	return nil
}

// Size 获取长度
func (hb *HeartBeat) Size() (n int) {
	return 0
}
