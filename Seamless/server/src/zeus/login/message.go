package login

// UserLoginReq 登录消息格式
type UserLoginReq struct {
	User      string
	Password  string
	Token     string
	Channel   string
	ClientVer string
	Data      []byte
}

// UserLoginAck 登录消息返回格式
type UserLoginAck struct {
	UID       uint64
	Token     string
	LobbyAddr string
	Result    int
	ResultMsg string
	HB        bool
	// Config    interface{}
}

// UserCreateReq 创建帐号消息
type UserCreateReq struct {
	User     string
	Password string
}

// UserCreateAck 创建帐号返回
type UserCreateAck struct {
	UID       uint64
	Result    int
	ResultMsg string
}
