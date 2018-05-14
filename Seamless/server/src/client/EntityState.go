package main

import (
	"errors"
	"math"
	"reflect"
	"zeus/common"
	"zeus/linmath"

	log "github.com/cihub/seelog"
)

// IEntityState 实体状态快照
type IEntityState interface {
	GetPos() linmath.Vector3
	SetPos(linmath.Vector3)
	GetRota() linmath.Vector3
	SetRota(linmath.Vector3)

	GetTimeStamp() uint32
	SetTimeStamp(uint32)

	SetDirty(bool)
	IsDirty() bool

	SetModify(bool)
	IsModify() bool

	Clone() IEntityState
	CopyTo(IEntityState)
	Combine([]byte)
	Delta(IEntityState) ([]byte, bool)
	Marshal() []byte
}

const (
	EntityStateMask_Max = 32

	EntityStateMask_Pos_X = 0
	EntityStateMask_Pos_Y = 1
	EntityStateMask_Pos_Z = 2

	EntityStateMask_Rota_X = 3
	EntityStateMask_Rota_Y = 4
	EntityStateMask_Rota_Z = 5

	EntityStateMask_Param1 = 6
	EntityStateMask_Param2 = 7

	EntityStateMask_Reserve = 10
)

var ErrorMaskOffsetExceed = errors.New("mask offset exceed")
var ErrorMaskOffsetInvalid = errors.New("mask offset invalid")

// EntityState 由后代的EntityState包含
type EntityState struct {
	isDirty   bool
	isModify  bool
	TimeStamp uint32

	Pos  linmath.Vector3
	Rota linmath.Vector3

	Param1 uint64
	Param2 uint64

	// Events []byte
}

// GetPos 获取位置
func (s *EntityState) GetPos() linmath.Vector3 {
	return s.Pos
}

// SetPos 设置位置
func (s *EntityState) SetPos(pos linmath.Vector3) {
	s.Pos = pos
}

// GetRota 获取旋转
func (s *EntityState) GetRota() linmath.Vector3 {
	return s.Rota
}

// SetRota 设置旋转
func (s *EntityState) SetRota(rota linmath.Vector3) {
	s.Rota = rota
}

// GetTimeStamp 获取时间戳
func (s *EntityState) GetTimeStamp() uint32 {
	return s.TimeStamp
}

// SetTimeStamp 设置时间戳
func (s *EntityState) SetTimeStamp(timeStamp uint32) {
	s.TimeStamp = timeStamp
}

// SetParam1Uint64 设置参数1
func (s *EntityState) SetParam1Uint64(v uint64) {
	s.Param1 = v
}

// GetParam1Uint64 获取参数1
func (s *EntityState) GetParam1Uint64() uint64 {
	return s.Param1
}

// SetParam2Uint64 设置参数2
func (s *EntityState) SetParam2Uint64(v uint64) {
	s.Param2 = v
}

// GetParam2Uint64 获取参数2
func (s *EntityState) GetParam2Uint64() uint64 {
	return s.Param2
}

// SetParam1Uint32 设置参数1
func (s *EntityState) SetParam1Uint32(v uint32) {
	s.Param1 = uint64(v)
}

// GetParam1Uint32 获取参数1
func (s *EntityState) GetParam1Uint32() uint32 {
	return uint32(s.Param1)
}

// SetParam2Uint32 设置参数2
func (s *EntityState) SetParam2Uint32(v uint32) {
	s.Param2 = uint64(v)
}

// GetParam2Uint32 获取参数2
func (s *EntityState) GetParam2Uint32() uint32 {
	return uint32(s.Param2)
}

// SetParam1Float64 设置参数1
func (s *EntityState) SetParam1Float64(v float64) {
	s.Param1 = math.Float64bits(v)
}

// GetParam1Float64 获取参数1
func (s *EntityState) GetParam1Float64() float64 {
	return math.Float64frombits(s.Param1)
}

// SetParam2Float64 设置参数2
func (s *EntityState) SetParam2Float64(v float64) {
	s.Param2 = math.Float64bits(v)
}

// GetParam2Float64 获取参数2
func (s *EntityState) GetParam2Float64() float64 {
	return math.Float64frombits(s.Param2)
}

// SetParam1Float32 设置参数1
func (s *EntityState) SetParam1Float32(v float32) {
	s.Param1 = uint64(math.Float32bits(v))
}

// GetParam1Float32 获取参数1
func (s *EntityState) GetParam1Float32() float32 {
	return math.Float32frombits(uint32(s.Param1))
}

// SetParam2Float32 设置参数2
func (s *EntityState) SetParam2Float32(v float32) {
	s.Param2 = uint64(math.Float32bits(v))
}

// GetParam2Float32 获取参数2
func (s *EntityState) GetParam2Float32() float32 {
	return math.Float32frombits(uint32(s.Param2))
}

// SetDirty 设置脏标记
func (s *EntityState) SetDirty(dirty bool) {
	s.isDirty = dirty
}

// IsDirty 是否脏标记
func (s *EntityState) IsDirty() bool {
	return s.isDirty
}

// SetModify 设置被修改过
func (s *EntityState) SetModify(m bool) {
	s.isModify = m
	if m == true {
		s.isDirty = true
	}
}

// IsModify 是否被修改过
func (s *EntityState) IsModify() bool {
	return s.isModify
}

// SetBaseValue 设置基础值
func (s *EntityState) SetBaseValue(mask int, bs *common.ByteStream) {
	switch mask {
	case EntityStateMask_Pos_X:
		s.Pos.X, _ = bs.ReadFloat32()
	case EntityStateMask_Pos_Y:
		s.Pos.Y, _ = bs.ReadFloat32()
	case EntityStateMask_Pos_Z:
		s.Pos.Z, _ = bs.ReadFloat32()
	case EntityStateMask_Rota_X:
		s.Rota.X, _ = bs.ReadFloat32()
	case EntityStateMask_Rota_Y:
		s.Rota.Y, _ = bs.ReadFloat32()
	case EntityStateMask_Rota_Z:
		s.Rota.Z, _ = bs.ReadFloat32()
	case EntityStateMask_Param1:
		s.Param1, _ = bs.ReadUInt64()
	case EntityStateMask_Param2:
		s.Param2, _ = bs.ReadUInt64()
	default:
		log.Error("Set base value failed ", mask)
	}
}

// CompareAndSetBaseValueDelta 比较基础值
func (s *EntityState) CompareAndSetBaseValueDelta(o *EntityState, mask *int32, maskoffset uint32, bs *common.ByteStream) bool {
	var oldfloat float32
	var newfloat float32
	var olduint uint64
	var newuint uint64
	var oldbytes []byte
	var newbytes []byte
	var t int

	switch maskoffset {
	case EntityStateMask_Pos_X:
		oldfloat = s.Pos.X
		newfloat = o.Pos.X
		t = 1
	case EntityStateMask_Pos_Y:
		oldfloat = s.Pos.Y
		newfloat = o.Pos.Y
		t = 1
	case EntityStateMask_Pos_Z:
		oldfloat = s.Pos.Z
		newfloat = o.Pos.Z
		t = 1
	case EntityStateMask_Rota_X:
		oldfloat = s.Rota.X
		newfloat = o.Rota.X
		t = 1
	case EntityStateMask_Rota_Y:
		oldfloat = s.Rota.Y
		newfloat = o.Rota.Y
		t = 1
	case EntityStateMask_Rota_Z:
		oldfloat = s.Rota.Z
		newfloat = o.Rota.Z
		t = 1
	case EntityStateMask_Param1:
		olduint = s.Param1
		newuint = o.Param1
		t = 2
	case EntityStateMask_Param2:
		olduint = s.Param2
		newuint = o.Param2
		t = 2
	default:
		return true
	}

	if t == 1 {
		if math.Abs(float64(oldfloat-newfloat)) <= 0.001 {
			return true
		}
		bs.WriteFloat32(newfloat)
	} else if t == 2 {
		if olduint == newuint {
			return true

		}
		bs.WriteUInt64(newuint)
	} else if t == 3 {
		if reflect.DeepEqual(oldbytes, newbytes) {
			return true
		}
		bs.WriteBytes(newbytes)
	}

	(*mask) |= 1 << maskoffset
	return false
}

// WriteBaseValue 设置基础状态
func (s *EntityState) WriteBaseValue(mask *int32, maskoffset uint32, bs *common.ByteStream) bool {
	var newfloat float32
	var newuint uint64
	var newbytes []byte
	var t int

	switch maskoffset {
	case EntityStateMask_Pos_X:
		newfloat = s.Pos.X
		t = 1
	case EntityStateMask_Pos_Y:
		newfloat = s.Pos.Y
		t = 1
	case EntityStateMask_Pos_Z:
		newfloat = s.Pos.Z
		t = 1
	case EntityStateMask_Rota_X:
		newfloat = s.Rota.X
		t = 1
	case EntityStateMask_Rota_Y:
		newfloat = s.Rota.Y
		t = 1
	case EntityStateMask_Rota_Z:
		newfloat = s.Rota.Z
		t = 1
	case EntityStateMask_Param1:
		newuint = s.Param1
		t = 2
	case EntityStateMask_Param2:
		newuint = s.Param2
		t = 2
	default:
		return true
	}

	if t == 1 {
		bs.WriteFloat32(newfloat)
	} else if t == 2 {
		bs.WriteUInt64(newuint)
	} else if t == 3 {
		bs.WriteBytes(newbytes)
	}

	(*mask) |= 1 << maskoffset
	return false
}
