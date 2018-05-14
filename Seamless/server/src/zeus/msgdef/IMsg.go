package msgdef

// IMsg 消息接口，所有的消息类都必须实现的接口
type IMsg interface {
	MarshalTo(data []byte) (n int, err error)
	Unmarshal(data []byte) error
	Size() (n int)
	Name() string
}
