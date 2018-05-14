package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"zeus/linmath"
)

// ByteStream 包装一些字节流的读写方法，方便二进制消息的序列和反序列
type ByteStream struct {
	data     []byte
	readPos  uint32
	writePos uint32
}

// NewByteStream 创建一个新的字节流
func NewByteStream(data []byte) *ByteStream {
	return &ByteStream{
		data:     data,
		readPos:  0,
		writePos: 0,
	}
}

func (bs *ByteStream) readCheck(c uint32) error {
	if bs.data == nil {
		return errors.New("data is nil")
	}

	if bs.readPos+c > uint32(len(bs.data)) {
		return errors.New("no enough read space")
	}

	return nil
}

// ReadEnd 是否读完
func (bs *ByteStream) ReadEnd() bool {
	return bs.readPos != uint32(len(bs.data))
}

// ReadByte 读取一个字节
func (bs *ByteStream) ReadByte() (byte, error) {
	if err := bs.readCheck(1); err != nil {
		return 0, err
	}

	v := bs.data[bs.readPos]
	bs.readPos = bs.readPos + 1

	return v, nil
}

// ReadBool 读bool
func (bs *ByteStream) ReadBool() (bool, error) {
	v, err := bs.ReadByte()
	if err != nil {
		return false, err
	}

	return (v != 0), nil
}

// ReadInt8 读取一个int8
func (bs *ByteStream) ReadInt8() (int8, error) {
	v, err := bs.ReadByte()
	return int8(v), err
}

// ReadInt16 读取一个int16
func (bs *ByteStream) ReadInt16() (int16, error) {
	v, err := bs.ReadUInt16()
	return int16(v), err
}

// ReadInt32 读取一个int32
func (bs *ByteStream) ReadInt32() (int32, error) {
	v, err := bs.ReadUInt32()
	return int32(v), err
}

// ReadInt64 读取一个int64
func (bs *ByteStream) ReadInt64() (int64, error) {
	v, err := bs.ReadUInt64()
	return int64(v), err
}

// ReadUInt16 读取一个UInt16
func (bs *ByteStream) ReadUInt16() (uint16, error) {
	if err := bs.readCheck(2); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint16(bs.data[bs.readPos : bs.readPos+2])
	bs.readPos = bs.readPos + 2

	return v, nil
}

// ReadUInt32 读取一个Int
func (bs *ByteStream) ReadUInt32() (uint32, error) {
	if err := bs.readCheck(4); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint32(bs.data[bs.readPos : bs.readPos+4])
	bs.readPos = bs.readPos + 4

	return v, nil
}

// ReadUInt64 读取一个Uint64
func (bs *ByteStream) ReadUInt64() (uint64, error) {
	if err := bs.readCheck(8); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint64(bs.data[bs.readPos : bs.readPos+8])
	bs.readPos = bs.readPos + 8

	return v, nil
}

// ReadStr 读取一个string
func (bs *ByteStream) ReadStr() (string, error) {

	len, err := bs.ReadUInt16()
	if err != nil {
		return "", err
	}

	if err = bs.readCheck(uint32(len)); err != nil {
		return "", err
	}

	v := string(bs.data[bs.readPos : bs.readPos+uint32(len)])
	bs.readPos = bs.readPos + uint32(len)

	return v, nil
}

// ReadBytes 读取一个byte[]
func (bs *ByteStream) ReadBytes() ([]byte, error) {
	len, err := bs.ReadUInt16()
	if err != nil {
		return nil, err
	}

	if len == 0 {
		return nil, nil
	}

	if err := bs.readCheck(uint32(len)); err != nil {
		return nil, err
	}

	b := make([]byte, len)

	copy(b, bs.data[bs.readPos:bs.readPos+uint32(len)])
	bs.readPos = bs.readPos + uint32(len)

	return b, nil
}

// ReadFloat32 读取float32
func (bs *ByteStream) ReadFloat32() (float32, error) {
	u, err := bs.ReadUInt32()
	if err != nil {
		return 0, err
	}

	f := math.Float32frombits(u)
	return f, nil
}

// ReadFloat64 读取float64
func (bs *ByteStream) ReadFloat64() (float64, error) {
	u, err := bs.ReadUInt64()
	if err != nil {
		return 0, err
	}

	f := math.Float64frombits(u)
	return f, nil
}

// ReadVector3 读取一个坐标
func (bs *ByteStream) ReadVector3() (linmath.Vector3, error) {
	pos := linmath.Vector3{}
	var err error

	pos.X, err = bs.ReadFloat32()
	if err != nil {
		return pos, err
	}
	pos.Y, err = bs.ReadFloat32()
	if err != nil {
		return pos, err
	}
	pos.Z, err = bs.ReadFloat32()
	if err != nil {
		return pos, err
	}

	return pos, nil
}

func (bs *ByteStream) writeCheck(c uint32) error {
	if bs.data == nil {
		return errors.New("data is nil")
	}

	if bs.writePos+c > uint32(len(bs.data)) {
		return errors.New("no enough write space")
	}

	return nil
}

// WriteByte 写字节
func (bs *ByteStream) WriteByte(v byte) error {
	if err := bs.writeCheck(1); err != nil {
		return err
	}

	bs.data[bs.writePos] = v
	bs.writePos = bs.writePos + 1

	return nil
}

// WriteBool 写bool, 1代表true, 0代表false
func (bs *ByteStream) WriteBool(v bool) error {
	if v {
		return bs.WriteByte(1)
	}
	return bs.WriteByte(0)
}

// WriteInt8 写Int8
func (bs *ByteStream) WriteInt8(v int8) error {
	return bs.WriteByte(byte(v))
}

// WriteInt16 写Int16
func (bs *ByteStream) WriteInt16(v int16) error {
	return bs.WriteUInt16(uint16(v))
}

// WriteInt32 写Int32
func (bs *ByteStream) WriteInt32(v int32) error {
	return bs.WriteUInt32(uint32(v))
}

// WriteInt64 写Int64
func (bs *ByteStream) WriteInt64(v int64) error {
	return bs.WriteUInt64(uint64(v))
}

// WriteUInt16 写Uint16
func (bs *ByteStream) WriteUInt16(v uint16) error {
	if err := bs.writeCheck(2); err != nil {
		return err
	}

	binary.LittleEndian.PutUint16(bs.data[bs.writePos:bs.writePos+2], v)
	bs.writePos = bs.writePos + 2

	return nil
}

// WriteUInt32 写Uint32
func (bs *ByteStream) WriteUInt32(v uint32) error {

	if err := bs.writeCheck(4); err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(bs.data[bs.writePos:bs.writePos+4], v)
	bs.writePos = bs.writePos + 4

	return nil
}

// WriteUInt64 写Uint64
func (bs *ByteStream) WriteUInt64(v uint64) error {

	if err := bs.writeCheck(8); err != nil {
		return err
	}

	binary.LittleEndian.PutUint64(bs.data[bs.writePos:bs.writePos+8], v)
	bs.writePos = bs.writePos + 8

	return nil
}

// WriteStr 写string
func (bs *ByteStream) WriteStr(v string) error {

	if err := bs.writeCheck(uint32(len(v) + 2)); err != nil {
		return err
	}

	bs.WriteUInt16(uint16(len(v)))

	if len(v) != 0 {
		copy(bs.data[bs.writePos:bs.writePos+uint32(len(v))], v)
		bs.writePos = bs.writePos + uint32(len(v))
	}

	return nil
}

// WriteBytes 写[]byte
func (bs *ByteStream) WriteBytes(v []byte) error {

	if v == nil {
		bs.WriteUInt16(0)
		return nil
	}

	if err := bs.writeCheck(uint32(len(v) + 2)); err != nil {
		return err
	}

	bs.WriteUInt16(uint16(len(v)))

	copy(bs.data[bs.writePos:bs.writePos+uint32(len(v))], v)

	bs.writePos = bs.writePos + uint32(len(v))
	return nil
}

// WriteFloat32 写Float32
func (bs *ByteStream) WriteFloat32(f float32) error {
	u := math.Float32bits(f)
	return bs.WriteUInt32(u)
}

// WriteFloat64 写Float64
func (bs *ByteStream) WriteFloat64(f float64) error {
	u := math.Float64bits(f)
	return bs.WriteUInt64(u)
}

// WriteVector3 写位置信息
func (bs *ByteStream) WriteVector3(pos linmath.Vector3) error {
	if err := bs.WriteFloat32(pos.X); err != nil {
		return err
	}
	if err := bs.WriteFloat32(pos.Y); err != nil {
		return err
	}
	if err := bs.WriteFloat32(pos.Z); err != nil {
		return err
	}

	return nil
}

func (bs *ByteStream) WriteMap(val map[uint32]uint32) error {
	for k, v := range val {
		if err := bs.WriteUInt32(k); err != nil {
			return err
		}

		if err := bs.WriteUInt32(v); err != nil {
			return err
		}
	}
	return nil
}

// CalcSize 计算序列化所需长度
func CalcSize(content interface{}) int {
	size := 0

	v := reflect.ValueOf(content).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()

		switch field.(type) {
		case uint8, int8, bool:
			size++
		case uint16, int16:
			size += 2
		case uint32, int32, float32:
			size += 4
		case uint64, int64, float64:
			size += 8
		case string:
			size += 2
			size += len(field.(string))
		case []byte:
			size += 2
			size += len(field.([]byte))
		case linmath.Vector3:
			size += 12
		default:
			panic(fmt.Sprintf("不支持的类型 %+v", field))
		}
	}

	return size
}

// Marshal 序列化content
func (bs *ByteStream) Marshal(content interface{}) error {
	v := reflect.ValueOf(content).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()

		var err error
		switch field.(type) {
		case bool:
			err = bs.WriteBool(v.Field(i).Interface().(bool))
		case uint8:
			err = bs.WriteByte(v.Field(i).Interface().(uint8))
		case uint16:
			err = bs.WriteUInt16(v.Field(i).Interface().(uint16))
		case uint32:
			err = bs.WriteUInt32(v.Field(i).Interface().(uint32))
		case uint64:
			err = bs.WriteUInt64(v.Field(i).Interface().(uint64))
		case string:
			err = bs.WriteStr(v.Field(i).Interface().(string))
		case float32:
			err = bs.WriteFloat32(v.Field(i).Interface().(float32))
		case float64:
			err = bs.WriteFloat64(v.Field(i).Interface().(float64))
		case []uint8:
			err = bs.WriteBytes(v.Field(i).Interface().([]byte))
		case linmath.Vector3:
			err = bs.WriteVector3(v.Field(i).Interface().(linmath.Vector3))
		default:
			panic(fmt.Sprintf("不支持的类型 %t", field))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal 反序列化
func (bs *ByteStream) Unmarshal(content interface{}) error {
	v := reflect.ValueOf(content).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()

		var err error
		var value interface{}
		switch field.(type) {
		case bool:
			value, err = bs.ReadBool()
		case uint8:
			value, err = bs.ReadByte()
		case uint16:
			value, err = bs.ReadUInt16()
		case uint32:
			value, err = bs.ReadUInt32()
		case uint64:
			value, err = bs.ReadUInt64()
		case string:
			value, err = bs.ReadStr()
		case float32:
			value, err = bs.ReadFloat32()
		case float64:
			value, err = bs.ReadFloat64()
		case []byte:
			value, err = bs.ReadBytes()
		case linmath.Vector3:
			value, err = bs.ReadVector3()
		default:
			panic(fmt.Sprintf("不支持的类型 %+v", field))
		}

		if err != nil {
			return err
		}
		v.Field(i).Set(reflect.ValueOf(value))
	}

	return nil
}

// GetUsedSlice 获取已经写入的部分Slice
func (bs *ByteStream) GetUsedSlice() []byte {
	if bs.data == nil || bs.writePos == 0 {
		return nil
	}

	return bs.data[0:bs.writePos]
}

func (bs *ByteStream) WriteInt32AtPos(value int32, pos uint32) {
	var oldWPos = bs.writePos

	bs.writePos = pos
	bs.WriteInt32(value)

	bs.writePos = oldWPos
}
