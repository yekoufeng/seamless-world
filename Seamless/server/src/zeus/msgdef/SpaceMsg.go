package msgdef

import (
	"encoding/binary"
	"fmt"
	"zeus/common"
	"zeus/linmath"
)

// EnterCellReq 请求进入空间
type EnterCellReq struct {
	SrvID      uint64
	CellID     uint64
	EntityType string
	EntityID   uint64
	DBID       uint64
	InitParam  []byte
	OldSrvID   uint64
	OldCellID  uint64
	Pos        linmath.Vector3
}

func (msg *EnterCellReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *EnterCellReq) Name() string {
	return "EnterCellReq"
}

// MarshalTo 序列化
func (msg *EnterCellReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *EnterCellReq) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
// 8 + 8 + 2 + len(string) + 8 + 8 + 8 + 8
func (msg *EnterCellReq) Size() (n int) {
	//return 50 + len(msg.EntityType)
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////////

// LeaveCellReq 请求离开空间
type LeaveCellReq struct {
	EntityID uint64
}

func (msg *LeaveCellReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *LeaveCellReq) Name() string {
	return "LeaveCellReq"
}

// MarshalTo 序列化
func (msg *LeaveCellReq) MarshalTo(data []byte) (n int, err error) {
	binary.LittleEndian.PutUint64(data[0:8], msg.EntityID)
	return msg.Size(), nil
}

// Unmarshal 反序列化
func (msg *LeaveCellReq) Unmarshal(data []byte) error {
	msg.EntityID = binary.LittleEndian.Uint64(data[0:8])
	return nil
}

// Size 获取长度
func (msg *LeaveCellReq) Size() (n int) {
	return 8
}

// SpaceUserConnect 空间玩家专用连接成功，由客户端发给服务器
type SpaceUserConnect struct {
	UID     uint64
	SpaceID uint64
}

func (msg *SpaceUserConnect) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *SpaceUserConnect) Name() string {
	return "SpaceUserConnect"
}

// MarshalTo 序列化
func (msg *SpaceUserConnect) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *SpaceUserConnect) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *SpaceUserConnect) Size() (n int) {
	return common.CalcSize(msg)
}

// SpaceUserConnectSucceedRet 空间玩家专用连接成功
type SpaceUserConnectSucceedRet struct {
}

func (msg *SpaceUserConnectSucceedRet) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *SpaceUserConnectSucceedRet) Name() string {
	return "SpaceUserConnectSucceedRet"
}

// MarshalTo 序列化
func (msg *SpaceUserConnectSucceedRet) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *SpaceUserConnectSucceedRet) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *SpaceUserConnectSucceedRet) Size() (n int) {
	return common.CalcSize(msg)
}

// SyncClock 对时
type SyncClock struct {
	TimeStamp uint32
}

func (msg *SyncClock) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *SyncClock) Name() string {
	return "SyncClock"
}

// MarshalTo 序列化
func (msg *SyncClock) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *SyncClock) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *SyncClock) Size() (n int) {
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////

// SyncUserState 同步信息
type SyncUserState struct {
	EntityID uint64
	Data     []byte
}

func (msg *SyncUserState) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *SyncUserState) Name() string {
	return "SyncUserState"
}

// MarshalTo 序列化
func (msg *SyncUserState) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *SyncUserState) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *SyncUserState) Size() (n int) {
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////////////////

func NewAOISyncUserState() *AOISyncUserState {
	return &AOISyncUserState{
		//Data: make([]byte, 0, 1000),
		Num:  0,
		EIDS: make([]uint64, 0, 10),
		EDS:  make([][]byte, 0, 10),
	}
}

// AOISyncUserState AOI同步信息
type AOISyncUserState struct {
	Num  uint32
	EIDS []uint64
	EDS  [][]byte
}

func (msg *AOISyncUserState) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *AOISyncUserState) Name() string {
	return "AOISyncUserState"
}

// MarshalTo 序列化
func (msg *AOISyncUserState) MarshalTo(data []byte) (n int, err error) {
	//bw := common.NewByteStream(data)
	//return msg.Size(), bw.Marshal(msg)

	size := msg.Size()
	bw := common.NewByteStream(data)
	bw.WriteUInt32(msg.Num)

	for i := 0; i < int(msg.Num); i++ {
		bw.WriteUInt64(msg.EIDS[i])
		bw.WriteBytes(msg.EDS[i])
	}

	return size, nil
}

// Unmarshal 反序列化
func (msg *AOISyncUserState) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	var err error
	msg.Num, err = br.ReadUInt32()
	if err != nil {
		return err
	}
	var oneData []byte
	for i := uint32(0); i < msg.Num; i++ {
		eid, err := br.ReadUInt64()
		if err != nil {
			return err
		}
		msg.EIDS = append(msg.EIDS, eid)

		oneData, err = br.ReadBytes()
		if err != nil {
			return err
		}
		msg.EDS = append(msg.EDS, oneData)
	}
	return nil
}

// Size 获取长度
func (msg *AOISyncUserState) Size() (n int) {

	sum := 0
	for i := 0; i < len(msg.EDS); i++ {
		sum += len(msg.EDS[i])
	}

	return 4 + len(msg.EIDS)*8 + len(msg.EIDS)*2 + sum
}

// AddData 增加一个玩家数据
func (msg *AOISyncUserState) AddData(id uint64, data []byte) {
	msg.Num++
	msg.EIDS = append(msg.EIDS, id)
	msg.EDS = append(msg.EDS, data)
}

// Clear 清理
func (msg *AOISyncUserState) Clear() {
	msg.Num = 0
	msg.EIDS = msg.EIDS[0:0]
	msg.EDS = msg.EDS[0:0]
}

////////////////////////////////////////////////////////////////

// AdjustUserState 纠正状态
type AdjustUserState struct {
	Data []byte
}

func (msg *AdjustUserState) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *AdjustUserState) Name() string {
	return "AdjustUserState"
}

// MarshalTo 序列化
func (msg *AdjustUserState) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *AdjustUserState) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *AdjustUserState) Size() (n int) {
	return common.CalcSize(msg)
}

////////////////////////////////////////////////////////////////

// CreateGhostReq 请求创建ghost
type CreateGhostReq struct {
	EntityType   string
	EntityID     uint64
	DBID         uint64
	InitParam    []byte
	CellID       uint64
	Pos          linmath.Vector3
	RealServerID uint64
	RealCellID   uint64
	PropNum      uint32
	Props        []byte
}

func (msg *CreateGhostReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *CreateGhostReq) Name() string {
	return "CreateGhostReq"
}

// MarshalTo 序列化
func (msg *CreateGhostReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *CreateGhostReq) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *CreateGhostReq) Size() (n int) {
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////////

// DeleteGhostReq 请求删除ghost
type DeleteGhostReq struct {
	EntityID uint64
}

func (msg *DeleteGhostReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *DeleteGhostReq) Name() string {
	return "DeleteGhostReq"
}

// MarshalTo 序列化
func (msg *DeleteGhostReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *DeleteGhostReq) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *DeleteGhostReq) Size() (n int) {
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////////

// TransferRealReq 把real转移至新的cell
type TransferRealReq struct {
	EntityID     uint64
	PropNum      uint32
	Props        []byte
	Pos          linmath.Vector3
	RealServerID uint64
	GhostNum     uint32
	GhostData    []byte
}

func (msg *TransferRealReq) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *TransferRealReq) Name() string {
	return "TransferRealReq"
}

// MarshalTo 序列化
func (msg *TransferRealReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *TransferRealReq) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *TransferRealReq) Size() (n int) {
	return common.CalcSize(msg)
}

//////////////////////////////////////////////////////////////

// NewRealNotify 通知有新的real出现
type NewRealNotify struct {
	RealServerID uint64
	RealCellID   uint64
}

func (msg *NewRealNotify) String() string {
	return fmt.Sprintf("%+v", *msg)
}

// Name 获取消息名称
func (msg *NewRealNotify) Name() string {
	return "NewRealNotify"
}

// MarshalTo 序列化
func (msg *NewRealNotify) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *NewRealNotify) Unmarshal(data []byte) error {
	br := common.NewByteStream(data)
	return br.Unmarshal(msg)
}

// Size 获取长度
func (msg *NewRealNotify) Size() (n int) {
	return common.CalcSize(msg)
}
