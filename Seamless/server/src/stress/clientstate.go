package main

import (
	"encoding/binary"
	"math"
	"zeus/common"
	"zeus/linmath"
	"zeus/space"

	log "github.com/cihub/seelog"
)

const (
	RoomPlayerStateMask_Speed_X = space.EntityStateMask_Reserve + 1
	RoomPlayerStateMask_Speed_Y = space.EntityStateMask_Reserve + 2
	RoomPlayerStateMask_Speed_Z = space.EntityStateMask_Reserve + 3

	RoomPlayerStateMask_Base_State   = space.EntityStateMask_Reserve + 7
	RoomPlayerStateMask_Action_State = space.EntityStateMask_Reserve + 9

	RoomPlayerStateMask_Input_LJS_X   = space.EntityStateMask_Reserve + 10
	RoomPlayerStateMask_Input_LJS_Y   = space.EntityStateMask_Reserve + 11
	RoomPlayerStateMask_Input_RJS_X   = space.EntityStateMask_Reserve + 12
	RoomPlayerStateMask_Input_RJS_Y   = space.EntityStateMask_Reserve + 13
	RoomPlayerStateMask_Input_Button  = space.EntityStateMask_Reserve + 14
	RoomPlayerStateMask_WheelF_Rota_X = space.EntityStateMask_Reserve + 15
	RoomPlayerStateMask_WheelF_Rota_Y = space.EntityStateMask_Reserve + 16
	RoomPlayerStateMask_WheelF_Rota_Z = space.EntityStateMask_Reserve + 17
	RoomPlayerStateMask_WheelB_Rota_X = space.EntityStateMask_Reserve + 18
	RoomPlayerStateMask_WheelB_Rota_Y = space.EntityStateMask_Reserve + 19
	RoomPlayerStateMask_WheelB_Rota_Z = space.EntityStateMask_Reserve + 20
)

const (
	//RoomPlayerBaseState_Inplane 跳伞准备
	RoomPlayerBaseState_Inplane = 1
	//RoomPlayerBaseState_Glide 跳伞俯冲
	RoomPlayerBaseState_Glide = 2
	//RoomPlayerBaseState_Parachute 跳伞
	RoomPlayerBaseState_Parachute = 3
	//RoomPlayerBaseState_Stand 正常状态(站立，移动)
	RoomPlayerBaseState_Stand = 4
	//RoomPlayerBaseState_Down 匍匐
	RoomPlayerBaseState_Down = 5
	//RoomPlayerBaseState_Ride 载具
	RoomPlayerBaseState_Ride = 6
	//RoomPlayerBaseState_Passenger 乘客
	RoomPlayerBaseState_Passenger = 7
	//RoomPlayerBaseState_Swim 游泳
	RoomPlayerBaseState_Swim = 8
	//RoomPlayerBaseState_Dead 死亡
	RoomPlayerBaseState_Dead = 9
	//RoomPlayerBaseState_WillDie 被击倒
	RoomPlayerBaseState_WillDie = 10
	//RoomPlayerBaseState_Watch 观战
	RoomPlayerBaseState_Watch = 11
	//RoomPlayerBaseState_Crouch 蹲
	RoomPlayerBaseState_Crouch = 12
	//RoomPlayerBaseState_Fall 跌落
	RoomPlayerBaseState_Fall = 13
	//RoomPlayerBaseState_Jump 跳跃
	RoomPlayerBaseState_Jump = 14
	//RoomPlayerBaseState_LeaveMap 离开地图
	RoomPlayerBaseState_LeaveMap = 15
	//RoomPlayerBaseState_LoadingMap 加载地图
	RoomPlayerBaseState_LoadingMap = 100
)

// RoomPlayerState 场景玩家状态
type RoomPlayerState struct {
	space.EntityState

	Speed       linmath.Vector3
	BaseState   byte
	ActionState byte
	LeftJS      linmath.Vector2
	RightJS     linmath.Vector2
	Button      byte
	WheelFRota  linmath.Vector3
	WheelBRota  linmath.Vector3
}

// NewRoomPlayerState 新建状态
func NewRoomPlayerState() space.IEntityState {
	// state := &RoomPlayerState{}
	// state.Init(state)

	// state.Bind("Speed.X", RoomPlayerStateMask_Speed_X)
	// state.Bind("Speed.Y", RoomPlayerStateMask_Speed_Y)
	// state.Bind("Speed.Z", RoomPlayerStateMask_Speed_Z)

	// state.Bind("BaseState", RoomPlayerStateMask_Base_State)
	// state.Bind("ActionState", RoomPlayerStateMask_Action_State)

	// state.Bind("LeftJS.X", RoomPlayerStateMask_Input_LJS_X)
	// state.Bind("LeftJS.Y", RoomPlayerStateMask_Input_LJS_Y)
	// state.Bind("RightJS.X", RoomPlayerStateMask_Input_RJS_X)
	// state.Bind("RightJS.Y", RoomPlayerStateMask_Input_RJS_Y)
	// state.Bind("Button", RoomPlayerStateMask_Input_Button)

	// state.Bind("WheelFRota.X", RoomPlayerStateMask_WheelF_Rota_X)
	// state.Bind("WheelFRota.Y", RoomPlayerStateMask_WheelF_Rota_Y)
	// state.Bind("WheelFRota.Z", RoomPlayerStateMask_WheelF_Rota_Z)
	// state.Bind("WheelBRota.X", RoomPlayerStateMask_WheelB_Rota_X)
	// state.Bind("WheelBRota.Y", RoomPlayerStateMask_WheelB_Rota_Y)
	// state.Bind("WheelBRota.Z", RoomPlayerStateMask_WheelB_Rota_Z)

	return &RoomPlayerState{}
}

func (state *RoomPlayerState) CopyTo(n space.IEntityState) {
	o := n.(*RoomPlayerState)
	o.TimeStamp = state.TimeStamp
	o.Pos = state.Pos
	o.Rota = state.Rota
	o.Param1 = state.Param1
	o.Param2 = state.Param2
	o.Speed = state.Speed
	o.BaseState = state.BaseState
	o.ActionState = state.ActionState
	o.LeftJS = state.LeftJS
	o.RightJS = state.RightJS
	o.Button = state.Button
	o.WheelFRota = state.WheelFRota
	o.WheelBRota = state.WheelBRota
}

// Clone 克隆
func (state *RoomPlayerState) Clone() space.IEntityState {
	o := &RoomPlayerState{}
	state.CopyTo(o)
	return o
}

// Combine 合并
func (state *RoomPlayerState) Combine(data []byte) {
	bs := common.NewByteStream(data)

	state.TimeStamp, _ = bs.ReadUInt32()
	mask, _ := bs.ReadInt32()

	for i := 0; i < space.EntityStateMask_Max; i++ {
		if (mask & (1 << uint(i))) != 0 {
			if i <= space.EntityStateMask_Reserve {
				state.SetBaseValue(i, bs)
			} else {
				state.SetExtValue(i, bs)
			}
		}
	}
}

// Delta 求差异
func (state *RoomPlayerState) Delta(o space.IEntityState) ([]byte, bool) {
	n := o.(*RoomPlayerState)

	bs := common.NewByteStream(make([]byte, space.EntityStateMask_Max*4))
	var mask int32
	isEqual := true

	bs.WriteUInt32(o.GetTimeStamp())
	bs.WriteInt32(mask)

	for i := 0; i < space.EntityStateMask_Max; i++ {
		if i <= space.EntityStateMask_Reserve {
			if !state.CompareAndSetBaseValueDelta(&n.EntityState, &mask, uint32(i), bs) {
				isEqual = false
			}
		} else {
			if !state.CompareAndSetExtValueDelta(n, &mask, uint32(i), bs) {
				isEqual = false
			}
		}
	}

	used := bs.GetUsedSlice()
	binary.LittleEndian.PutUint32(used[4:8], uint32(mask))
	return used, isEqual
}

// Marshal 序列化
func (state *RoomPlayerState) Marshal() []byte {
	bs := common.NewByteStream(make([]byte, space.EntityStateMask_Max*4))
	var mask int32
	bs.WriteUInt32(state.TimeStamp)
	bs.WriteInt32(mask)

	for i := 0; i < space.EntityStateMask_Max; i++ {
		if i <= space.EntityStateMask_Reserve {
			state.WriteBaseValue(&mask, uint32(i), bs)
		} else {
			state.WriteExtValue(&mask, uint32(i), bs)
		}
	}

	used := bs.GetUsedSlice()
	binary.LittleEndian.PutUint32(used[4:8], uint32(mask))

	return used
}

// CompareAndSetExtValueDelta 比较和设置额外的值
func (state *RoomPlayerState) CompareAndSetExtValueDelta(o *RoomPlayerState, mask *int32, maskoffset uint32, bs *common.ByteStream) bool {
	var oldfloat float32
	var newfloat float32
	var oldByte byte
	var newByte byte
	var t int

	switch maskoffset {
	case RoomPlayerStateMask_Speed_X:
		oldfloat = state.Speed.X
		newfloat = o.Speed.X
		t = 1
	case RoomPlayerStateMask_Speed_Y:
		oldfloat = state.Speed.Y
		newfloat = o.Speed.Y
		t = 1
	case RoomPlayerStateMask_Speed_Z:
		oldfloat = state.Speed.Z
		newfloat = o.Speed.Z
		t = 1
	case RoomPlayerStateMask_Base_State:
		oldByte = state.BaseState
		newByte = o.BaseState
		t = 2
	case RoomPlayerStateMask_Action_State:
		oldByte = state.ActionState
		newByte = o.ActionState
		t = 2
	case RoomPlayerStateMask_Input_LJS_X:
		oldfloat = state.LeftJS.X
		newfloat = o.LeftJS.X
		t = 1
	case RoomPlayerStateMask_Input_LJS_Y:
		oldfloat = state.LeftJS.Y
		newfloat = o.LeftJS.Y
		t = 1
	case RoomPlayerStateMask_Input_RJS_X:
		oldfloat = state.RightJS.X
		newfloat = o.RightJS.X
		t = 1
	case RoomPlayerStateMask_Input_RJS_Y:
		oldfloat = state.RightJS.Y
		newfloat = o.RightJS.Y
		t = 1
	case RoomPlayerStateMask_Input_Button:
		oldByte = state.Button
		newByte = o.Button
		t = 2
	case RoomPlayerStateMask_WheelF_Rota_X:
		oldfloat = state.WheelFRota.X
		newfloat = o.WheelFRota.X
		t = 1
	case RoomPlayerStateMask_WheelF_Rota_Y:
		oldfloat = state.WheelFRota.Y
		newfloat = o.WheelFRota.Y
		t = 1
	case RoomPlayerStateMask_WheelF_Rota_Z:
		oldfloat = state.WheelFRota.Z
		newfloat = o.WheelFRota.Z
		t = 1
	case RoomPlayerStateMask_WheelB_Rota_X:
		oldfloat = state.WheelBRota.X
		newfloat = o.WheelBRota.X
		t = 1
	case RoomPlayerStateMask_WheelB_Rota_Y:
		oldfloat = state.WheelBRota.Y
		newfloat = o.WheelBRota.Y
		t = 1
	case RoomPlayerStateMask_WheelB_Rota_Z:
		oldfloat = state.WheelBRota.Z
		newfloat = o.WheelBRota.Z
		t = 1
	default:
		return true
	}

	if t == 1 {
		if math.Abs(float64(oldfloat-newfloat)) <= 0.001 {
			return true
		}
		bs.WriteFloat32(newfloat)
	} else if t == 2 {
		if oldByte == newByte {
			return true
		}
		bs.WriteByte(newByte)
	}

	(*mask) |= 1 << maskoffset
	return false
}

// WriteExtValue 设置额外的值
func (state *RoomPlayerState) WriteExtValue(mask *int32, maskoffset uint32, bs *common.ByteStream) bool {
	var newfloat float32
	var newByte byte
	var t int

	switch maskoffset {
	case RoomPlayerStateMask_Speed_X:
		newfloat = state.Speed.X
		t = 1
	case RoomPlayerStateMask_Speed_Y:
		newfloat = state.Speed.Y
		t = 1
	case RoomPlayerStateMask_Speed_Z:
		newfloat = state.Speed.Z
		t = 1
	case RoomPlayerStateMask_Base_State:
		newByte = state.BaseState
		t = 2
	case RoomPlayerStateMask_Action_State:
		newByte = state.ActionState
		t = 2
	case RoomPlayerStateMask_Input_LJS_X:
		newfloat = state.LeftJS.X
		t = 1
	case RoomPlayerStateMask_Input_LJS_Y:
		newfloat = state.LeftJS.Y
		t = 1
	case RoomPlayerStateMask_Input_RJS_X:
		newfloat = state.RightJS.X
		t = 1
	case RoomPlayerStateMask_Input_RJS_Y:
		newfloat = state.RightJS.Y
		t = 1
	case RoomPlayerStateMask_Input_Button:
		newByte = state.Button
		t = 2
	case RoomPlayerStateMask_WheelF_Rota_X:
		newfloat = state.WheelFRota.X
		t = 1
	case RoomPlayerStateMask_WheelF_Rota_Y:
		newfloat = state.WheelFRota.Y
		t = 1
	case RoomPlayerStateMask_WheelF_Rota_Z:
		newfloat = state.WheelFRota.Z
		t = 1
	case RoomPlayerStateMask_WheelB_Rota_X:
		newfloat = state.WheelBRota.X
		t = 1
	case RoomPlayerStateMask_WheelB_Rota_Y:
		newfloat = state.WheelBRota.Y
		t = 1
	case RoomPlayerStateMask_WheelB_Rota_Z:
		newfloat = state.WheelBRota.Z
		t = 1
	default:
		return true
	}

	if t == 1 {
		bs.WriteFloat32(newfloat)
	} else if t == 2 {
		bs.WriteByte(newByte)
	}

	(*mask) |= 1 << maskoffset
	return false
}

// SetExtValue 设置额外的值
func (state *RoomPlayerState) SetExtValue(mask int, bs *common.ByteStream) {
	switch mask {
	case RoomPlayerStateMask_Speed_X:
		state.Speed.X, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Speed_Y:
		state.Speed.Y, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Speed_Z:
		state.Speed.Z, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Base_State:
		state.BaseState, _ = bs.ReadByte()
	case RoomPlayerStateMask_Action_State:
		state.ActionState, _ = bs.ReadByte()
	case RoomPlayerStateMask_Input_LJS_X:
		state.LeftJS.X, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Input_LJS_Y:
		state.LeftJS.Y, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Input_RJS_X:
		state.RightJS.X, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Input_RJS_Y:
		state.RightJS.Y, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_Input_Button:
		state.Button, _ = bs.ReadByte()
	case RoomPlayerStateMask_WheelF_Rota_X:
		state.WheelFRota.X, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_WheelF_Rota_Y:
		state.WheelFRota.Y, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_WheelF_Rota_Z:
		state.WheelFRota.Z, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_WheelB_Rota_X:
		state.WheelBRota.X, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_WheelB_Rota_Y:
		state.WheelBRota.Y, _ = bs.ReadFloat32()
	case RoomPlayerStateMask_WheelB_Rota_Z:
		state.WheelBRota.Z, _ = bs.ReadFloat32()
	default:
		log.Error("Set ext value failed ", mask)
	}
}
