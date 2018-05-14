package sess

const (
	msgIDSize       = 2
	minCompressSize = 100

	// MsgHeadSize consist message length , compression type and message id
	MsgHeadSize = 6
)

const (
	maxTCPPacket = 100 * 1024

	// MaxMsgBuffer 消息最大长度
	MaxMsgBuffer = 100 * 1024
)
