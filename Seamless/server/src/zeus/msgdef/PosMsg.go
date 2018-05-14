package msgdef

/*
	所有位置移动相关的消息
*/

// const (
// 	UM_PX = 1 << 0
// 	UM_PY = 1 << 1
// 	UM_PZ = 1 << 2
// 	UM_RX = 1 << 3
// 	UM_RY = 1 << 4
// 	UM_RZ = 1 << 5

// 	MAX_SYNC_NUM = 5
// )

// // UserMove 玩家位置移动
// type UserMove struct {
// 	Mask int32
// 	Data []byte
// }

// func (msg *UserMove) String() string {
// 	return fmt.Sprintf("%+v", *msg)
// }

// // MarshalTo 序列化
// func (msg *UserMove) MarshalTo(data []byte) (n int, err error) {
// 	return msg.Size(), nil
// }

// // Unmarshal 反序列化
// func (msg *UserMove) Unmarshal(data []byte) error {
// 	br := common.NewByteStream(data)
// 	msg.Mask, _ = br.ReadInt32()
// 	msg.Data, _ = br.ReadBytes()
// 	return nil
// }

// // Size 获取长度
// func (msg *UserMove) Size() (n int) {
// 	return 4 + len(msg.Data) + 2
// }

// // Name 获取名字
// func (msg *UserMove) Name() string {
// 	return "UserMove"
// }

// // EntityPosSet 实体位置设置
// type EntityPosSet struct {
// 	ID  uint64
// 	Pos linmath.Vector3
// }

// func (msg *EntityPosSet) String() string {
// 	return fmt.Sprintf("%+v", *msg)
// }

// // MarshalTo 序列化
// func (msg *EntityPosSet) MarshalTo(data []byte) (n int, err error) {

// 	bw := common.NewByteStream(data)
// 	bw.WriteUInt64(msg.ID)
// 	bw.WriteFloat32(msg.Pos.X)
// 	bw.WriteFloat32(msg.Pos.Y)
// 	bw.WriteFloat32(msg.Pos.Z)

// 	return msg.Size(), nil
// }

// // Unmarshal 反序列化
// func (msg *EntityPosSet) Unmarshal(data []byte) error {
// 	br := common.NewByteStream(data)
// 	msg.ID, _ = br.ReadUInt64()
// 	msg.Pos.X, _ = br.ReadFloat32()
// 	msg.Pos.Y, _ = br.ReadFloat32()
// 	msg.Pos.Z, _ = br.ReadFloat32()
// 	return nil
// }

// // Size 获取长度
// func (msg *EntityPosSet) Size() (n int) {
// 	return 8 + 3*4
// }

// // Name 获取名字
// func (msg *EntityPosSet) Name() string {
// 	return "EntityPosSet"
// }
