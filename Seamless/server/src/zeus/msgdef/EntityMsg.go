package msgdef

import (
	"encoding/binary"
	"fmt"
	"zeus/common"
)

///////////////////////////////////////////////////////////////

// ClientTransport 客户端中转消息
// 当客户端消息要投递给非网关对象时，用该消息包装
type ClientTransport struct {
	SrvType    uint8
	MsgFlag    uint8
	MsgContent []byte
}

func (msg *ClientTransport) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *ClientTransport) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *ClientTransport) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// []byte长度额外2个字节, 两个uint8 2个字节
// 固定长度4字节
func (msg *ClientTransport) Size() (n int) {
	//return 4 + len(msg.MsgContent)
	return common.CalcSize(msg)
}

// Name 获取名字
func (msg *ClientTransport) Name() string {
	return "ClientTransport"
}

// CreateEntityReq 请求创建实体消息
type CreateEntityReq struct {
	EntityType string
	EntityID   uint64
	CellID     uint64
	DBID       uint64
	InitParam  []byte
	SrcSrvType uint8
	SrcSrvID   uint64
	CallbackID uint32
}

func (msg *CreateEntityReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *CreateEntityReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *CreateEntityReq) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *CreateEntityReq) Size() (n int) {
	//return 39 + len(msg.EntityType)
	return common.CalcSize(msg)
}

// Name 获取名字
func (msg *CreateEntityReq) Name() string {
	return "CreateEntityReq"
}

////////////////////////////////////////////////////////////

// CreateEntityRet 创建实体返回
type CreateEntityRet struct {
	SrvType    uint8
	EntityID   uint64
	SpaceID    uint64
	CallbackID uint32
	ErrorStr   string
}

func (msg *CreateEntityRet) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *CreateEntityRet) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *CreateEntityRet) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 1 + 8 + 8 + 4 + 2 + len(string)
func (msg *CreateEntityRet) Size() (n int) {
	//return 23 + len(msg.ErrorStr)
	return common.CalcSize(msg)
}

// Name 获取名字
func (msg *CreateEntityRet) Name() string {
	return "CreateEntityRet"
}

///////////////////////////////////////////////////////////////

// DestroyEntityReq 请求删除实体消息
type DestroyEntityReq struct {
	EntityID   uint64
	SrcSrvType uint8
	SrcSrvID   uint64
	CallbackID uint32
	CellID     uint64
}

func (msg *DestroyEntityReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *DestroyEntityReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *DestroyEntityReq) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 8 + 1 + 8 + 4 + 8
func (msg *DestroyEntityReq) Size() (n int) {
	//return 29
	return common.CalcSize(msg)
}

// Name 获取名字
func (msg *DestroyEntityReq) Name() string {
	return "DestroyEntityReq"
}

////////////////////////////////////////////////////////////

// DestroyEntityRet 销毁实体返回
type DestroyEntityRet struct {
	SrvType    uint8
	SrvID      uint64
	EntityID   uint64
	SpaceID    uint64
	CallbackID uint32
	ErrorStr   string
}

func (msg *DestroyEntityRet) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *DestroyEntityRet) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *DestroyEntityRet) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 1 + 8 + 8 + 8 + 4 + 2 + len(string)
func (msg *DestroyEntityRet) Size() (n int) {
	//return 31 + len(msg.ErrorStr)
	return common.CalcSize(msg)
}

// Name 获取名字
func (msg *DestroyEntityRet) Name() string {
	return "DestroyEntityRet"
}

///////////////////////////////////////////////////////////

// EntityMsgTransport 分布式实体之间传递消息用
type EntityMsgTransport struct {
	SrvType    uint8
	EntityID   uint64
	CellID     uint64
	IsGateway  bool //是不是gateway中转过来的数据
	MsgContent []byte
}

func (msg *EntityMsgTransport) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *EntityMsgTransport) MarshalTo(data []byte) (n int, err error) {
	data[0] = msg.SrvType
	binary.LittleEndian.PutUint64(data[1:9], msg.EntityID)
	binary.LittleEndian.PutUint64(data[9:17], msg.CellID)
	copy(data[17:], msg.MsgContent)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *EntityMsgTransport) Unmarshal(data []byte) error {
	msg.SrvType = data[0]
	msg.EntityID = binary.LittleEndian.Uint64(data[1:9])
	msg.CellID = binary.LittleEndian.Uint64(data[9:17])
	msg.MsgContent = data[17:]
	return nil
}

// Size 获取长度
func (msg *EntityMsgTransport) Size() (n int) {
	return 17 + len(msg.MsgContent)
}

// Name 获取名字
func (msg *EntityMsgTransport) Name() string {
	return "EntityMsgTransport"
}

// EntityVarData 分布式实体变量的消息
type EntityVarData struct {
	Identifier string
	Variant    []byte
}

func (msg *EntityVarData) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// MarshalTo 序列化
func (msg *EntityVarData) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *EntityVarData) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 2 + len(string) + 2 + len([]byte)
func (msg *EntityVarData) Size() (n int) {
	//return 4 + len(msg.Identifier) + len(msg.Variant)
	return common.CalcSize(msg)
}

// Name 获取名字
func (msg *EntityVarData) Name() string {
	return "EntityVarData"
}

///////////////////////////////////////////////////////

// EntityMsgChange 分布式实体之间同步数据使用
type EntityMsgChange struct {
	VarData []EntityVarData
}

func (msg *EntityMsgChange) String() string {
	var str string
	for _, v := range msg.VarData {
		str += v.String()
		str += "\n"
	}
	return str
}

// MarshalTo 序列化
func (msg *EntityMsgChange) MarshalTo(data []byte) (n int, err error) {
	count := uint16(len(msg.VarData))
	binary.LittleEndian.PutUint16(data[0:2], count)
	pos := 2
	for _, v := range msg.VarData {
		i, _ := v.MarshalTo(data[pos:])
		pos += i
	}
	return msg.Size(), nil
}

// Unmarshal 反序列化
//这里代码如果data数据有问题，会导致v.Unmarshal(data[pos:]) crash，先记录下
func (msg *EntityMsgChange) Unmarshal(data []byte) error {
	count := binary.LittleEndian.Uint16(data[0:2])
	pos := 2
	for i := uint16(0); i < count; i++ {
		v := new(EntityVarData)
		v.Unmarshal(data[pos:])
		msg.VarData = append(msg.VarData, *v)
		pos += v.Size()
	}

	return nil
}

// Size 获取长度
func (msg *EntityMsgChange) Size() (n int) {
	i := 2
	for _, v := range msg.VarData {
		i += v.Size()
	}
	return i
}

// Name 获取名字
func (msg *EntityMsgChange) Name() string {
	return "EntityMsgChange"
}

/////////////////////////////////////////////////////////////////////////////

// RPCMsg 分布式实体之间调用方法
type RPCMsg struct {
	ServerType  uint8
	SrcEntityID uint64
	MethodName  string
	Data        []byte //proto消息参数
}

func (msg *RPCMsg) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 消息名
func (msg *RPCMsg) Name() string {
	return "RPCMsg"
}

// MarshalTo 序列化
func (msg *RPCMsg) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *RPCMsg) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 消息长度
// 1 + 8 + 2 + len(string) + 2 + len([]byte)
func (msg *RPCMsg) Size() (n int) {
	//return 13 + len(msg.Data) + len(msg.MethodName)
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////

// CellEntityMsg 客户端直接投递到spaceEntity时的包装消息
type CellEntityMsg struct {
	UID  uint64
	Data []byte //具体的消息
}

func (msg *CellEntityMsg) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 消息名
func (msg *CellEntityMsg) Name() string {
	return "CellEntityMsg"
}

// MarshalTo 序列化
func (msg *CellEntityMsg) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *CellEntityMsg) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 消息长度
func (msg *CellEntityMsg) Size() (n int) {
	return common.CalcSize(msg)
}

///////////////////////////////////////////////////////////

// EntitySrvInfoNotify 客户端直接投递到spaceEntity时的包装消息
type EntitySrvInfoNotify struct {
}

func (msg *EntitySrvInfoNotify) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 消息名
func (msg *EntitySrvInfoNotify) Name() string {
	return "EntitySrvInfoNotify"
}

// MarshalTo 序列化
func (msg *EntitySrvInfoNotify) MarshalTo(data []byte) (n int, err error) {
	return 1, nil
}

// Unmarshal 反序列化
func (msg *EntitySrvInfoNotify) Unmarshal(data []byte) error {
	return nil
}

// Size 消息长度
func (msg *EntitySrvInfoNotify) Size() (n int) {
	return 1
}

/////////////////////////////////////////////////////////////////////////////

// EntityEvent 分布式实体之间调用方法
type EntityEvent struct {
	SrcEntityID uint64
	EventName   string
	Data        []byte //参数
}

func (msg *EntityEvent) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 消息名
func (msg *EntityEvent) Name() string {
	return "EntityEvent"
}

// MarshalTo 序列化
func (msg *EntityEvent) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *EntityEvent) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 消息长度
func (msg *EntityEvent) Size() (n int) {
	return common.CalcSize(msg)
}
