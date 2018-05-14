package msgdef

import (
	"fmt"
	"zeus/common"
	"zeus/linmath"
)

// EnterAOI 玩家进入AOI
type EnterAOI struct {
	EntityID   uint64
	EntityType string
	State      []byte
	PropNum    uint16
	Properties []byte

	Pos  linmath.Vector3
	Rota linmath.Vector3
}

func (msg *EnterAOI) String() string {
	return fmt.Sprintf("{EntityID:%d EntityType:%s PropNum:%d}", msg.EntityID, msg.EntityType, msg.PropNum)
}

// Name 获取消息名称
func (msg *EnterAOI) Name() string {
	return "EnterAOI"
}

// MarshalTo 序列化
func (msg *EnterAOI) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	bw.Marshal(msg)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *EnterAOI) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	br.Unmarshal(msg)
	return nil
}

// Size 获取长度
// 8 + 2 + len(string) + 4*3 + 4*3 + 2 + 2 + len([]byte)
func (msg *EnterAOI) Size() (n int) {
	return common.CalcSize(msg)
}

////////////////////////////////////////////////////////////////////

// LeaveAOI 离开AOI
type LeaveAOI struct {
	EntityID uint64
}

func (msg *LeaveAOI) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *LeaveAOI) Name() string {
	return "LeaveAOI"
}

// MarshalTo 序列化
func (msg *LeaveAOI) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	bw.Marshal(msg)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *LeaveAOI) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	br.Unmarshal(msg)
	return nil
}

// Size 获取长度
func (msg *LeaveAOI) Size() (n int) {
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////////////////

// EnterCell 玩家进入场景
type EnterCell struct {
	CellID    uint64
	MapName   string
	EntityID  uint64
	Addr      string
	TimeStamp uint32
}

func (msg *EnterCell) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *EnterCell) Name() string {
	return "EnterCell"
}

// MarshalTo 序列化
func (msg *EnterCell) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *EnterCell) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *EnterCell) Size() (n int) {
	return common.CalcSize(msg)
}

////////////////////////////////////////////////////////////////////

// LeaveCell 离开空间
type LeaveCell struct {
}

func (msg *LeaveCell) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *LeaveCell) Name() string {
	return "LeaveCell"
}

// MarshalTo 序列化
func (msg *LeaveCell) MarshalTo(data []byte) (n int, err error) {
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *LeaveCell) Unmarshal(data []byte) error {
	return nil
}

// Size 获取长度
func (msg *LeaveCell) Size() (n int) {
	return 1
}

////////////////////////////////////////////////////////////////////

// EntityAOIS 进入AOI范围
type EntityAOIS struct {
	Num  uint32
	data [][]byte
}

func (msg *EntityAOIS) GetData() [][]byte {
	return msg.data
}

// NewEntityAOISMsg 新建进入AOI消息
func NewEntityAOISMsg() *EntityAOIS {
	return &EntityAOIS{
		0,
		make([][]byte, 0, 1),
	}
}

func (msg *EntityAOIS) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *EntityAOIS) Name() string {
	return "EntityAOIS"
}

// MarshalTo 序列化
func (msg *EntityAOIS) MarshalTo(data []byte) (n int, err error) {
	size := msg.Size()
	bw := common.NewByteStream(data)
	bw.WriteUInt32(msg.Num)

	for i := 0; i < int(msg.Num); i++ {
		bw.WriteBytes(msg.data[i])
	}

	return size, nil
}

// Unmarshal 反序列化
func (msg *EntityAOIS) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	var err error
	msg.Num, err = br.ReadUInt32()
	if err != nil {
		return err
	}
	var oneData []byte
	for i := uint32(0); i < msg.Num; i++ {
		oneData, err = br.ReadBytes()
		if err != nil {
			return err
		}
		msg.data = append(msg.data, oneData)
	}
	return nil
}

// Size 获取长度
func (msg *EntityAOIS) Size() (n int) {
	sum := 0
	for i := 0; i < len(msg.data); i++ {
		sum += len(msg.data[i])
	}

	return 4 + sum + 2*len(msg.data)
}

// AddData 增加一个玩家数据
func (msg *EntityAOIS) AddData(data []byte) {
	msg.Num++
	msg.data = append(msg.data, data)

	// seelog.Debug("EntityAOIS AddData ", msg)
}

// Clear 清理
func (msg *EntityAOIS) Clear() {
	msg.Num = 0
	msg.data = msg.data[0:0]
}

// ////////////////////////////////////////////////////////////////
